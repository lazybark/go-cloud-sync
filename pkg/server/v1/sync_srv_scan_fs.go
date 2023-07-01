package v1

import (
	"fmt"

	"github.com/lazybark/go-cloud-sync/pkg/storage"
)

func (s *FSWServer) rescanOnce() {
	local, err := s.readLocalObjects()
	if err != nil {
		s.extErc <- fmt.Errorf("[rescanOnce]%w", err)
		return
	}
	fmt.Printf("Root directory read. Found %d objects\n", len(local))

	err = s.stor.RefillDatabase(local)
	if err != nil {
		s.extErc <- fmt.Errorf("[rescanOnce]%w", err)
		return
	}
	fmt.Println("Database refilled")
}

func (s *FSWServer) readLocalObjects() (objs []storage.FSObjectStored, err error) {
	local, err := s.fp.ProcessDirectory(s.conf.root)
	if err != nil {
		return nil, fmt.Errorf("[FSWATCHER][DiffListWithServer]: %w", err)
	}
	for _, o := range local {
		if o.Path == "?ROOT_DIR?" {
			//Ignore all objs in server root
			continue
		}
		objs = append(objs, storage.FSObjectStored{
			Path:        o.Path,
			Name:        o.Name,
			IsDir:       o.IsDir,
			Hash:        o.Hash,
			Owner:       s.GetOwnerFromPath(o.Path),
			Ext:         o.Ext,
			Size:        o.Size,
			FSUpdatedAt: o.UpdatedAt,
		})
	}
	return
}
