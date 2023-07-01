package storage

type IServerStorageV1 interface {
	CreateObject(obj FSObjectStored) error
	RemoveObject(obj FSObjectStored, recursive bool) (err error)
	UpdateObject(obj FSObjectStored) error
	AddOrUpdateObject(obj FSObjectStored) error
	RefillDatabase(objs []FSObjectStored) error
	GetUsersObjects(owner string) ([]FSObjectStored, error)
	GetObject(path, name string) (obj FSObjectStored, err error)
	LockObject(path, name string) error
	UnLockObject(path, name string) error
}
