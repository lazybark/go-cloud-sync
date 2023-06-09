package server

import (
	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/proto"
)

func (s *FSWServer) processAuth(log string, pwd string, c *SyncConnection) {
	user, sessionKey := s.createSession(log, pwd)
	if sessionKey == "" {
		c.SendError(proto.ErrCodeWrongCreds)
		return
	}
	c.uid = user

	hash, err := hashAndSaltPassword([]byte(sessionKey))
	if err != nil {
		s.extErc <- err
		return
	}

	err = c.SendMessage(proto.MessageAuthAnswer{Success: true, AuthKey: sessionKey}, proto.MessageTypeAuthAns)
	if err != nil {
		s.extErc <- err
		return
	}
	c.clientTokenHash = hash
	c.basicSession = true //To notify stat routines
}
