package sqlitestorage

import (
	"github.com/lazybark/go-cloud-sync/pkg/storage"
	"gorm.io/gorm"
)

func (s *SQLiteStorage) CreateObject(obj storage.FSObjectStored) error {
	if err := s.db.Create(&obj).Error; err != nil {
		return err
	}
	return nil
}

func (s *SQLiteStorage) UpdateObject(obj storage.FSObjectStored) error {
	return s.AddOrUpdateObject(obj)
}

func (s *SQLiteStorage) AddOrUpdateObject(obj storage.FSObjectStored) error {
	o := storage.FSObjectStored{}
	if err := s.db.Where("name = ? and path = ?", obj.Name, obj.Path).First(&o).Error; err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if o.ID != 0 {
		//Only specific fields can be updated
		o.Hash = obj.Hash
		o.Size = obj.Size
		o.UpdatedAt = obj.UpdatedAt

		if err := s.db.Save(&o).Error; err != nil {
			return err
		}
	} else {
		if err := s.CreateObject(obj); err != nil {
			return err
		}
	}

	return nil
}
