package client

import (
	"fmt"
	"path/filepath"

	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/proto"
)

func (c *FSWClient) Init(evc chan (proto.FSEvent), erc chan (error), login, pwd string) error {
	c.extEvChannel = evc
	c.extErc = erc
	c.ActionsBuffer = make(map[string]bool)

	fmt.Println("Starting filesystem watcher")
	err := c.w.Init(c.cfg.Root, c.evc, c.erc)
	if err != nil {
		return fmt.Errorf("[SYNC][INIT] can not start watcher: %w", err)
	}

	fmt.Println("Connecting to server")
	err = c.link.Init(5555, "localhost", login, pwd)
	if err != nil {
		return fmt.Errorf("[SYNC][INIT] can not connect to server: %w", err)
	}

	fmt.Printf("Client Init on %s\n", c.cfg.Root)

	return nil
}

func (c *FSWClient) Start() error {
	go c.watcherRoutine()
	c.resyncOnce()

	err := c.w.Start()
	if err != nil {
		return fmt.Errorf("[SYNC][START] starting watcher: %w", err)
	}
	//go c.resyncRoutine(time.Minute)
	err = c.link.SendSyncMessage(nil, proto.MessageTypeSyncStart)
	if err != nil {
		return fmt.Errorf("[SYNC][START]%w", err)
	}
	fmt.Println("Client started")

	var maa proto.ExchangeMessage
	var pathFullUnescaped string
	for {
		maa, err = c.link.Await()
		if err != nil {
			return fmt.Errorf("[SYNC][START]%w", err)
		}
		if maa.Type == proto.MessageTypeError {
			a, err := maa.ReadError()
			if err != nil {
				c.extErc <- fmt.Errorf("[SYNC][START]%w", err)
			}
			c.extErc <- fmt.Errorf("sync error #%d: %s", a.ErrorCode, a.Error)
		} else if maa.Type == proto.MessageTypeSyncEvent {
			a, err := maa.ReadSyncEvent()
			if err != nil {
				c.extErc <- fmt.Errorf("[SYNC][START]%w", err)
			}

			pathFullUnescaped = filepath.Join(c.fp.GetPathUnescaped(a.Object))
			if c.IsInActionBuffer(pathFullUnescaped) {
				continue
			}

			fmt.Println(a.Object, a.Event)
			if a.Event == proto.Remove {
				c.AddToActionBuffer(pathFullUnescaped)
				err = c.fp.Remove(a.Object)
				if err != nil {
					c.extErc <- fmt.Errorf("[SYNC][START]%w", err)
				}
				c.RemoveFromActionBuffer(pathFullUnescaped)
			} else if a.Event == proto.Create {
				go c.DownloadObject(a.Object)
			} else if a.Event == proto.Write {
				go c.DownloadObject(a.Object)
			}
		} else {
			c.extErc <- fmt.Errorf("[SYNC]unexpected answer type '%s'", maa.Type)
		}
	}

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
