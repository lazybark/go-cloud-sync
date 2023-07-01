package fselink

import (
	"encoding/json"
	"fmt"

	proto "github.com/lazybark/go-cloud-sync/pkg/fselink/proto/v1"
)

func AwaitAnswer(sc SyncReciever, m *proto.ExchangeMessage) error {
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

func UnpackMessage(m proto.ExchangeMessage, expectedType proto.ExchangeMessageType, payload any) error {
	if m.Type != expectedType {
		return fmt.Errorf("[UnpackMessage] unexpected message type '%s'", m.Type)
	}

	err := json.Unmarshal(m.Payload, payload)
	if err != nil {
		return fmt.Errorf("[AwaitAnswer] %w", err)
	}

	return nil
}
