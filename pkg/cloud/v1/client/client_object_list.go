package client

import (
	"fmt"

	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/proto"
)

func (c *FSWClient) GetServerObjList() (l []proto.FSObject, err error) {
	err = c.link.SendSyncMessage(nil, proto.MessageTypeFullSyncRequest)
	if err != nil {
		c.extErc <- fmt.Errorf("[GetServerObjList]: %w", err)
		return
	}
	maa, err := c.link.Await()
	if err != nil {
		c.extErc <- fmt.Errorf("[GetServerObjList]%w", err)
		return
	}
	if maa.Type == proto.MessageTypeError {
		a, err := maa.ReadError()
		if err != nil {
			return l, fmt.Errorf("[GetServerObjList] %w", err)
		}
		return l, fmt.Errorf("sync error #%d: %s", a.ErrorCode, a.Error)

	} else if maa.Type == proto.MessageTypeFullSyncReply {
		a, err := maa.ReadFullSyncReply()
		if err != nil {
			return l, fmt.Errorf("[GetServerObjList] %w", err)
		}
		l = a.Objects

		return l, err
	}

	return l, fmt.Errorf("[GetObjList] unexpected answer type '%s'", maa.Type)
}

func (c *FSWClient) GetLocalObjects() (objs []proto.FSObject, err error) {
	objs, err = c.fp.ProcessDirectory(c.cfg.Root)
	if err != nil {
		err = fmt.Errorf("[DiffListWithServer]: %w", err)
		return
	}
	return
}
