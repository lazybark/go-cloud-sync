package sqlitestorage

import (
	"gorm.io/gorm"
)

type SQLiteStorage struct {
	db        *gorm.DB
	escSymbol string
}
