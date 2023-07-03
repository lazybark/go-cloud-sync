package server

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/proto"
	"github.com/lazybark/go-tls-server/v3/server"
)

type SyncConnection struct {
	uid             string
	clientTokenHash string
	sendEvents      bool
	basicSession    bool

	//sendMutex makes sure that messages will be written in the specific connection in order.
	//
	sendMutex     *sync.Mutex
	tlsConnection *server.Connection
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
	//HERE WILL BE STAT DATA for server frontend. In case sc.basicSession == true, we set session as
	//done.
	//And we delete all connection with same token after this one is closed
	sc.tlsConnection.Close()
}

func (sc *SyncConnection) SendMessage(payload any, mType proto.ExchangeMessageType) error {
	if sc.IsClosed() {
		return nil
	}
	plb, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("[SendMessage] %w", err)
	}
	mess, err := json.Marshal(proto.ExchangeMessage{Type: mType, Payload: plb})
	if err != nil {
		return fmt.Errorf("[SendMessage] %w", err)
	}

	sc.sendMutex.Lock()
	defer sc.sendMutex.Unlock()
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

func (s *FSWServer) notifyClients(o proto.FSObject, e proto.FSAction, sourceCID string, sourceUID string) {
	fmt.Println(e, sourceCID)
	o.Path = s.ExtractOwnerFromPath(o.Path, sourceUID)
	for cid, c := range s.connPool {
		if !c.sendEvents {
			continue
		}
		//Ignore same connection
		if sourceCID == cid {
			continue
		}
		//Ignore other users
		if c.uid != sourceUID {
			continue
		}
		err := c.SendMessage(proto.MessageSyncEvent{Object: o, Event: e}, proto.MessageTypeSyncEvent)
		if err != nil {
			s.extErc <- err
			return
		}
	}
}
