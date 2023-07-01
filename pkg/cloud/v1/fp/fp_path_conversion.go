package fp

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/proto"
)

func (f *FileProcessor) GetPathUnescaped(obj proto.FSObject) string {
	return filepath.Join(f.UnEscapePath(obj.Path), obj.Name)
}

func (f *FileProcessor) UnescapePath(obj proto.FSObject) string {
	return filepath.Join(f.UnEscapePath(obj.Path))
}

func (fp *FileProcessor) EscapePath(path string) string {
	path = strings.ReplaceAll(path, fp.root, "?ROOT_DIR?")
	return strings.ReplaceAll(path, string(filepath.Separator), fp.escSymbol)
}

func (fp *FileProcessor) UnEscapePath(path string) string {
	path = strings.ReplaceAll(path, "?ROOT_DIR?", fp.root)
	return strings.ReplaceAll(path, fp.escSymbol, string(filepath.Separator))
}
func (fp *FileProcessor) UnEscapePathWithUser(path string, user string) string {
	path = strings.ReplaceAll(path, "?ROOT_DIR?", fp.root+fp.escSymbol+user+fp.escSymbol)
	return strings.ReplaceAll(path, fp.escSymbol, string(filepath.Separator))
}

func (fp *FileProcessor) ConvertPathName(obj proto.FSObject) (dir, name string, err error) {
	dir, name = filepath.Split(obj.Path)
	dir, err = filepath.Abs(dir)
	if err != nil {
		return
	}
	dir = fp.EscapePath(dir)

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
