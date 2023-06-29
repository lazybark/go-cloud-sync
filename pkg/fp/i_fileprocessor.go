package fp

import (
	"os"

	v1 "github.com/lazybark/go-cloud-sync/pkg/fp/v1"
	"github.com/lazybark/go-cloud-sync/pkg/fse"
)

type Fileprocessor interface {
	ProcessObject(obj fse.FSObject, checkHash bool) (fse.FSObject, error)
	ConvertPathName(obj fse.FSObject) (dir, name string, err error)
	ProcessDirectory(path string) ([]fse.FSObject, error)
	GetPathUnescaped(obj fse.FSObject) string
	UnescapePath(obj fse.FSObject) string
	CreateFileInCache() (file *os.File, err error)
	DeleteFileInCache(path string) (err error)
	OpenToRead(path string) (file *os.File, err error)
}

func NewFPv1(escSymbol, root, cacheRoot string) Fileprocessor {
	return v1.NewFP(escSymbol, root, cacheRoot)
}
