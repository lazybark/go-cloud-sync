package v1

import (
	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/proto"
)

func (s *FSWServer) processFullSyncRequest(c *syncConnection) {
	uo, err := s.stor.GetUsersObjects(c.uid)
	if err != nil {
		s.extErc <- err
		return
	}
	var l []proto.FSObject
	for _, ol := range uo {
		l = append(l, proto.FSObject{
			Path:      s.ExtractOwnerFromPath(ol.Path, c.uid),
			Name:      ol.Name,
			IsDir:     ol.IsDir,
			Hash:      ol.Hash,
			Ext:       ol.Ext,
			Size:      ol.Size,
			UpdatedAt: ol.FSUpdatedAt,
		})
	}
	err = c.SendMessage(proto.MessageFullSyncReply{Success: true, Objects: l}, proto.MessageTypeFullSyncReply)
	if err != nil {
		s.extErc <- err
		return
	}
}
