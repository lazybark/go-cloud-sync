package sqlitestorage

import (
	"github.com/lazybark/go-cloud-sync/pkg/storage"
	"gorm.io/gorm"
)

func (s *SQLiteStorage) GetObject(path, name string) (obj storage.FSObjectStored, err error) {
	var o storage.FSObjectStored
	if err := s.db.Where("name = ? and path = ?", name, path).First(&o).Error; err != nil && err != gorm.ErrRecordNotFound {
		return o, err
	}
	if o.ID == 0 {
		return o, storage.ErrNotExists
	}

	return o, nil
}

func (s *SQLiteStorage) GetUsersObjects(owner string) ([]storage.FSObjectStored, error) {
	var objs []storage.FSObjectStored
	err := s.db.Where("owner = ?", owner).Find(&objs).Error

	return objs, err
}
