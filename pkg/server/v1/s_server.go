package v1

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lazybark/go-cloud-sync/pkg/fse"
	"github.com/lazybark/go-cloud-sync/pkg/fselink"
	proto "github.com/lazybark/go-cloud-sync/pkg/fselink/proto/v1"
	"github.com/lazybark/go-cloud-sync/pkg/storage"
)

func (s *FSWServer) ExtractOwnerFromPath(p string, o string) string {
	dirs := strings.Split(p, s.conf.escSymbol)
	if len(dirs) < 2 {
		return p
	} else {
		return strings.ReplaceAll(p, "?ROOT_DIR?,"+o, "?ROOT_DIR?")
	}
}

func (s *FSWServer) GetOwnerFromPath(p string) string {
	dirs := strings.Split(p, s.conf.escSymbol)
	if len(dirs) < 2 {
		return ""
	} else {
		return dirs[1]
	}
}

func (s *FSWServer) ReadLocalObjects() (objs []storage.FSObjectStored, err error) {
	local, err := s.fp.ProcessDirectory(s.conf.root)
	if err != nil {
		return nil, fmt.Errorf("[FSWATCHER][DiffListWithServer]: %w", err)
	}
	for _, o := range local {
		if o.Path == "?ROOT_DIR?" {
			//Ignore all objs in server root
			continue
		}
		objs = append(objs, storage.FSObjectStored{
			Path:        o.Path,
			Name:        o.Name,
			IsDir:       o.IsDir,
			Hash:        o.Hash,
			Owner:       s.GetOwnerFromPath(o.Path),
			Ext:         o.Ext,
			Size:        o.Size,
			FSUpdatedAt: o.UpdatedAt,
		})
	}
	return
}

func (s *FSWServer) rescanOnce() {
	local, err := s.ReadLocalObjects()
	if err != nil {
		s.extErc <- fmt.Errorf("[rescanOnce]%w", err)
		return
	}
	fmt.Printf("Root directory read. Found %d objects\n", len(local))

	err = s.stor.RefillDatabase(local)
	if err != nil {
		s.extErc <- fmt.Errorf("[rescanOnce]%w", err)
		return
	}
	fmt.Println("Database refilled")
}

