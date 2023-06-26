package v1

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/lazybark/go-cloud-sync/pkg/fse"
	"github.com/lazybark/go-cloud-sync/pkg/fselink"
	"github.com/lazybark/go-cloud-sync/pkg/storage"
	"github.com/lazybark/go-cloud-sync/pkg/watcher"
	gts "github.com/lazybark/go-tls-server/v2/server"
)

type FSWServer struct {
	//w is an instance of filesystem watcher to recieve events from
	w watcher.IFilesystemWatcher

	//extEvc is the channel to send notifications to external routine
	extEvc chan (fse.FSEvent)

	//evc is the channel to recieve events from internal routine
	evc chan (fse.FSEvent)

	//extErc is the channel to send errors to external routine
	extErc chan (error)
	//erc is the channel to recieve errors from internal routine
	erc chan (error)

	srvMessChan chan *gts.Message
	srvErrChan  chan error
	srvConnChan chan *gts.Connection

	//root is the full path to directory where Watcher will watch for events (subdirs included)
	root string

	//stor is the storage watcher uses to store & process filesystem hashes
	stor storage.IStorage

	htsrv fselink.FSEServerPool

	isActive bool
}

// NewServer returns new filesystem watcher
func NewServer(stor storage.IStorage) *FSWServer {
	s := &FSWServer{}
	s.evc = make(chan (fse.FSEvent))
	s.erc = make(chan (error))

	s.w = watcher.NewWatcher()
	s.stor = stor
	s.htsrv = fselink.NewServer()
	s.srvConnChan = make(chan *gts.Connection)
	s.srvErrChan = make(chan error)
	s.srvMessChan = make(chan *gts.Message)

	return s
}

// Init sets initial config to the Watcher
func (s *FSWServer) Init(root string, evc chan (fse.FSEvent), erc chan (error)) error {
	s.root = root
	s.extEvc = evc
	s.extErc = erc

	err := s.w.Init(root, s.evc, s.erc)
	if err != nil {
		return fmt.Errorf("[FSWATCHER][INIT] can not init watcher: %w", err)
	}

	err = s.htsrv.Init(s.srvMessChan, s.srvConnChan, s.srvErrChan)
	if err != nil {
		return fmt.Errorf("[FSWATCHER][INIT] can not init server: %w", err)
	}

	return nil
}

// Start launches the filesystem watcher routine. You need to call Init() before.
func (s *FSWServer) Start() error {
	go s.htsrv.Listen("localhost", "5555")
	go s.watcherRoutine()
	s.isActive = true
	return s.w.Start()
}

// Stop stops the filesystem watcher and closes all channels
func (s *FSWServer) Stop() error {
	err := s.w.Stop()
	if err != nil {
		return fmt.Errorf("[FSWATCHER][STOP] can not stop watcher: %w", err)
	}
	close(s.extEvc)
	close(s.extErc)
	close(s.evc)
	close(s.erc)

	s.isActive = false

	return nil
}

func (s *FSWServer) watcherRoutine() {
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
				l := []fse.FSObject{
					{Path: `?ROOT_DIR?`, Name: "New folder (4)", Hash: "", UpdatedAt: time.Now().Add(time.Minute * +14), IsDir: true},
					{Path: `?ROOT_DIR?`, Name: "SOME WEIRD FOLDER", Hash: "", UpdatedAt: time.Now().Add(time.Minute * +14), IsDir: true},
					{Path: `?ROOT_DIR?,SOME WEIRD FOLDER`, Name: "file.jpg", Hash: "asdfghfdsfgh", UpdatedAt: time.Now().Add(time.Minute * +14), IsDir: false},
				}
				err = fselink.SendSyncMessage(mess.Conn(), fselink.MessageFullSyncReply{Success: true, Objects: l}, fselink.MessageTypeFullSyncReply)
				if err != nil {
					s.extErc <- err
					continue
				}
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
			fmt.Println(connection)

		}
	}

}

func (s *FSWServer) checkCredentials(log, pwd string) bool {
	return true
}

func (s *FSWServer) checkToken(t string, ct string) bool {
	return t == ct
}
