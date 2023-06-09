package server

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lazybark/go-cloud-sync/pkg/storage"
	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/proto"
)

func (s *FSWServer) processPushFile(c *SyncConnection, m proto.ExchangeMessage) {
	a, err := m.ReadObjectData()
	if err != nil {
		c.SendError(proto.ErrMessageReadingFailed)
		s.extErc <- err
		return
	}

	a.Object.Path = strings.ReplaceAll(a.Object.Path, "?ROOT_DIR?", "?ROOT_DIR?,"+c.uid)
	pathFullUnescaped := filepath.Join(s.fp.GetPathUnescaped(a.Object))
	action := proto.Create

	dbObj, err := s.stor.GetObject(a.Object.Path, a.Object.Name)
	if err != nil && err != storage.ErrNotExists {
		c.SendError(proto.ErrInternalServerError)
		s.extErc <- err
		return
	}
	//If object does not exist yet
	if dbObj.ID != 0 {
		action = proto.Write
		err = s.stor.LockObject(a.Object.Path, a.Object.Name)
		if err != nil {
			s.extErc <- err
			c.SendError(proto.ErrInternalServerError)
			return
		}
		//Do not sync dirs
		if a.Object.IsDir {
			err = c.SendMessage(nil, proto.MessageTypeClose)
			if err != nil {
				s.extErc <- err
				return
			}
			c.Close()
			return
		}
	}
	if a.Object.IsDir {
		fmt.Println("CREATING DIR")
		//CREATE DIR HERE
		if err := os.MkdirAll(pathFullUnescaped, os.ModePerm); err != nil {
			s.extErc <- err
			c.SendError(proto.ErrInternalServerError)
			return
		}
		err = s.stor.AddOrUpdateObject(storage.FSObjectStored{
			Name:        a.Object.Name,
			Path:        a.Object.Path,
			Owner:       c.uid,
			Hash:        a.Object.Hash,
			IsDir:       a.Object.IsDir,
			Ext:         a.Object.Ext,
			Size:        a.Object.Size,
			FSUpdatedAt: a.Object.UpdatedAt,
		})
		if err != nil {
			s.extErc <- err
			c.SendError(proto.ErrInternalServerError)
			return
		}

		err = c.SendMessage(nil, proto.MessageTypeClose)
		if err != nil {
			s.extErc <- err
			return
		}
		c.Close()
		return
	}

	err = c.SendMessage(nil, proto.MessageTypePeerReady)
	if err != nil {
		s.extErc <- err
		return
	}

	file, err := s.fp.NewEmptyCache(a.Object)
	if err != nil {
		s.extErc <- err
		c.SendError(proto.ErrInternalServerError)
		return
	}
	for !c.IsClosed() {
		m, err := c.Await()
		if err != nil {
			s.extErc <- err
		}
		if m.Type == proto.MessageTypeError {
			a, err := m.ReadError()
			if err != nil {
				s.extErc <- err
				return
			}
			s.extErc <- fmt.Errorf("sync error #%d: %s", a.ErrorCode, a.Error)
			return
		} else if m.Type == proto.MessageTypeFileParts {
			a, err := m.ReadFilePart()
			if err != nil {
				s.extErc <- err
				return
			}
			_, err = file.Write(a.Payload)
			if err != nil {
				s.extErc <- err
				return
			}

		} else if m.Type == proto.MessageTypeFileEnd {
			break
		} else {
			s.extErc <- fmt.Errorf("[DownloadObject] unexpected answer type '%s'", m.Type)
			return
		}
	}

	err = file.Close()
	if err != nil {
		s.extErc <- err
		c.SendError(proto.ErrInternalServerError)
		return
	}

	err = s.fp.ReplaceFromCache(file)
	if err != nil {
		s.extErc <- err
		c.SendError(proto.ErrInternalServerError)
		err = file.Remove()
		if err != nil {
			s.extErc <- err
		}
		return
	}

	err = s.stor.AddOrUpdateObject(storage.FSObjectStored{
		Name:        a.Object.Name,
		Path:        a.Object.Path,
		Owner:       c.uid,
		Hash:        a.Object.Hash,
		IsDir:       a.Object.IsDir,
		Ext:         a.Object.Ext,
		Size:        a.Object.Size,
		FSUpdatedAt: a.Object.UpdatedAt,
	})
	if err != nil {
		s.extErc <- err
		c.SendError(proto.ErrInternalServerError)
		return
	}

	err = c.SendMessage(nil, proto.MessageTypeClose)
	if err != nil {
		s.extErc <- err
		c.SendError(proto.ErrInternalServerError)
		return
	}

	c.Close()

	//Now notify other clients for this user
	go s.notifyClients(a.Object, action, c.ID(), c.uid)

	if dbObj.ID != 0 {
		err = s.stor.UnLockObject(a.Object.Path, a.Object.Name)
		if err != nil {
			s.extErc <- err
			c.SendError(proto.ErrInternalServerError)
			return
		}
	}

}
