package client

import (
	"fmt"
	"path/filepath"

	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/proto"
)

func (c *FSWClient) DownloadObject(obj proto.FSObject) {
	pathFullUnescaped := filepath.Join(c.fp.GetPathUnescaped(obj))

	if c.IsInActionBuffer(pathFullUnescaped) {
		return
	}
	fmt.Println("DOWNLOADING:", pathFullUnescaped)

	c.AddToActionBuffer(pathFullUnescaped)
	defer c.RemoveFromActionBuffer(pathFullUnescaped)

	//Create file in cache
	file, err := c.fp.NewEmptyCache(obj)
	if err != nil {
		c.extErc <- fmt.Errorf("[DownloadObject]: %w", err)
		return
	}

	newConn, err := c.link.NewConnectionInSession()
	if err != nil {
		c.extErc <- fmt.Errorf("[DownloadObject]: %w", err)
		return
	}

	err = newConn.SendSyncMessage(proto.MessageObject{Object: obj}, proto.MessageTypeGetFile)
	if err != nil {
		c.extErc <- fmt.Errorf("[DownloadObject]: %w", err)
		return
	}

	var maa proto.ExchangeMessage
	for {
		maa, err = newConn.Await()
		if err != nil {
			c.extErc <- fmt.Errorf("[DownloadObject]: %w", err)
			return
		}
		if maa.Type == proto.MessageTypeError {
			a, err := maa.ReadError()
			if err != nil {
				c.extErc <- fmt.Errorf("[DownloadObject]: %w", err)
				return
			}
			c.extErc <- fmt.Errorf("sync error #%d: %s", a.ErrorCode, a.Error)
			return
		} else if maa.Type == proto.MessageTypeFileParts {
			a, err := maa.ReadFilePart()
			if err != nil {
				c.extErc <- fmt.Errorf("[DownloadObject]%w", err)
				return
			}
			_, err = file.Write(a.Payload)
			if err != nil {
				c.extErc <- fmt.Errorf("[DownloadObject]%w", err)
				return
			}

		} else if maa.Type == proto.MessageTypeFileEnd {
			break
		} else {
			c.extErc <- fmt.Errorf("[DownloadObject] unexpected answer type '%s'", maa.Type)
			return
		}
	}

	err = file.Close()
	if err != nil {
		c.extErc <- err
		return
	}

	err = c.fp.ReplaceFromCache(file)
	if err != nil {
		c.extErc <- err
		err = file.Remove()
		if err != nil {
			c.extErc <- err
		}
		return
	}
}

func (c *FSWClient) IsInActionBuffer(object string) bool {
	c.ActionsBufferMutex.Lock()
	_, yes := c.ActionsBuffer[object]
	c.ActionsBufferMutex.Unlock()
	return yes
}

func (c *FSWClient) AddToActionBuffer(object string) {
	c.ActionsBufferMutex.Lock()
	c.ActionsBuffer[object] = true
	c.ActionsBufferMutex.Unlock()
}

func (c *FSWClient) RemoveFromActionBuffer(object string) {
	c.ActionsBufferMutex.Lock()
	delete(c.ActionsBuffer, object)
	c.ActionsBufferMutex.Unlock()
}
