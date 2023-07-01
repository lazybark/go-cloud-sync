package client

import (
	"fmt"

	"github.com/lazybark/go-cloud-sync/pkg/fselink/v1/proto"
)

func (sc *LinkClient) GetObjList() (l []proto.FSObject, err error) {
	err = sc.SendSyncMessage(nil, proto.MessageTypeFullSyncRequest)
	if err != nil {
		err = fmt.Errorf("[GetObjList] %w", err)
		return
	}
	maa, err := sc.Await()
	if err != nil {
		err = fmt.Errorf("[GetObjList] %w", err)
		return
	}
	if maa.Type == proto.MessageTypeError {
		a, err := maa.ReadError()
		if err != nil {
			return l, fmt.Errorf("[GetObjList] %w", err)
		}
		return l, fmt.Errorf("sync error #%d: %s", a.ErrorCode, a.Error)
	} else if maa.Type == proto.MessageTypeFullSyncReply {
		a, err := maa.ReadFullSyncReply()
		if err != nil {
			return l, fmt.Errorf("[GetObjList] %w", err)
		}
		l = a.Objects

		return l, err
	}

	return l, fmt.Errorf("[GetObjList] unexpected answer type '%s'", maa.Type)
}
