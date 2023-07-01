package client

import (
	"fmt"

	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/proto"
)

func (sc *LinkClient) DeleteObject(obj proto.FSObject) (err error) {
	link, err := NewClient()
	if err != nil {
		return fmt.Errorf("[DeleteObject]%w", err)
	}

	err = link.Init(sc.serverPort, sc.serverAddr, sc.login, sc.pwd)
	if err != nil {
		return fmt.Errorf("[DeleteObject]%w", err)
	}
	defer link.c.Close()

	link.SetAuthKey(sc.akey)

	err = link.SendSyncMessage(proto.MessageObject{Object: obj}, proto.MessageTypeDeleteObject)
	if err != nil {
		return fmt.Errorf("[DeleteObject]%w", err)
	}

	maa, err := link.Await()
	if err != nil {
		return fmt.Errorf("[DeleteObject]%w", err)
	}
	if maa.Type == proto.MessageTypeError {
		a, err := maa.ReadError()
		if err != nil {
			return fmt.Errorf("[DeleteObject]%w", err)
		}
		return fmt.Errorf("sync error #%d: %s", a.ErrorCode, a.Error)
	} else if maa.Type == proto.MessageTypeClose {
		//We try to close only after server says that we can do it (no more answers expected)
		return
	} else {
		return fmt.Errorf("[DeleteObject] unexpected answer type '%s'", maa.Type)
	}
}
