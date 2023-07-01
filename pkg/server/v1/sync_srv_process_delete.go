package v1

import (
	"os"
	"strings"

	"github.com/lazybark/go-cloud-sync/pkg/storage"
	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/proto"
)

func (s *FSWServer) processDelete(c *syncConnection, m proto.ExchangeMessage) {
	a, err := m.ReadObjectData()
	if err != nil {
		c.SendError(proto.ErrMessageReadingFailed)
		s.extErc <- err
		return
	}

	a.Object.Path = strings.ReplaceAll(a.Object.Path, "?ROOT_DIR?", "?ROOT_DIR?,"+c.uid)

	//fileName := s.fp.GetPathUnescaped(mu.Object)
	dbObj, err := s.stor.GetObject(a.Object.Path, a.Object.Name)
	if err != nil && err != storage.ErrNotExists {
		c.SendError(proto.ErrInternalServerError)
		s.extErc <- err
		return
	}

	if dbObj.ID != 0 {
		err = s.stor.RemoveObject(dbObj, true)
		if err != nil && err != storage.ErrNotExists {
			c.SendError(proto.ErrInternalServerError)
			s.extErc <- err
			return
		}
	}

	err = os.RemoveAll(s.fp.GetPathUnescaped(a.Object))
	if err != nil && err != storage.ErrNotExists {
		c.SendError(proto.ErrInternalServerError)
		s.extErc <- err
		return
	}

	err = c.SendMessage(nil, proto.MessageTypeClose)
	if err != nil {
		s.extErc <- err
		return
	}

	err = c.tlsConnection.Close()
	if err != nil {
		s.extErc <- err
	}
}
