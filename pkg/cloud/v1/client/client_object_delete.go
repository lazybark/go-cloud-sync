package client

import (
	"fmt"
	"path/filepath"

	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/proto"
)

func (c *FSWClient) DeleteObject(obj proto.FSObject) {
	pathFullUnescaped := filepath.Join(c.fp.GetPathUnescaped(obj))

	if c.IsInActionBuffer(pathFullUnescaped) {
		return
	}

	c.AddToActionBuffer(pathFullUnescaped)
	defer c.RemoveFromActionBuffer(pathFullUnescaped)

	newConn, err := c.link.NewConnectionInSession()
	if err != nil {
		c.extErc <- fmt.Errorf("[DeleteObject]: %w", err)
		return
	}

	err = newConn.SendSyncMessage(proto.MessageObject{Object: obj}, proto.MessageTypeDeleteObject)
	if err != nil {
		c.extErc <- fmt.Errorf("[DeleteObject]: %w", err)
		return
	}
	maa, err := newConn.Await()
	if err != nil {
		c.extErc <- fmt.Errorf("[DeleteObject]%w", err)
		return
	}
	if maa.Type == proto.MessageTypeError {
		a, err := maa.ReadError()
		if err != nil {
			c.extErc <- fmt.Errorf("[DeleteObject]%w", err)
		}
		c.extErc <- fmt.Errorf("sync error #%d: %s", a.ErrorCode, a.Error)
	} else if maa.Type == proto.MessageTypeClose {
		err = newConn.Close()
		if err != nil {
			c.extErc <- fmt.Errorf("[DeleteObject]%w", err)
		}
	} else {
		c.extErc <- fmt.Errorf("[DeleteObject] unexpected answer type '%s'", maa.Type)
		return
	}
}