func (s *FSWServer) watcherRoutine() {
	fmt.Println("Waiting for connections")
	var ev fse.FSEvent
	var m proto.ExchangeMessage
	for {
		select {
		case err, ok := <-s.srvErrChan:
			if !ok {
				return
			}
			s.extErc <- err
		case connection, ok := <-s.srvConnChan:
			if !ok {
				return
			}
			fmt.Println("NEW CONNECTION:")
			fmt.Println(connection)
			fmt.Println()
			go func() {
				for mess := range connection.MessageChan {
					err := json.Unmarshal(mess.Bytes(), &m)
					if err != nil {
						s.extErc <- err
						continue
					}

					if m.Type != proto.MessageTypeAuthReq && m.AuthKey == "" {
						s.sendError(mess.Conn(), proto.ErrForbidden)
						continue
					}

					if m.Type == proto.MessageTypeAuthReq {
						if !s.createSession("", "") {
							s.sendError(mess.Conn(), proto.ErrCodeWrongCreds)
							continue
						}
						err = fselink.SendSyncMessage(mess.Conn(), proto.MessageAuthAnswer{Success: true, AuthKey: "SOMEKEY"}, proto.MessageTypeAuthAns)
						if err != nil {
							s.extErc <- err
							continue
						}
						fmt.Println("SENT AUTH")
					} else if m.Type == proto.MessageTypeEvent {
						user, ok := s.checkToken(m.AuthKey)
						if !ok || user == "" {
							s.sendError(mess.Conn(), proto.ErrForbidden)
							continue
						}
						err := json.Unmarshal(m.Payload, &ev)
						if err != nil {
							s.extErc <- err
						} else {
							fmt.Println(ev)
						}
					} else if m.Type == proto.MessageTypeFullSyncRequest {
						user, ok := s.checkToken(m.AuthKey)
						if !ok || user == "" {
							s.sendError(mess.Conn(), proto.ErrForbidden)
							continue
						}
						uo, err := s.stor.GetUsersObjects(user)
						if err != nil {
							s.extErc <- err
							continue
						}
						var l []fse.FSObject
						for _, ol := range uo {
							l = append(l, fse.FSObject{
								Path:      s.ExtractOwnerFromPath(ol.Path, user),
								Name:      ol.Name,
								IsDir:     ol.IsDir,
								Hash:      ol.Hash,
								Ext:       ol.Ext,
								Size:      ol.Size,
								UpdatedAt: ol.FSUpdatedAt,
							})
						}
						err = fselink.SendSyncMessage(mess.Conn(), proto.MessageFullSyncReply{Success: true, Objects: l}, proto.MessageTypeFullSyncReply)
						if err != nil {
							s.extErc <- err
							continue
						}
					} else if m.Type == proto.MessageTypeGetFile {
						user, ok := s.checkToken(m.AuthKey)
						if !ok || user == "" {
							s.sendError(mess.Conn(), proto.ErrForbidden)
							continue
						}

						var mu proto.MessageGetFile
						err = fselink.UnpackMessage(m, proto.MessageTypeGetFile, &mu)
						if err != nil {
							s.extErc <- err
							continue
						}

						//Do not sync dirs
						if mu.Object.IsDir {
							s.sendError(mess.Conn(), proto.ErrWrongObjectType)
							continue
						}

						mu.Object.Path = strings.ReplaceAll(mu.Object.Path, "?ROOT_DIR?", "?ROOT_DIR?,"+user)

						fileName := s.fp.GetPathUnescaped(mu.Object)
						dbObj, err := s.stor.GetObject(mu.Object.Path, mu.Object.Name)
						if err != nil {
							s.extErc <- err
							continue
						}
						//Do not sync dirs (2) - if client's a smartass and still wants it somehow
						if dbObj.IsDir {
							s.sendError(mess.Conn(), proto.ErrWrongObjectType)
							continue
						}

						fileData, err := os.Open(fileName)
						if err != nil {
							s.sendError(mess.Conn(), proto.ErrFileReadingFailed)
							continue
						}

						// TLS record size can be up to 16KB but some extra bytes may apply
						// https://hpbn.co/transport-layer-security-tls/#optimize-tls-record-size
						buf := make([]byte, 15360)
						n := 0

						r := bufio.NewReader(fileData)

						for {
							n, err = r.Read(buf)
							if err == io.EOF {
								break
							}
							if err != nil {
								s.sendError(mess.Conn(), proto.ErrFileReadingFailed)
								break
							}

							err = fselink.SendSyncMessage(mess.Conn(), proto.MessageFilePart{Payload: buf[:n]}, proto.MessageTypeFileParts)
							if err != nil {
								s.extErc <- err
								break
							}
						}
						fileData.Close()

						err = fselink.SendSyncMessage(mess.Conn(), nil, proto.MessageTypeFileEnd)
						if err != nil {
							s.extErc <- err
							continue
						}

						fmt.Println("SENT FILE")
						err = mess.Conn().Close()
						if err != nil {
							s.extErc <- err
						}
						continue

					} else if m.Type == proto.MessageTypePushFile {
						user, ok := s.checkToken(m.AuthKey)
						if !ok || user == "" {
							s.sendError(mess.Conn(), proto.ErrForbidden)
							continue
						}

						var mu proto.MessagePushFile
						err = fselink.UnpackMessage(m, proto.MessageTypePushFile, &mu)
						if err != nil {
							s.sendError(mess.Conn(), proto.ErrMessageReadingFailed)
							s.extErc <- err
							continue
						}

						//Do not sync dirs
						if mu.Object.IsDir {
							s.sendError(mess.Conn(), proto.ErrWrongObjectType)
							continue
						}

						mu.Object.Path = strings.ReplaceAll(mu.Object.Path, "?ROOT_DIR?", "?ROOT_DIR?,"+user)

						//fileName := s.fp.GetPathUnescaped(mu.Object)
						dbObj, err := s.stor.GetObject(mu.Object.Path, mu.Object.Name)
						if err != nil && err != storage.ErrNotExists {
							s.sendError(mess.Conn(), proto.ErrInternalServerError)
							s.extErc <- err
							continue
						}
						//If object does not exist yet
						if dbObj.ID != 0 {
							err = s.stor.LockObject(mu.Object.Path, mu.Object.Name)
							if err != nil {
								s.extErc <- err
								s.sendError(mess.Conn(), proto.ErrInternalServerError)
								continue
							}
						}

						err = fselink.SendSyncMessage(mess.Conn(), nil, proto.MessageTypePeerReady)
						if err != nil {
							s.extErc <- err
							continue
						}

						//HERE WE SHOULD WAIT FOR FILE PARTS

						/*if dbObj.ID != 0 {
							err = s.stor.UnLockObject(mu.Object.Path, mu.Object.Name)
							if err != nil {
								s.extErc <- err
								s.sendError(mess.Conn(), proto.ErrInternalServerError)
								continue
							}
						}*/

						err = fselink.SendSyncMessage(mess.Conn(), nil, proto.MessageTypeClose)
						if err != nil {
							s.extErc <- err
							continue
						}

						fmt.Println("DOWNLOADED FILE")
						err = mess.Conn().Close()
						if err != nil {
							s.extErc <- err
						}
						continue

					} else {
						s.sendError(mess.Conn(), proto.ErrUnexpectedMessageType)
						continue
					}
				}
			}()
		}
	}

}

func (s *FSWServer) createSession(log, pwd string) bool {
	return true
}

func (s *FSWServer) checkToken(t string) (uid string, ok bool) {
	return "1", true
}

func (s *FSWServer) sendError(sm fselink.SyncMessenger, e proto.ErrorCode) {
	err := fselink.SendErrorMessage(sm, e)
	if err != nil {
		s.extErc <- err
	}
}
