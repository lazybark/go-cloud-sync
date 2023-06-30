package v1

import (
	"github.com/lazybark/go-cloud-sync/pkg/fselink"
	proto "github.com/lazybark/go-cloud-sync/pkg/fselink/proto/v1"
)

func (s *FSWServer) processAuth(log string, pwd string, c *syncConnection) {
	user, sessionKey := s.createSession(log, pwd)
	if sessionKey == "" {
		s.sendError(c.tlsConnection, proto.ErrCodeWrongCreds)
		return
	}
	c.uid = user

	hash, err := hashAndSaltPassword([]byte(sessionKey))
	if err != nil {
		s.extErc <- err
		return
	}

	err = fselink.SendSyncMessage(c.tlsConnection, proto.MessageAuthAnswer{Success: true, AuthKey: sessionKey}, proto.MessageTypeAuthAns)
	if err != nil {
		s.extErc <- err
		return
	}
	c.clientTokenHash = hash
}
