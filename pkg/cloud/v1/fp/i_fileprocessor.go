package fp

import (
	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/proto"
)

type FileprocessorV1 interface {
	ProcessObject(obj proto.FSObject, checkHash bool) (proto.FSObject, error)
	ConvertPathName(obj proto.FSObject) (dir, name string, err error)
	ProcessDirectory(path string) ([]proto.FSObject, error)
	GetPathUnescaped(obj proto.FSObject) string
	UnescapePath(obj proto.FSObject) string
	OpenToRead(path string) (file *File, err error)

	NewEmptyCache(obj proto.FSObject) (file *File, err error)
	ReplaceFromCache(file *File) error
}

func NewFPv1(escSymbol, root, cacheRoot string) FileprocessorV1 {
	return NewFileProcessorV1(escSymbol, root, cacheRoot)
}
