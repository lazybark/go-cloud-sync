package proto

import (
	"encoding/json"
	"fmt"
)

type MessageObject struct {
	Object FSObject
}

func (em ExchangeMessage) ReadObjectData() (MessageObject, error) {
	var m MessageObject
	if em.Type != MessageTypePushFile && em.Type != MessageTypeGetFile && em.Type != MessageTypeDeleteObject {
		return m, fmt.Errorf("[ExchangeMessage][ReadObjectData] unexpected message type '%s'", em.Type)
	}

	err := json.Unmarshal(em.Payload, &m)
	if err != nil {
		return m, fmt.Errorf("[ExchangeMessage][ReadObjectData] %w", err)
	}

	return m, nil
}
