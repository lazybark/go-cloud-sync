package v1

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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

func (s *FSWServer) watcherRoutine() {
	fmt.Println("Waiting for connections")
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
			//Add connection to pool to be able to control it after
			c := syncConnection{tlsConnection: connection}
			s.addToPool(&c)

			go func() {
				for mess := range connection.MessageChan {
					//GET MESSAGE HEAD
					err := json.Unmarshal(mess.Bytes(), &m)
					if err != nil {
						s.extErc <- err
						continue
					}

					//SECURITY CHECKS
					if m.Type != proto.MessageTypeAuthReq {
						if m.AuthKey == "" {
							c.SendError(proto.ErrForbidden)
							continue
						}
						ok, err := s.checkToken(c.clientTokenHash, m.AuthKey)
						if err != nil {
							c.SendError(proto.ErrInternalServerError)
							s.extErc <- err
							continue
						}
						if !ok {
							c.SendError(proto.ErrForbidden)
							continue
						}
					}

					//MESSAGE PROCESSING
					if m.Type == proto.MessageTypeAuthReq {

						s.processAuth("", "", &c)

					} else if m.Type == proto.MessageTypeFullSyncRequest {

						s.processFullSyncRequest(&c)

					} else if m.Type == proto.MessageTypeGetFile {
						var mu proto.MessageObject
						err = fselink.UnpackMessage(m, proto.MessageTypeGetFile, &mu)
						if err != nil {
							s.extErc <- err
							continue
						}

						//Do not SEND dirs
						//Client should create dir after full sync request
						if mu.Object.IsDir {
							c.SendError(proto.ErrWrongObjectType)
							continue
						}

						mu.Object.Path = strings.ReplaceAll(mu.Object.Path, "?ROOT_DIR?", "?ROOT_DIR?,"+c.uid)

						fileName := s.fp.GetPathUnescaped(mu.Object)
						dbObj, err := s.stor.GetObject(mu.Object.Path, mu.Object.Name)
						if err != nil {
							s.extErc <- err
							continue
						}
						//Do not SEND dirs (2) - if client's a smartass and still wants it somehow
						if dbObj.IsDir {
							c.SendError(proto.ErrWrongObjectType)
							continue
						}

						fileData, err := os.Open(fileName)
						if err != nil {
							c.SendError(proto.ErrFileReadingFailed)
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
								c.SendError(proto.ErrFileReadingFailed)
								break
							}

							err = c.SendMessage(proto.MessageFilePart{Payload: buf[:n]}, proto.MessageTypeFileParts)
							if err != nil {
								s.extErc <- err
								break
							}
						}
						fileData.Close()

						err = c.SendMessage(nil, proto.MessageTypeFileEnd)
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
						var mu proto.MessageObject
						err = fselink.UnpackMessage(m, proto.MessageTypePushFile, &mu)
						if err != nil {
							c.SendError(proto.ErrMessageReadingFailed)
							s.extErc <- err
							continue
						}

						mu.Object.Path = strings.ReplaceAll(mu.Object.Path, "?ROOT_DIR?", "?ROOT_DIR?,"+c.uid)
						pathUnescaped := filepath.Join(s.fp.UnescapePath(mu.Object))
						pathFullUnescaped := filepath.Join(s.fp.GetPathUnescaped(mu.Object))

						//fileName := s.fp.GetPathUnescaped(mu.Object)
						dbObj, err := s.stor.GetObject(mu.Object.Path, mu.Object.Name)
						if err != nil && err != storage.ErrNotExists {
							c.SendError(proto.ErrInternalServerError)
							s.extErc <- err
							continue
						}
						//If object does not exist yet
						if dbObj.ID != 0 {
							err = s.stor.LockObject(mu.Object.Path, mu.Object.Name)
							if err != nil {
								s.extErc <- err
								c.SendError(proto.ErrInternalServerError)
								continue
							}
							//Do not sync dirs
							if mu.Object.IsDir {
								err = c.SendMessage(nil, proto.MessageTypeClose)
								if err != nil {
									s.extErc <- err
									continue
								}
								err = mess.Conn().Close()
								if err != nil {
									s.extErc <- err
									continue
								}
								continue
							}
						}
						if mu.Object.IsDir {
							fmt.Println("CREATING DIR")
							//CREATE DIR HERE
							if err := os.MkdirAll(pathFullUnescaped, os.ModePerm); err != nil {
								s.extErc <- err
								c.SendError(proto.ErrInternalServerError)
								return
							}
							err = s.stor.AddOrUpdateObject(storage.FSObjectStored{
								Name:        mu.Object.Name,
								Path:        mu.Object.Path,
								Owner:       c.uid,
								Hash:        mu.Object.Hash,
								IsDir:       mu.Object.IsDir,
								Ext:         mu.Object.Ext,
								Size:        mu.Object.Size,
								FSUpdatedAt: mu.Object.UpdatedAt,
							})
							if err != nil {
								s.extErc <- err
								c.SendError(proto.ErrInternalServerError)
								continue
							}

							err = c.SendMessage(nil, proto.MessageTypeClose)
							if err != nil {
								s.extErc <- err
								continue
							}
							err = mess.Conn().Close()
							if err != nil {
								s.extErc <- err
								continue
							}
							continue
						}

						err = c.SendMessage(nil, proto.MessageTypePeerReady)
						if err != nil {
							s.extErc <- err
							continue
						}

						destFile, err := s.fp.CreateFileInCache()
						if err != nil {
							s.extErc <- err
							c.SendError(proto.ErrInternalServerError)
							return
						}
						//HERE WE SHOULD WAIT FOR FILE PARTS
						for filePartMessage := range connection.MessageChan {
							err := json.Unmarshal(filePartMessage.Bytes(), &m)
							if err != nil {
								s.extErc <- err
								continue
							}
							if m.Type == proto.MessageTypeError {
								var se proto.MessageError
								err := fselink.UnpackMessage(m, proto.MessageTypeError, &se)
								if err != nil {
									s.extErc <- err
									return
								}
								s.extErc <- err
								return
							} else if m.Type == proto.MessageTypeFileParts {
								var mp proto.MessageFilePart
								err := fselink.UnpackMessage(m, proto.MessageTypeFileParts, &mp)
								if err != nil {
									s.extErc <- err
									return
								}
								_, err = destFile.Write(mp.Payload)
								if err != nil {
									s.extErc <- err
									return
								}

							} else if m.Type == proto.MessageTypeFileEnd {
								break
							} else {
								s.extErc <- fmt.Errorf("[DownloadObject] unexpected answer type '%s'", m.Type)
								return
							}
						}

						destFile.Close()

						if err := os.MkdirAll(pathUnescaped, os.ModePerm); err != nil {
							s.extErc <- err
							err = s.fp.DeleteFileInCache(destFile.Name())
							if err != nil {
								s.extErc <- err
							}
							return
						}
						err = os.Rename(destFile.Name(), pathFullUnescaped)
						if err != nil {
							s.extErc <- err
							err = s.fp.DeleteFileInCache(destFile.Name())
							if err != nil {
								s.extErc <- err
							}
							return
						}
						err = os.Chtimes(pathFullUnescaped, mu.Object.UpdatedAt, mu.Object.UpdatedAt)
						if err != nil {
							s.extErc <- err
							c.SendError(proto.ErrInternalServerError)
							continue
						}

						err = s.stor.AddOrUpdateObject(storage.FSObjectStored{
							Name:        mu.Object.Name,
							Path:        mu.Object.Path,
							Owner:       c.uid,
							Hash:        mu.Object.Hash,
							IsDir:       mu.Object.IsDir,
							Ext:         mu.Object.Ext,
							Size:        mu.Object.Size,
							FSUpdatedAt: mu.Object.UpdatedAt,
						})
						if err != nil {
							s.extErc <- err
							c.SendError(proto.ErrInternalServerError)
							continue
						}

						err = c.SendMessage(nil, proto.MessageTypeClose)
						if err != nil {
							s.extErc <- err
							c.SendError(proto.ErrInternalServerError)
							continue
						}

						err = mess.Conn().Close()
						if err != nil {
							s.extErc <- err
							continue
						}
						if dbObj.ID != 0 {
							err = s.stor.UnLockObject(mu.Object.Path, mu.Object.Name)
							if err != nil {
								s.extErc <- err
								c.SendError(proto.ErrInternalServerError)
								continue
							}
						}
						fmt.Println("DOWNLOADED FILE")
						continue

					} else if m.Type == proto.MessageTypeDeleteObject {
						s.processDelete(&c, m)
					} else {
						c.SendError(proto.ErrUnexpectedMessageType)
						continue
					}
				}
			}()
		}
	}

}

func (s *FSWServer) createSession(log, pwd string) (user, sessionKey string) {
	return "1", "AUTH_KEY"
}

func (s *FSWServer) checkToken(hash, token string) (ok bool, err error) {
	return comparePasswords(hash, token)
}
