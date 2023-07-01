package v1

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/proto"
)

// ProcessDirectory returns full list of objects in the directory recursively.
// Unavailable for hashing are skipped.
func (f *FileProcessor) ProcessDirectory(path string) (objs []proto.FSObject, err error) {
	unescaped := f.UnEscapePath(path)

	ok := f.checkPathConsistency(unescaped)
	if !ok {
		err = fmt.Errorf("[ProcessDirectory] provided path is not valid or consistent")
		return
	}
	objs, err = f.scanDir(path)
	if err != nil {
		return
	}

	return
}

func (fs *FileProcessor) scanDir(path string) (objs []proto.FSObject, err error) {
	contents, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	var o proto.FSObject
	var fullPath string
	var objs1 []proto.FSObject
	for _, item := range contents {
		fullPath = filepath.Join(path, item.Name())
		o = proto.FSObject{Path: fullPath}
		o, err = fs.ProcessObject(o, true)
		if err != nil {
			return
		}
		objs = append(objs, o)
		//Recursively scan sub dirs
		if o.IsDir {
			objs1, err = fs.scanDir(fullPath)
			if err != nil {
				return
			}
			objs = append(objs, objs1...)
		}
	}
	return
}
