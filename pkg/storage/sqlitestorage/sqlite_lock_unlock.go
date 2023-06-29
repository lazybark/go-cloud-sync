package sqlitestorage

import (
	"github.com/lazybark/go-cloud-sync/pkg/storage"

	"gorm.io/gorm"
)

func (s *SQLiteStorage) LockObject(path, name string) error {
	o := storage.FSObjectStored{}
	if err := s.db.Where("name = ? and path = ?", name, path).First(&o).Error; err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if o.ID == 0 {
		return storage.ErrNotExists

	}
	o.IsLocked = true

	if err := s.db.Save(&o).Error; err != nil {
		return err
	}

	return nil
}

func (s *SQLiteStorage) UnLockObject(path, name string) error {
	o := storage.FSObjectStored{}
	if err := s.db.Where("name = ? and path = ?", name, path).First(&o).Error; err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if o.ID == 0 {
		return storage.ErrNotExists

	}
	o.IsLocked = false

	if err := s.db.Save(&o).Error; err != nil {
		return err
	}

	return nil
}
