package proto

import (
	"encoding/json"
	"fmt"
)

type Credentials struct {
	Login    string
	Password string
}

type MessageAuthAnswer struct {
	Success bool
	AuthKey string
}

func (em ExchangeMessage) ReadAuthAnswer() (MessageAuthAnswer, error) {
	var m MessageAuthAnswer
	if em.Type != MessageTypeAuthAns {
		return m, fmt.Errorf("[ExchangeMessage][ReadAuthAnswer] unexpected message type '%s'", em.Type)
	}

	err := json.Unmarshal(em.Payload, &m)
	if err != nil {
		return m, fmt.Errorf("[ExchangeMessage][ReadAuthAnswer] %w", err)
	}

	return m, nil
}
