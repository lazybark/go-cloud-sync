package v1

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lazybark/go-cloud-sync/pkg/fp"
	"github.com/lazybark/go-cloud-sync/pkg/fse"
	"github.com/lazybark/go-cloud-sync/pkg/fselink"
	"github.com/lazybark/go-cloud-sync/pkg/storage"
	"github.com/lazybark/go-cloud-sync/pkg/watcher"
	gts "github.com/lazybark/go-tls-server/v2/server"
)

// NewServer returns new filesystem watcher
func NewServer(stor storage.IStorage) *FSWServer {
	s := &FSWServer{}

	s.evc = make(chan (fse.FSEvent))
	s.erc = make(chan (error)) //Not used for now
	s.srvConnChan = make(chan *gts.Connection)
	s.srvErrChan = make(chan error)
	s.srvMessChan = make(chan *gts.Message)

	s.w = watcher.NewWatcher()
	s.stor = stor
	s.htsrv = fselink.NewServer()

	return s
}

// Init sets initial config to the Watcher
func (s *FSWServer) Init(root, cacheRoot, host, port, escSymbol string, evc chan (fse.FSEvent), erc chan (error)) error {
	s.conf.root = root
	s.conf.cacheRoot = cacheRoot
	s.conf.host = host
	s.conf.port = port
	s.conf.escSymbol = escSymbol
	s.extEvc = evc
	s.extErc = erc

	s.fp = fp.NewFPv1(escSymbol, root, cacheRoot)

	err := s.w.Init(root, s.extEvc, s.extErc)
	if err != nil {
		return fmt.Errorf("[INIT]: %w", err)
	}

	err = s.htsrv.Init(s.srvMessChan, s.srvConnChan, s.srvErrChan)
	if err != nil {
		return fmt.Errorf("[INIT]: %w", err)
	}

	fmt.Printf("Server started on %s:%s\n", host, port)
	fmt.Printf("Root path:%s\n", root)
	fmt.Printf("Cache path:%s\n", cacheRoot)

	return nil
}

func (s *FSWServer) ExtractOwnerFromPath(p string) string {
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
			Owner:       s.ExtractOwnerFromPath(o.Path),
			Ext:         o.Ext,
			Size:        o.Size,
			FSUpdatedAt: o.UpdatedAt,
		})
	}
	return
}

// Start launches the filesystem watcher routine. You need to call Init() before.
func (s *FSWServer) Start() error {
	s.rescanOnce()
	fmt.Println("Root directory processed")

	go s.htsrv.Listen(s.conf.host, s.conf.port)
	go s.watcherRoutine()
	s.isActive = true
	err := s.w.Start()
	if err != nil {
		return fmt.Errorf("[SERVER][START]%w", err)
	}
	return nil
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

// Stop stops the filesystem watcher and closes all channels
func (s *FSWServer) Stop() error {
	err := s.w.Stop()
	if err != nil {
		return fmt.Errorf("[STOP] can not stop watcher: %w", err)
	}
	close(s.extEvc)
	close(s.extErc)
	close(s.evc)
	close(s.erc)

	s.isActive = false
	fmt.Println("Server stopped")

	return nil
}

func (s *FSWServer) watcherRoutine() {
	fmt.Println("Waiting for connections")
	var ev fse.FSEvent
	var m fselink.ExchangeMessage
	for {
		select {
		case mess, ok := <-s.srvMessChan:
			if !ok {
				return
			}
			err := json.Unmarshal(mess.Bytes(), &m)
			if err != nil {
				s.extErc <- err
				continue
			}
			if m.Type == fselink.MessageTypeAuthReq {
				if !s.checkCredentials("", "") {
					err := fselink.SendErrorMessage(mess.Conn(), fselink.ErrCodeWrongCreds)
					if err != nil {
						s.extErc <- err
						continue
					}
					continue
				}
				err = fselink.SendSyncMessage(mess.Conn(), fselink.MessageAuthAnswer{Success: true, AuthKey: "SOMEKEY"}, fselink.MessageTypeAuthAns)
				if err != nil {
					s.extErc <- err
					continue
				}
				fmt.Println("SENT AUTH")
			} else if m.Type == fselink.MessageTypeEvent {
				if !s.checkToken(m.AuthKey, m.AuthKey) {
					err := fselink.SendErrorMessage(mess.Conn(), fselink.ErrForbidden)
					if err != nil {
						s.extErc <- err
						continue
					}
					continue
				}
				err := json.Unmarshal(m.Payload, &ev)
				if err != nil {
					s.extErc <- err
				} else {
					//s.extEvc <- ev
					fmt.Println(ev)
				}
			} else if m.Type == fselink.MessageTypeFullSyncRequest {
				if !s.checkToken(m.AuthKey, m.AuthKey) {
					err := fselink.SendErrorMessage(mess.Conn(), fselink.ErrForbidden)
					if err != nil {
						s.extErc <- err
						continue
					}
					continue
				}
				uo, err := s.stor.GetUsersObjects("1")
				if err != nil {
					s.extErc <- err
					continue
				}
				var l []fse.FSObject
				for _, ol := range uo {
					l = append(l, fse.FSObject{
						Path:      ol.Path,
						Name:      ol.Name,
						IsDir:     ol.IsDir,
						Hash:      ol.Hash,
						Ext:       ol.Ext,
						Size:      ol.Size,
						UpdatedAt: ol.FSUpdatedAt,
					})
				}
				err = fselink.SendSyncMessage(mess.Conn(), fselink.MessageFullSyncReply{Success: true, Objects: l}, fselink.MessageTypeFullSyncReply)
				if err != nil {
					s.extErc <- err
					continue
				}
			} else if m.Type == fselink.MessageTypeGetFile {
				if !s.checkToken(m.AuthKey, m.AuthKey) {
					err := fselink.SendErrorMessage(mess.Conn(), fselink.ErrForbidden)
					if err != nil {
						s.extErc <- err
						continue
					}
					continue
				}
				//user := "1"

				var mu fselink.MessageGetFile
				err = fselink.UnpackMessage(m, fselink.MessageTypeGetFile, &mu)
				if err != nil {
					s.extErc <- err
					continue
				}

				fileName := s.fp.GetPathUnescaped(mu.Object)
				stat, err := os.Stat(fileName)
				if err != nil {
					s.extErc <- err
					continue
				}
				if stat.IsDir() {
					err := fselink.SendErrorMessage(mess.Conn(), fselink.ErrWrongObjectType)
					if err != nil {
						s.extErc <- err
						continue
					}
					continue
				}

				fileData, err := os.Open(fileName)
				if err != nil {
					err := fselink.SendErrorMessage(mess.Conn(), fselink.ErrFileReadingFailed)
					if err != nil {
						s.extErc <- err
					}
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
						err := fselink.SendErrorMessage(mess.Conn(), fselink.ErrFileReadingFailed)
						if err != nil {
							s.extErc <- err
						}
						break
					}

					err = fselink.SendSyncMessage(mess.Conn(), fselink.MessageFilePart{Payload: buf[:n]}, fselink.MessageTypeFileParts)
					if err != nil {
						s.extErc <- err
						break
					}
				}
				fileData.Close()

				err = fselink.SendSyncMessage(mess.Conn(), nil, fselink.MessageTypeFileEnd)
				if err != nil {
					s.extErc <- err
					continue
				}
				fmt.Println("SENT FILE")

			} else {

			}
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
		}
	}

}

func (s *FSWServer) checkCredentials(log, pwd string) bool {
	return true
}

func (s *FSWServer) checkToken(t string, ct string) bool {
	return t == ct
}
