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
	plb, err := json.Marshal(MessageError{Error: e.String(), ErrorCode: e.Int()})
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

func AwaitAnswer(sc SyncReciever, m *ExchangeMessage) error {
	ans, err := sc.AwaitAnswer()
	if err != nil {
		return fmt.Errorf("[AwaitAnswer] %w", err)
	}
	err = json.Unmarshal(ans.Bytes(), m)
	if err != nil {
		return fmt.Errorf("[AwaitAnswer] %w", err)
	}

	return nil
}

func UnpackMessage(m ExchangeMessage, expectedType ExchangeMessageType, payload any) error {
	if m.Type != expectedType {
		return fmt.Errorf("[UnpackMessage] unexpected message type '%s'", m.Type)
	}

	err := json.Unmarshal(m.Payload, payload)
	if err != nil {
		return fmt.Errorf("[AwaitAnswer] %w", err)
	}

	return nil
}

/*
func AwaitAnswer(sc SyncReciever, payloadType any, expectedType ExchangeMessageType) (MessageError, error) {
	ans, err := sc.AwaitAnswer()
	if err != nil {
		return MessageError{}, fmt.Errorf("[AwaitAnswer] %w", err)
	}
	var ms ExchangeMessage
	err = json.Unmarshal(ans.Bytes(), &ms)
	if err != nil {
		return MessageError{}, fmt.Errorf("[AwaitAnswer] %w", err)
	}
	if ms.Type != expectedType && ms.Type != MessageTypeError {
		return MessageError{}, fmt.Errorf("unexpected message type '%s'", ms.Type)
	}
	if ms.Type == MessageTypeError {
		var e MessageError
		err = json.Unmarshal(ms.Payload, &e)
		if err != nil {
			return e, fmt.Errorf("[AwaitAnswer] %w", err)
		}
	}

	err = json.Unmarshal(ms.Payload, payloadType)
	if err != nil {
		return MessageError{}, fmt.Errorf("[AwaitAnswer] %w", err)
	}
	return MessageError{}, nil
}*/
