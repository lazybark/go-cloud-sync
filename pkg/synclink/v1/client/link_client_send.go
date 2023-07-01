package client

import (
	"encoding/json"
	"fmt"

	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/proto"
)

func (sc *LinkClient) SendByte(b []byte) (int, error) {
	err := sc.c.SendByte(b)
	if err != nil {
		err = fmt.Errorf("[SendByte]: %w", err)
	}
	return len(b), err
}

func (sc *LinkClient) SendSyncMessage(payload any, mType proto.ExchangeMessageType) error {
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
