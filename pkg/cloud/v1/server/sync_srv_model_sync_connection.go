package server

import (
	"encoding/json"
	"fmt"

	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/proto"
	"github.com/lazybark/go-tls-server/v3/server"
)

type SyncConnection struct {
	uid             string
	clientTokenHash string
	tlsConnection   *server.Connection
}

func (sc *SyncConnection) Await() (proto.ExchangeMessage, error) {
	var m proto.ExchangeMessage
	ans := <-sc.tlsConnection.MessageChan
	err := json.Unmarshal(ans.Bytes(), &m)
	if err != nil {
		return m, fmt.Errorf("[Await] %w", err)
	}
	return m, nil
}

func (sc *SyncConnection) ID() string {
	return sc.tlsConnection.Id()
}

func (sc *SyncConnection) IsClosed() bool {
	return sc.tlsConnection.Closed()
}

func (sc *SyncConnection) Close() {
	sc.tlsConnection.Close()
}

func (sc *SyncConnection) SendMessage(payload any, mType proto.ExchangeMessageType) error {
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

func (sc *SyncConnection) SendError(e proto.ErrorCode) error {
	pLoad := proto.MessageError{Error: e.String(), ErrorCode: e.Int()}

	err := sc.SendMessage(pLoad, proto.MessageTypeError)
	if err != nil {
		return fmt.Errorf("[SendError] %w", err)
	}
	return nil
}

func (s *FSWServer) addToPool(c *SyncConnection) {
	s.connPoolMutex.Lock()
	s.connPool[c.ID()] = c
	s.connPoolMutex.Unlock()
}

// remFromPool removes connection pointer from pool, so it becomes unavailable to reach
func (s *FSWServer) remFromPool(c *SyncConnection) {
	s.connPoolMutex.Lock()
	delete(s.connPool, c.ID())
	s.connPoolMutex.Unlock()
}
