package v1

import (
	proto "github.com/lazybark/go-cloud-sync/pkg/fselink/proto/v1"
)

func (s *FSWServer) processAuth(log string, pwd string, c *syncConnection) {
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
}
