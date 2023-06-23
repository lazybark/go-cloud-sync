package v1

import (
	"fmt"

	"github.com/lazybark/go-cloud-sync/pkg/fp"
	"github.com/lazybark/go-cloud-sync/pkg/fse"
	"github.com/lazybark/go-cloud-sync/pkg/fselink"
	"github.com/lazybark/go-cloud-sync/pkg/storage"
	"github.com/lazybark/go-cloud-sync/pkg/watcher"
)

type FSWClient struct {
	//w is an instance of filesystem watcher to recieve events from
	w watcher.IFilesystemWatcher

	//extEvc is the channel to send notifications to external routine
	extEvChannel chan (fse.FSEvent)

	//evc is the channel to recieve events from internal routine
	evc chan (fse.FSEvent)

	//extErc is the channel to send errors to external routine
	extErc chan (error)
	//erc is the channel to recieve errors from internal routine
	erc chan (error)

	//root is the full path to directory where Watcher will watch for events (subdirs included)
	root string

	//stor is the storage watcher uses to store & process filesystem hashes
	fp fp.Fileprocessor

	link fselink.FSEClientLink

	db storage.IStorage
}

func NewClient(db storage.IStorage, root string) *FSWClient {
	s := &FSWClient{
		db:   db,
		root: root,
	}
	s.evc = make(chan (fse.FSEvent))
	s.erc = make(chan (error))

	s.w = watcher.NewWatcher()
	s.fp = fp.NewFPv1(",", root)
	s.link = fselink.NewClient()

	return s
}

// Init sets initial config to the Watcher
func (s *FSWClient) Init(evc chan (fse.FSEvent), erc chan (error)) error {
	s.extEvChannel = evc
	s.extErc = erc

	err := s.w.Init(s.root, s.evc, s.erc)
	if err != nil {
		return fmt.Errorf("[FSWATCHER][INIT] can not start watcher: %w", err)
	}

	err = s.link.Init(5555, "localhost", "log", "pass")
	if err != nil {
		return fmt.Errorf("[FSWATCHER][INIT] can not connect to server: %w", err)
	}

	return nil
}

// Start launches the filesystem watcher routine. You need to call Init() before.
func (s *FSWClient) Start() error {
	go s.watcherRoutine()
	//go s.rescanRoutine()
	return s.w.Start()
}

// Stop stops the filesystem watcher and closes all channels
func (s *FSWClient) Stop() error {
	err := s.w.Stop()
	if err != nil {
		return fmt.Errorf("[FSWATCHER][STOP] can not stop watcher: %w", err)
	}
	close(s.extEvChannel)
	close(s.extErc)
	close(s.evc)
	close(s.erc)

	return nil
}

func (s *FSWClient) watcherRoutine() {
	done := make(chan bool)
	go func() {
		defer close(done)

		for {
			select {
			case event, ok := <-s.evc:
				if !ok {
					return
				}

				//Process event with storage
				if event.Action == fse.Create || event.Action == fse.Write {
					obj, err := s.fp.ProcessObject(event.Object, true)
					if err != nil {
						s.extErc <- fmt.Errorf("[FSWATCHER][WATCH][%s] processing error: %w", event.Action.String(), err)
					}
					//We skip objects that have no access to avoid syncing temp junks of apps
					if obj.Hash == "" && !obj.IsDir {
						continue
					}
					err = s.db.AddOrUpdateObject(obj)
					if err != nil {
						s.extErc <- fmt.Errorf("[FSWATCHER][WATCH][%s] processing error: %w", event.Action.String(), err)
					}
					if event.Action == fse.Create && obj.IsDir {
						err = s.w.Add(event.Object.Path)
						if err != nil {
							s.extErc <- fmt.Errorf("[FSWATCHER][WATCH][%s] adding watcher failed: %w", event.Action.String(), err)
						}
					}
				} else if event.Action == fse.Remove {
					obj := event.Object
					path, name, err := s.fp.ConvertPathName(obj)
					if err != nil {
						s.extErc <- fmt.Errorf("[FSWATCHER][WATCH][ConvertPathName] processing error: %w", err)
					}
					obj.Path = path
					obj.Name = name

					err = s.db.RemoveObject(obj, true)
					if err != nil {
						s.extErc <- fmt.Errorf("[FSWATCHER][WATCH][RemoveObject] processing error: %w", err)
					}
				}
				//Sent event to external code
				s.extEvChannel <- event

				//Now send to the server
				err := s.link.SendEvent(event)
				if err != nil {
					s.extErc <- fmt.Errorf("[FSWATCHER][WATCH][%s] error notifying server: %w", event.Action.String(), err)
				}

			case err, ok := <-s.erc:
				if !ok {
					return
				}
				s.extErc <- fmt.Errorf("[FSWATCHER][WATCH] fs watcher error: %w", err)
			}
		}
	}()

	<-done
}

func (s *FSWClient) Add(dir string) error {
	err := s.w.Add(dir)
	if err != nil {
		s.erc <- fmt.Errorf("[FSWATCHER][WATCH] fs watcher add failed: %w", err)
	}
	return nil
}

/*
// rescanRoutine checks rescanBuffer and rescans every object once per a minute.
// Object that was scanned is being deleted from the buffer.
func (c *FSWClient) rescanRoutine() {

}
*/
