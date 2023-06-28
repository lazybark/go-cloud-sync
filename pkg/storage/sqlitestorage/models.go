package sqlitestorage

import "time"

type (
	FSObject struct {
		ID          uint   `gorm:"primaryKey"`
		Name        string `gorm:"uniqueIndex:file"`
		Path        string `gorm:"uniqueIndex:file"`
		Owner       string `gorm:"uniqueIndex:file"`
		Hash        string
		IsDir       bool
		Ext         string
		Size        int64
		FSUpdatedAt time.Time
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}
)
