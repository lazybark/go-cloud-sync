package v1

import (
	"fmt"
	"os"
)

func (f *FileProcessor) OpenToRead(path string) (file *os.File, err error) {
	flags := os.O_RDONLY
	theFile, err := os.OpenFile(path, flags, 0666)
	if err != nil {
		return nil, fmt.Errorf("[OpenToRead][OpenFile]: %w", err)
	}

	return theFile, nil
}
