// storage has implementations of fse.IFilesystemWatcher interface
package storage

import "github.com/lazybark/go-cloud-sync/pkg/fse"

type IStorage interface {
	//CreateObject creates a record about a filesystem object in database
	CreateObject(obj fse.FSObject) (hash string, err error)

	//RemoveObject removes the record about a filesystem object in database
	RemoveObject(obj fse.FSObject) (err error)

	//UpdateObject updates the record about a filesystem object in database
	UpdateObject(obj fse.FSObject) (newHash string, err error)

	//AddOrUpdateObject creates a record about a filesystem object in database or updates if
	//record exists. Returns success = true if any of actions were successful
	AddOrUpdateObject(obj fse.FSObject) (hash string, err error)
}
