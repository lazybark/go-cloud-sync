package proto

import (
	"encoding/json"
	"fmt"

	"github.com/lazybark/go-cloud-sync/pkg/fse"
)

type MessageFullSyncReply struct {
	Success bool
	Objects []fse.FSObject
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
