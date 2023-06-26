package v1

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lazybark/go-cloud-sync/pkg/fse"
	"github.com/lazybark/go-helpers/hasher"
)

type FileProcessor struct {
	escSymbol string
	root      string
}

func NewFP(escSymbol, root string) *FileProcessor {
	fp := FileProcessor{
		escSymbol: escSymbol,
		root:      root,
	}

	return &fp
}

func (f *FileProcessor) GetPathUnescaped(obj fse.FSObject) string {
	return filepath.Join(f.UnEscapePath(obj.Path), obj.Name)
}

// ProcessDirectory returns full list of objects in the directory recursively.
// Unavailable for hashing are skipped.
func (f *FileProcessor) ProcessDirectory(path string) (objs []fse.FSObject, err error) {
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

func (fs *FileProcessor) scanDir(path string) (objs []fse.FSObject, err error) {
	contents, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	var o fse.FSObject
	var fullPath string
	var objs1 []fse.FSObject
	for _, item := range contents {
		fullPath = filepath.Join(path, item.Name())
		o = fse.FSObject{Path: fullPath}
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

func (f *FileProcessor) checkPathConsistency(path string) bool {
	// Check if path belongs to root dir
	if ok := strings.Contains(path, f.root); !ok {
		return ok
	}
	// Only absolute paths are available
	if ok := filepath.IsAbs(path); !ok {
		return ok
	}
	// Must exist and be a directory
	dir, err := os.Stat(path)
	if err != nil {
		return false
	}
	if !dir.IsDir() {
		return false
	}

	return true
}

func (fp *FileProcessor) ProcessObject(obj fse.FSObject, checkHash bool) (fse.FSObject, error) {
	if obj.IsProcessed {
		return obj, nil
	}

	var err error
	oInfo, err := os.Stat(obj.Path)
	if err != nil {
		return obj, fmt.Errorf("[CreateObject] object reading failed '%s': %w", obj.Path, err)
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

func (fp *FileProcessor) EscapePath(path string) string {
	path = strings.ReplaceAll(path, fp.root, "?ROOT_DIR?")
	return strings.ReplaceAll(path, string(filepath.Separator), fp.escSymbol)
}

func (fp *FileProcessor) UnEscapePath(path string) string {
	path = strings.ReplaceAll(path, "?ROOT_DIR?", fp.root)
	return strings.ReplaceAll(path, fp.escSymbol, string(filepath.Separator))
}

func (fp *FileProcessor) ConvertPathName(obj fse.FSObject) (dir, name string, err error) {
	dir, name = filepath.Split(obj.Path)
	dir, err = filepath.Abs(dir)
	if err != nil {
		return
	}
	dir = fp.EscapePath(dir)

	return
}

func (fp *FileProcessor) CheckFileHash(obj fse.FSObject) string {
	var sleep int
	var hash string
	var err error
	for {
		hash, err = hasher.HashFilePath(obj.Path, hasher.SHA256, 8192)
		if err != nil {
			if sleep >= 3 {
				//If object isn't readable - it's just ignored until next action.
				//Deprecating rescanBuffer for now. Seems useless as we still recieve info about new action
				//and we can create object after it's modified.
				break
			}
			time.Sleep(time.Second * 1)
			sleep++
		} else {
			break
		}
	}

	return hash
}

func (fp *FileProcessor) CheckDirHash(obj fse.FSObject) string {
	return ""
}
