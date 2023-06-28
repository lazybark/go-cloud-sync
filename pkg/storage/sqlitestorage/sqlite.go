package sqlitestorage

import (
	"fmt"
	"log"
	"os"

	"github.com/lazybark/go-cloud-sync/pkg/storage"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"
)

func NewSQLite(path, escSymbol string) (*SQLiteStorage, error) {
	s := &SQLiteStorage{
		escSymbol:  escSymbol,
		rootSymbol: "?ROOT_DIR?",
	}

	db, err := gorm.Open(sqlite.Open("fs_watcher_server.db"), &gorm.Config{Logger: gLogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gLogger.Config{LogLevel: gLogger.Silent},
	)})
	if err != nil {
		return nil, err
	}
	s.db = db

	err = s.db.AutoMigrate(storage.FSObjectStored{})
	if err != nil {
		return nil, err
	}

	return s, nil
}

type SQLiteStorage struct {
	db         *gorm.DB
	escSymbol  string
	rootSymbol string
}

func (s *SQLiteStorage) RefillDatabase(objs []storage.FSObjectStored) error {
	err := s.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&storage.FSObjectStored{}).Error
	if err != nil {
		return err
	}
	/*
		var recs []FSObject
		var owner int
		for _, o := range objs {
			fmt.Println(o)
			owner, err = ExtractOwnerFromPath(o.Path, s.escSymbol)
			if err != nil {
				return err
			}
			if owner == 0 {
				continue
			}
			//We don't care about files in server's root - they don't belong to anyone
			if o.Path == s.rootSymbol {
				continue
			}

			recs = append(recs, ConvertObjectsToDB(o, owner))
		}

		fmt.Println(recs)*/

	return s.db.Create(&objs).Error
}

/*
func ConvertObjectsToDB(obj fse.FSObject, owner int) FSObject {
	return FSObject{
		Path:        obj.Path,
		Name:        obj.Name,
		IsDir:       obj.IsDir,
		Hash:        obj.Hash,
		Owner:       owner,
		Ext:         obj.Ext,
		Size:        obj.Size,
		FSUpdatedAt: obj.UpdatedAt,
	}
}*/

func (s *SQLiteStorage) CreateObject(obj storage.FSObjectStored) error {
	if err := s.db.Create(&obj).Error; err != nil {
		return err
	}
	return nil
}

func (s *SQLiteStorage) RemoveObject(obj storage.FSObjectStored, recursive bool) error {
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

func (s *SQLiteStorage) UpdateObject(obj storage.FSObjectStored) error {
	return s.AddOrUpdateObject(obj)
}

func (s *SQLiteStorage) AddOrUpdateObject(obj storage.FSObjectStored) error {
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

func (s *SQLiteStorage) GetUsersObjects(owner string) ([]storage.FSObjectStored, error) {
	var objs []storage.FSObjectStored
	err := s.db.Where("owner = ?", owner).Find(&objs).Error

	return objs, err
}
