package v1

import (
	"os"
	"strings"

	"github.com/lazybark/go-cloud-sync/pkg/fselink"
	proto "github.com/lazybark/go-cloud-sync/pkg/fselink/proto/v1"
	"github.com/lazybark/go-cloud-sync/pkg/storage"
)

func (s *FSWServer) processDelete(c *syncConnection, m proto.ExchangeMessage) {
	var mu proto.MessageDeleteObject
	err := fselink.UnpackMessage(m, proto.MessageTypeDeleteObject, &mu)
	if err != nil {
		c.SendError(proto.ErrMessageReadingFailed)
		s.extErc <- err
		return
	}

	mu.Object.Path = strings.ReplaceAll(mu.Object.Path, "?ROOT_DIR?", "?ROOT_DIR?,"+c.uid)

	//fileName := s.fp.GetPathUnescaped(mu.Object)
	dbObj, err := s.stor.GetObject(mu.Object.Path, mu.Object.Name)
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

	err = os.RemoveAll(s.fp.GetPathUnescaped(mu.Object))
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
