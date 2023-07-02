package fp

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofrs/uuid"
	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/proto"
)

func (fp *FileProcessor) NewEmptyCache(obj proto.FSObject) (*File, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("[MakeCache][NewV4]%w", err)
	}

	tempPath := fp.cacheRoot + string(filepath.Separator) + "cache_getfile" + u.String()
	if err := os.MkdirAll(fp.cacheRoot, os.ModePerm); err != nil {
		return nil, fmt.Errorf("[MakeCache][MkdirAll]%w", err)
	}
	theFile, err := os.Create(tempPath)
	if err != nil {
		return nil, fmt.Errorf("[MakeCache][Create]%w", err)
	}

	file := &File{
		o:    obj,
		file: theFile,
	}

	return file, nil
}

func (fp *FileProcessor) ReplaceFromCache(file *File) error {
	pathUnescaped := filepath.Join(fp.UnescapePath(file.o))
	pathFullUnescaped := filepath.Join(fp.GetPathUnescaped(file.o))

	err := os.MkdirAll(pathUnescaped, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.Rename(file.file.Name(), pathFullUnescaped)
	if err != nil {
		return err
	}
	err = os.Chtimes(pathFullUnescaped, file.o.UpdatedAt, file.o.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (f *File) Write(b []byte) (int, error) {
	return f.file.Write(b)
}

func (f *File) Close() error {
	return f.file.Close()
}

func (f *File) Remove() error {
	return os.RemoveAll(f.file.Name())
}
