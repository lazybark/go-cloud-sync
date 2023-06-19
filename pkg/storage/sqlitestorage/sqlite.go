package sqlitestorage

import "github.com/lazybark/go-cloud-sync/pkg/fse"

func NewSQLite(path string) *SQLiteStorage {
	s := &SQLiteStorage{}
	return s
}

type SQLiteStorage struct {
}

func (s *SQLiteStorage) CreateObject(obj fse.FSObject) (string, error) {
	var hash string
	return hash, nil
}

func (s *SQLiteStorage) RemoveObject(obj fse.FSObject) error {
	return nil
}

func (s *SQLiteStorage) UpdateObject(obj fse.FSObject) (string, error) {
	var newHash string
	return newHash, nil
}

func (s *SQLiteStorage) AddOrUpdateObject(obj fse.FSObject) (string, error) {
	var hash string
	return hash, nil
}
