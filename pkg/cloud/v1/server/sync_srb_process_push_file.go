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
	pathUnescaped := filepath.Join(s.fp.UnescapePath(a.Object))
	pathFullUnescaped := filepath.Join(s.fp.GetPathUnescaped(a.Object))

	//fileName := s.fp.GetPathUnescaped(mu.Object)
	dbObj, err := s.stor.GetObject(a.Object.Path, a.Object.Name)
	if err != nil && err != storage.ErrNotExists {
		c.SendError(proto.ErrInternalServerError)
		s.extErc <- err
		return
	}
	//If object does not exist yet
	if dbObj.ID != 0 {
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

	destFile, err := s.fp.CreateFileInCache()
	if err != nil {
		s.extErc <- err
		c.SendError(proto.ErrInternalServerError)
		return
	}
	//HERE WE SHOULD WAIT FOR FILE PARTS: intercept all connection messages until error or filebytes ended
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
			_, err = destFile.Write(a.Payload)
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

	destFile.Close()

	if err := os.MkdirAll(pathUnescaped, os.ModePerm); err != nil {
		s.extErc <- err
		err = s.fp.DeleteFileInCache(destFile.Name())
		if err != nil {
			s.extErc <- err
		}
		return
	}
	err = os.Rename(destFile.Name(), pathFullUnescaped)
	if err != nil {
		s.extErc <- err
		err = s.fp.DeleteFileInCache(destFile.Name())
		if err != nil {
			s.extErc <- err
		}
		return
	}
	err = os.Chtimes(pathFullUnescaped, a.Object.UpdatedAt, a.Object.UpdatedAt)
	if err != nil {
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
		c.SendError(proto.ErrInternalServerError)
		return
	}

	c.Close()
	if dbObj.ID != 0 {
		err = s.stor.UnLockObject(a.Object.Path, a.Object.Name)
		if err != nil {
			s.extErc <- err
			c.SendError(proto.ErrInternalServerError)
			return
		}
	}
	fmt.Println("DOWNLOADED FILE")
}
