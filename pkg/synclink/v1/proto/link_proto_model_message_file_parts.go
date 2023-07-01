package proto

import (
	"encoding/json"
	"fmt"
)

type MessageFilePart struct {
	Payload []byte
}

func (em ExchangeMessage) ReadFilePart() (MessageFilePart, error) {
	var m MessageFilePart
	if em.Type != MessageTypeFileParts {
		return m, fmt.Errorf("[ExchangeMessage][ReadFilePart] unexpected message type '%s'", em.Type)
	}

	err := json.Unmarshal(em.Payload, &m)
	if err != nil {
		return m, fmt.Errorf("[ExchangeMessage][ReadFilePart] %w", err)
	}

	return m, nil
}
