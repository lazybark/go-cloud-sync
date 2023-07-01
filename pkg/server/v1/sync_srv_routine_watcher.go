package v1

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/lazybark/go-cloud-sync/pkg/storage"
	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/proto"
)

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
						a, err := m.ReadObjectData()
						if err != nil {
							s.extErc <- err
							continue
						}

						//Do not SEND dirs
						//Client should create dir after full sync request
						if a.Object.IsDir {
							c.SendError(proto.ErrWrongObjectType)
							continue
						}

						a.Object.Path = strings.ReplaceAll(a.Object.Path, "?ROOT_DIR?", "?ROOT_DIR?,"+c.uid)

						fileName := s.fp.GetPathUnescaped(a.Object)
						dbObj, err := s.stor.GetObject(a.Object.Path, a.Object.Name)
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
						a, err := m.ReadObjectData()
						if err != nil {
							c.SendError(proto.ErrMessageReadingFailed)
							s.extErc <- err
							continue
						}

						a.Object.Path = strings.ReplaceAll(a.Object.Path, "?ROOT_DIR?", "?ROOT_DIR?,"+c.uid)
						pathUnescaped := filepath.Join(s.fp.UnescapePath(a.Object))
						pathFullUnescaped := filepath.Join(s.fp.GetPathUnescaped(a.Object))

						//fileName := s.fp.GetPathUnescaped(mu.Object)
						dbObj, err := s.stor.GetObject(a.Object.Path, a.Object.Name)
						if err != nil && err != storage.ErrNotExists {
							c.SendError(proto.ErrInternalServerError)
							s.extErc <- err
							continue
						}
						//If object does not exist yet
						if dbObj.ID != 0 {
							err = s.stor.LockObject(a.Object.Path, a.Object.Name)
							if err != nil {
								s.extErc <- err
								c.SendError(proto.ErrInternalServerError)
								continue
							}
							//Do not sync dirs
							if a.Object.IsDir {
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
						if a.Object.IsDir {
							fmt.Println("CREATING DIR")
							//CREATE DIR HERE
							if err := os.MkdirAll(pathFullUnescaped, os.ModePerm); err != nil {
								s.extErc <- err
								c.SendError(proto.ErrInternalServerError)
								return
							}
							err = s.stor.AddOrUpdateObject(storage.FSObjectStored{
								Name:        a.Object.Name,
								Path:        a.Object.Path,
								Owner:       c.uid,
								Hash:        a.Object.Hash,
								IsDir:       a.Object.IsDir,
								Ext:         a.Object.Ext,
								Size:        a.Object.Size,
								FSUpdatedAt: a.Object.UpdatedAt,
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
								a, err := m.ReadError()
								if err != nil {
									s.extErc <- err
									return
								}
								s.extErc <- fmt.Errorf("sync error #%d: %s", a.ErrorCode, a.Error)
								return
							} else if m.Type == proto.MessageTypeFileParts {
								a, err := m.ReadFilePart()
								if err != nil {
									s.extErc <- err
									return
								}
								_, err = destFile.Write(a.Payload)
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
						err = os.Chtimes(pathFullUnescaped, a.Object.UpdatedAt, a.Object.UpdatedAt)
						if err != nil {
							s.extErc <- err
							c.SendError(proto.ErrInternalServerError)
							continue
						}

						err = s.stor.AddOrUpdateObject(storage.FSObjectStored{
							Name:        a.Object.Name,
							Path:        a.Object.Path,
							Owner:       c.uid,
							Hash:        a.Object.Hash,
							IsDir:       a.Object.IsDir,
							Ext:         a.Object.Ext,
							Size:        a.Object.Size,
							FSUpdatedAt: a.Object.UpdatedAt,
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
							err = s.stor.UnLockObject(a.Object.Path, a.Object.Name)
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
