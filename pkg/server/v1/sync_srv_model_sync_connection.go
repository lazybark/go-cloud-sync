package v1

import "github.com/lazybark/go-tls-server/v3/server"

type syncConnection struct {
	uid             string
	clientTokenHash string
	tlsConnection   *server.Connection
}

func (sc *syncConnection) ID() string {
	return sc.tlsConnection.Id()
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
