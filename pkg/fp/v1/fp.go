package v1

import (
	"fmt"
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

func (fp *FileProcessor) ProcessObject(obj fse.FSObject, checkHash bool) (fse.FSObject, error) {
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
