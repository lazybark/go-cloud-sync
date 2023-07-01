package v1

import (
	"os"
	"path/filepath"

	"github.com/gofrs/uuid"
)

func (f *FileProcessor) CreateFileInCache() (file *os.File, err error) {
	u, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	tempPath := f.cacheRoot + string(filepath.Separator) + "cache_getfile" + u.String()
	if err := os.MkdirAll(f.cacheRoot, os.ModePerm); err != nil {
		return nil, err
	}
	theFile, err := os.Create(tempPath)
	if err != nil {
		return nil, err
	}

	return theFile, nil
}
