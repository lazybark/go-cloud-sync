package sqlitestorage

import (
	"fmt"
	"log"
	"os"

	"github.com/lazybark/go-cloud-sync/pkg/fse"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"
)

func NewSQLite(path, escSymbol string) (*SQLiteStorage, error) {
	s := &SQLiteStorage{
		escSymbol: escSymbol,
	}

	db, err := gorm.Open(sqlite.Open("fs_watcher.db"), &gorm.Config{Logger: gLogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gLogger.Config{LogLevel: gLogger.Silent},
	)})
	if err != nil {
		return nil, err
	}
	s.db = db

	err = s.db.AutoMigrate(FSObject{})
	if err != nil {
		return nil, err
	}

	return s, nil
}

type SQLiteStorage struct {
	db        *gorm.DB
	escSymbol string
}

func (s *SQLiteStorage) RefillDatabase(objs []fse.FSObject) error {
	//return s.db.Create(&objs).Error
	return nil
}

func (s *SQLiteStorage) CreateObject(obj fse.FSObject) error {
	o := FSObject{
		Path:        obj.Path,
		Name:        obj.Name,
		IsDir:       obj.IsDir,
		Hash:        obj.Hash,
		Ext:         obj.Ext,
		Size:        obj.Size,
		FSUpdatedAt: obj.UpdatedAt,
	}
	if err := s.db.Create(&o).Error; err != nil {
		return err
	}
	return nil
}

func (s *SQLiteStorage) RemoveObject(obj fse.FSObject, recursive bool) error {
	var o FSObject
	if err := s.db.Where("name = ? and path = ?", obj.Name, obj.Path).First(&o).Error; err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if o.ID != 0 {
		if err := s.db.Delete(&o).Error; err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
	}
	fmt.Println(`DELETE FROM fs_objects WHERE "path" LIKE "%` + obj.Path + s.escSymbol + obj.Name + s.escSymbol + `%" OR "path" = "` + obj.Path + s.escSymbol + obj.Name + `"`)

	if o.IsDir && recursive {
		/*if err := s.db.Where("path LIKE %?%", obj.Path+s.escSymbol+obj.Name).Delete(&FSObject{}).Error; err != nil && err != gorm.ErrRecordNotFound {
			return err
		}*/
		//s.db.Delete(&FSObject{}, "path LIKE ?", "%"+obj.Path+s.escSymbol+obj.Name+"%")
		if err := s.db.Exec(`DELETE FROM fs_objects WHERE "path" LIKE "%` + obj.Path + s.escSymbol + obj.Name + s.escSymbol + `%" OR "path" = "` + obj.Path + s.escSymbol + obj.Name + `"`).Error; err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		/*childrenList := strings.Split(obj.Path+s.escSymbol+obj.Name, s.escSymbol)
		fmt.Println(childrenList)
		var childrenPaths [][]string
		for step := range childrenList {
			childrenPaths = append(childrenPaths, childrenList[0:step])
		}
		for _, c := range childrenPaths {
			fmt.Println(filepath.Join(c...))
			if err := s.db.Where("path = ?", filepath.Join(c...)).Delete(&o).Error; err != nil && err != gorm.ErrRecordNotFound {
				return err
			}
		}*/
	}

	return nil
}

func (s *SQLiteStorage) UpdateObject(obj fse.FSObject) error {
	return s.AddOrUpdateObject(obj)
}

func (s *SQLiteStorage) AddOrUpdateObject(obj fse.FSObject) error {
	o := FSObject{}
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
