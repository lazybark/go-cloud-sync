package fselink

import (
	"encoding/json"
	"fmt"

	proto "github.com/lazybark/go-cloud-sync/pkg/fselink/proto/v1"
)

func (sc *SyncClient) SendSyncMessage(payload any, mType proto.ExchangeMessageType) error {
	plb, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("[sendSyncMessage]%w", err)
	}
	mess, err := json.Marshal(proto.ExchangeMessage{Type: mType, AuthKey: sc.akey, Payload: plb})
	if err != nil {
		return fmt.Errorf("[sendSyncMessage]%w", err)
	}
	_, err = sc.SendByte(mess)
	if err != nil {
		return fmt.Errorf("[sendSyncMessage]%w", err)
	}
	return nil
}
