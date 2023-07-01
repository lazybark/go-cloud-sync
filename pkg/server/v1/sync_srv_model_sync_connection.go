package v1

import (
	"encoding/json"
	"fmt"

	proto "github.com/lazybark/go-cloud-sync/pkg/fselink/proto/v1"
	"github.com/lazybark/go-tls-server/v3/server"
)

type syncConnection struct {
	uid             string
	clientTokenHash string
	tlsConnection   *server.Connection
}

func (sc *syncConnection) ID() string {
	return sc.tlsConnection.Id()
}

func (sc *syncConnection) Close() {
	sc.tlsConnection.Close()
}

func (sc *syncConnection) SendMessage(payload any, mType proto.ExchangeMessageType) error {
	plb, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("[SendMessage] %w", err)
	}
	mess, err := json.Marshal(proto.ExchangeMessage{Type: mType, Payload: plb})
	if err != nil {
		return fmt.Errorf("[SendMessage] %w", err)
	}
	_, err = sc.tlsConnection.SendByte(mess)
	if err != nil {
		return fmt.Errorf("[SendMessage] %w", err)
	}
	return nil
}

func (sc *syncConnection) SendError(e proto.ErrorCode) error {
	pLoad := proto.MessageError{Error: e.String(), ErrorCode: e.Int()}

	err := sc.SendMessage(pLoad, proto.MessageTypeError)
	if err != nil {
		return fmt.Errorf("[SendError] %w", err)
	}
	return nil
}

func (s *FSWServer) addToPool(c *syncConnection) {
	s.connPoolMutex.Lock()
	s.connPool[c.ID()] = c
	s.connPoolMutex.Unlock()
}

// remFromPool removes connection pointer from pool, so it becomes unavailable to reach
func (s *FSWServer) remFromPool(c *syncConnection) {
	s.connPoolMutex.Lock()
	delete(s.connPool, c.ID())
	s.connPoolMutex.Unlock()
}
