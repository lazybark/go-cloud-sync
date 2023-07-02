package proto

import (
	"encoding/json"
	"fmt"
)

type MessageFullSyncReply struct {
	Success bool
	Objects []FSObject
}

func (em ExchangeMessage) ReadFullSyncReply() (MessageFullSyncReply, error) {
	var m MessageFullSyncReply
	if em.Type != MessageTypeFullSyncReply {
		return m, fmt.Errorf("[ExchangeMessage][ReadFullSyncReply] unexpected message type '%s'", em.Type)
	}

	err := json.Unmarshal(em.Payload, &m)
	if err != nil {
		return m, fmt.Errorf("[ExchangeMessage][ReadFullSyncReply] %w", err)
	}

	return m, nil
}

type MessageSyncEvent struct {
	Object FSObject
	Event  FSAction
}

func (em ExchangeMessage) ReadSyncEvent() (MessageSyncEvent, error) {
	var m MessageSyncEvent
	if em.Type != MessageTypeSyncEvent {
		return m, fmt.Errorf("[ExchangeMessage][ReadSyncEvent] unexpected message type '%s'", em.Type)
	}

	err := json.Unmarshal(em.Payload, &m)
	if err != nil {
		return m, fmt.Errorf("[ExchangeMessage][ReadSyncEvent] %w", err)
	}

	return m, nil
}
