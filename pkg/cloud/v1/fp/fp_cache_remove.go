package fp

import (
	"fmt"
	"os"
	"strings"
)

func (f *FileProcessor) DeleteFileInCache(path string) (err error) {
	if !strings.HasPrefix(path, f.cacheRoot) {
		return fmt.Errorf("[DeleteFileInCache] can not delete file outside of cache root ('%s')", path)
	}
	return os.RemoveAll(path)
}
