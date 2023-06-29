package sqlitestorage

import (
	"log"
	"os"

	"github.com/lazybark/go-cloud-sync/pkg/storage"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"
)

func NewSQLite(path, escSymbol string) (*SQLiteStorage, error) {
	s := &SQLiteStorage{
		escSymbol: escSymbol,
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
