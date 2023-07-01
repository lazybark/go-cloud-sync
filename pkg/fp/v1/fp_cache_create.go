package v1

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofrs/uuid"
)

func (f *FileProcessor) CreateFileInCache() (file *os.File, err error) {
	u, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("[CreateFileInCache][NewV4]%w", err)
	}
	tempPath := f.cacheRoot + string(filepath.Separator) + "cache_getfile" + u.String()
	if err := os.MkdirAll(f.cacheRoot, os.ModePerm); err != nil {
		return nil, fmt.Errorf("[CreateFileInCache][MkdirAll]%w", err)
	}
	theFile, err := os.Create(tempPath)
	if err != nil {
		return nil, fmt.Errorf("[CreateFileInCache][Create]%w", err)
	}

	return theFile, nil
}
