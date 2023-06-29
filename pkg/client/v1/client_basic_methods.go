package v1

import (
	"fmt"
	"time"

	"github.com/lazybark/go-cloud-sync/pkg/fse"
)

func (c *FSWClient) Init(evc chan (fse.FSEvent), erc chan (error)) error {
	c.extEvChannel = evc
	c.extErc = erc
	c.ActionsBuffer = make(map[string]bool)

	err := c.w.Init(c.cfg.Root, c.evc, c.erc)
	if err != nil {
		return fmt.Errorf("[SYNC][INIT] can not start watcher: %w", err)
	}

	err = c.link.Init(5555, "localhost", "log", "pass")
	if err != nil {
		return fmt.Errorf("[SYNC][INIT] can not connect to server: %w", err)
	}

	fmt.Printf("Client Init on %s\n", c.cfg.Root)

	return nil
}

func (c *FSWClient) Start() error {
	go c.watcherRoutine()
	c.rescanOnce()

	err := c.w.Start()
	if err != nil {
		return fmt.Errorf("[SYNC][START] starting watcher: %w", err)
	}
	go c.rescanRoutine(time.Minute)
	fmt.Println("Client started")

	return nil
}

// Stop stops the filesystem watcher and closes all channels
func (s *FSWClient) Stop() error {
	err := s.w.Stop()
	if err != nil {
		return fmt.Errorf("[SYNC][STOP] can not stop watcher: %w", err)
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
		s.erc <- fmt.Errorf("[SYNC][WATCH] fs watcher add failed: %w", err)
	}
	return nil
}