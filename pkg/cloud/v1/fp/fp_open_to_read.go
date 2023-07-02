package fp

import (
	"fmt"
	"os"
)

// OpenToRead returns File object that is available to read. It does not fill the Object data field.
// That means physical file will be readable, but all sync data about it must be gathered from
// another source
func (f *FileProcessor) OpenToRead(path string) (file *File, err error) {
	flags := os.O_RDONLY
	theFile, err := os.OpenFile(path, flags, 0666)
	if err != nil {
		return nil, fmt.Errorf("[OpenToRead][OpenFile]: %w", err)
	}

	file = &File{
		//Object property is empty here!
		file: theFile,
	}

	return file, nil
}
