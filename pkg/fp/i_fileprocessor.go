package fp

import (
	"os"

	v1 "github.com/lazybark/go-cloud-sync/pkg/fp/v1"
	"github.com/lazybark/go-cloud-sync/pkg/fselink/v1/proto"
)

type FileprocessorV1 interface {
	ProcessObject(obj proto.FSObject, checkHash bool) (proto.FSObject, error)
	ConvertPathName(obj proto.FSObject) (dir, name string, err error)
	ProcessDirectory(path string) ([]proto.FSObject, error)
	GetPathUnescaped(obj proto.FSObject) string
	UnescapePath(obj proto.FSObject) string
	CreateFileInCache() (file *os.File, err error)
	DeleteFileInCache(path string) (err error)
	OpenToRead(path string) (file *os.File, err error)
}

func NewFPv1(escSymbol, root, cacheRoot string) FileprocessorV1 {
	return v1.NewFP(escSymbol, root, cacheRoot)
}
