package storage

import "fmt"

type StorageError error

var (
	ErrNotExists StorageError = fmt.Errorf("object doesn't exist")
)
