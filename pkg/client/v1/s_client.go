package v1

import (
	"fmt"
	"log"
	"time"

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

	//stor is the storage watcher uses to store & process filesystem hashes
	fp fp.Fileprocessor

	link fselink.FSEClientLink

	db storage.IStorage

	cfg ClientConfig
}

type ClientConfig struct {
	//root is the full path to directory where Watcher will watch for events (subdirs included)
	Root string
}

func NewClient(db storage.IStorage, cacheDir, root string) *FSWClient {
	c := &FSWClient{
		db:  db,
		cfg: ClientConfig{Root: root},
	}
	c.evc = make(chan (fse.FSEvent))
	c.erc = make(chan (error))

	c.w = watcher.NewWatcher()
	c.fp = fp.NewFPv1(",", root)
	link, err := fselink.NewClient(cacheDir, root)
	if err != nil {
		log.Fatal(err)
	}
	c.link = link

	return c
}

// Init sets initial config to the Watcher
func (c *FSWClient) Init(evc chan (fse.FSEvent), erc chan (error)) error {
	c.extEvChannel = evc
	c.extErc = erc

	err := c.w.Init(c.cfg.Root, c.evc, c.erc)
	if err != nil {
		return fmt.Errorf("[FSWATCHER][INIT] can not start watcher: %w", err)
	}

	err = c.link.Init(5555, "localhost", "log", "pass")
	if err != nil {
		return fmt.Errorf("[FSWATCHER][INIT] can not connect to server: %w", err)
	}

	return nil
}

// Start launches the filesystem watcher routine. You need to call Init() before.
func (c *FSWClient) Start() error {
	go c.watcherRoutine()
	c.rescanOnce()
	/*local, err := c.GetLocalObjects()
	if err != nil {
		return fmt.Errorf("[FSWATCHER][START] local objects: %w", err)
	}
	/*err = c.db.RefillDatabase(local)
	if err != nil {
		return fmt.Errorf("[FSWATCHER][START] refill db: %w", err)
	}*/

	/*objsOnServer, err := c.GetServerObjList()
	if err != nil {
		return fmt.Errorf("[FSWATCHER][START][DiffListWithServer]: %w", err)
	}

	download, created, updated, err := c.GetDiffListWithServer(local, objsOnServer)
	if err != nil {
		return fmt.Errorf("[FSWATCHER][START] diff with server: %w", err)
	}
	for _, o := range download {
		go c.DownloadObject(o)
	}
	//Each obj in created & updated is treated as a new FS event
	for _, o := range created {
		o.Path = c.fp.GetPathUnescaped(o)
		go c.addFsEvent(fse.FSEvent{Action: fse.Create, Object: o})
	}
	for _, o := range updated {
		o.Path = c.fp.GetPathUnescaped(o)
		go c.addFsEvent(fse.FSEvent{Action: fse.Write, Object: o})
	}*/

	err := c.w.Start()
	if err != nil {
		return fmt.Errorf("[FSWATCHER][START] starting watcher: %w", err)
	}
	go c.rescanRoutine(time.Minute)

	return nil
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

func (s *FSWClient) Add(dir string) error {
	err := s.w.Add(dir)
	if err != nil {
		s.erc <- fmt.Errorf("[FSWATCHER][WATCH] fs watcher add failed: %w", err)
	}
	return nil
}
