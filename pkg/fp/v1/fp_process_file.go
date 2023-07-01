package v1

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/proto"
)

func (fp *FileProcessor) ProcessObject(obj proto.FSObject, checkHash bool) (proto.FSObject, error) {
	if obj.IsProcessed {
		return obj, nil
	}

	var err error
	oInfo, err := os.Stat(obj.Path)
	if err != nil {
		return obj, fmt.Errorf("[ProcessObject][Stat] '%s': %w", obj.Path, err)
	}
	obj.IsDir = oInfo.IsDir()

	if checkHash {
		if obj.IsDir {
			obj.Hash = fp.CheckDirHash(obj)
		} else {
			obj.Hash = fp.CheckFileHash(obj)
		}
	}

	obj.Ext = filepath.Ext(obj.Path)
	obj.Size = oInfo.Size()
	obj.UpdatedAt = oInfo.ModTime()
	obj.Path, obj.Name, err = fp.ConvertPathName(obj)
	if err != nil {
		return obj, nil
	}
	obj.IsProcessed = true

	return obj, nil
}
