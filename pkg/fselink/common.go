package fselink

import (
	"encoding/json"
	"fmt"
)

func SendSyncMessage(sc SyncMessenger, payload any, mType ExchangeMessageType) error {
	plb, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("[sendSyncMessage] %w", err)
	}
	mess, err := json.Marshal(ExchangeMessage{Type: mType, Payload: plb})
	if err != nil {
		return fmt.Errorf("[sendSyncMessage] %w", err)
	}
	_, err = sc.SendByte(mess)
	if err != nil {
		return fmt.Errorf("[sendSyncMessage] %w", err)
	}
	return nil
}

func SendErrorMessage(sc SyncMessenger, e ErrorCode) error {
	plb, err := json.Marshal(MessageError{Success: false, Error: e.String(), ErrorCode: e.Int()})
	if err != nil {
		return fmt.Errorf("[sendSyncMessage] %w", err)
	}
	mess, err := json.Marshal(ExchangeMessage{Type: MessageTypeError, Payload: plb})
	if err != nil {
		return fmt.Errorf("[sendSyncMessage] %w", err)
	}
	_, err = sc.SendByte(mess)
	if err != nil {
		return fmt.Errorf("[sendSyncMessage] %w", err)
	}
	return nil
}

func AwaitAnswer(sc SyncReciever, payloadType any) error {
	ans, err := sc.AwaitAnswer()
	if err != nil {
		return fmt.Errorf("[AwaitAnswer] %w", err)
	}
	var ms ExchangeMessage
	err = json.Unmarshal(ans.Bytes(), &ms)
	if err != nil {
		return fmt.Errorf("[AwaitAnswer] %w", err)
	}
	if ms.Type == MessageTypeError {
		fmt.Println(string(ms.Payload))
	}
	err = json.Unmarshal(ms.Payload, payloadType)
	if err != nil {
		return fmt.Errorf("[AwaitAnswer] %w", err)
	}
	return nil
}
