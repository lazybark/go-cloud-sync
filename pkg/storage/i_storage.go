// storage has implementations of fse.IFilesystemWatcher interface
package storage

import (
	"fmt"
	"time"
)

type IStorage interface {
	//CreateObject creates a record about a filesystem object in database
	CreateObject(obj FSObjectStored) error

	//RemoveObject removes the record about a filesystem object in database
	RemoveObject(obj FSObjectStored, recursive bool) (err error)

	//UpdateObject updates the record about a filesystem object in database
	UpdateObject(obj FSObjectStored) error

	//AddOrUpdateObject creates a record about a filesystem object in database or updates if
	//record exists. Returns success = true if any of actions were successful
	AddOrUpdateObject(obj FSObjectStored) error

	RefillDatabase(objs []FSObjectStored) error

	GetUsersObjects(owner string) ([]FSObjectStored, error)

	GetObject(path, name string) (obj FSObjectStored, err error)
	LockObject(path, name string) error
	UnLockObject(path, name string) error
}

type FSObjectStored struct {
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
	IsLocked    bool
}

type StorageError error

var (
	ErrNotExists StorageError = fmt.Errorf("object doesn't exist")
)
