package v1

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/lazybark/go-cloud-sync/pkg/fse"
)

type FSWatcher struct {
	//Instance of filesystem watcher to recieve events from
	w    *fsnotify.Watcher
	evc  chan (fse.FSEvent)
	erc  chan (error)
	root string
}

func NewWatcher() *FSWatcher {
	return &FSWatcher{}
}

func (fw *FSWatcher) Init(root string, evc chan (fse.FSEvent), erc chan (error)) error {
	fw.root = root
	fw.evc = evc
	fw.erc = erc

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("[FSWATCHER][INIT] can not init watcher: %w", err)
	}
	fw.w = watcher

	return nil
}

func (fw *FSWatcher) Start() error {
	if err := os.MkdirAll(fw.root, os.ModePerm); err != nil {
		return fmt.Errorf("[FSWATCHER][START] can not make dir: %w", err)
	}

	go fw.filesystemWatcherRoutine()
	return nil
}

// Stop stops the filesystem watcher and closes all channels
func (fw *FSWatcher) Stop() error {
	err := fw.w.Close()
	if err != nil {
		return fmt.Errorf("[FSWATCHER][STOP] can not stop watcher: %w", err)
	}
	close(fw.evc)
	close(fw.erc)

	return nil
}

// ConvertFSNotifyEventToFSEvent converts events from fsnotify.Event into fse.FSEvent.
// It fills only Action type and Object.Path. Other fields are left empty and meant to be filled
// by reciever routine
func ConvertFSNotifyEventToFSEvent(event fsnotify.Event) fse.FSEvent {
	e := fse.FSEvent{}

	if event.Op == fsnotify.Create {
		e.Action = fse.Create
	} else if event.Op == fsnotify.Write {
		e.Action = fse.Write
	} else if event.Op == fsnotify.Remove {
		e.Action = fse.Remove
	} else if event.Op == fsnotify.Rename {
		e.Action = fse.Rename
	} else if event.Op == fsnotify.Chmod {
		e.Action = fse.Chmod
	} else {
		e.Action = fse.NoAction
	}

	e.Object.Path = event.Name

	return e
}

func (fw *FSWatcher) filesystemWatcherRoutine() {
	w := fw.w
	done := make(chan bool)
	go func() {
		defer close(done)

		for {
			select {
			case event, ok := <-w.Events:
				if !ok {
					return
				}
				fw.evc <- ConvertFSNotifyEventToFSEvent(event)

			case err, ok := <-w.Errors:
				if !ok {
					return
				}
				fw.erc <- fmt.Errorf("[FSWATCHER][WATCH] fs watcher error: %w", err)
			}
		}
	}()

	err := w.Add(fw.root)
	if err != nil {
		fw.erc <- fmt.Errorf("[FSWATCHER][WATCH] fs watcher add failed: %w", err)
	}
	<-done
}
