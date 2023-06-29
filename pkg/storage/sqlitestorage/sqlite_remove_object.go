package sqlitestorage

import (
	"github.com/lazybark/go-cloud-sync/pkg/storage"
	"gorm.io/gorm"
)

func (s *SQLiteStorage) RemoveObject(obj storage.FSObjectStored, recursive bool) error {
	var o storage.FSObjectStored
	if err := s.db.Where("name = ? and path = ?", obj.Name, obj.Path).First(&o).Error; err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if o.ID != 0 {
		if err := s.db.Delete(&o).Error; err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
	}

	if o.IsDir && recursive {
		if err := s.db.Exec(`DELETE FROM fs_object_storeds WHERE "path" LIKE "%` + obj.Path + s.escSymbol + obj.Name + s.escSymbol + `%" OR "path" = "` + obj.Path + s.escSymbol + obj.Name + `"`).Error; err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
	}

	return nil
}
