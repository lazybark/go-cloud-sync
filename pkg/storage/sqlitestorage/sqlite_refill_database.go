package sqlitestorage

import (
	"github.com/lazybark/go-cloud-sync/pkg/storage"
	"gorm.io/gorm"
)

func (s *SQLiteStorage) RefillDatabase(objs []storage.FSObjectStored) error {
	err := s.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&storage.FSObjectStored{}).Error
	if err != nil {
		return err
	}

	if len(objs) == 0 {
		return nil
	}

	return s.db.Create(&objs).Error
}
