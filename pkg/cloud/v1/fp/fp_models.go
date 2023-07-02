package fp

import (
	"os"

	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/proto"
)

type FileProcessor struct {
	escSymbol string
	root      string
	cacheRoot string
}

func NewFileProcessorV1(escSymbol, root, cacheRoot string) *FileProcessor {
	fp := FileProcessor{
		escSymbol: escSymbol,
		root:      root,
		cacheRoot: cacheRoot,
	}

	return &fp
}

type File struct {
	o    proto.FSObject
	file *os.File
}
