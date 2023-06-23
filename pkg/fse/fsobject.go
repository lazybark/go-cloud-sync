package fse

import (
	"encoding/json"
	"fmt"
	"time"
)

// FSObject represents and filesystem object that can be watched by IFilesystemWatcher
type FSObject struct {
	Path      string
	Name      string
	IsDir     bool
	Hash      string
	Ext       string
	Size      int64
	UpdatedAt time.Time
}

func (o FSObject) JSON() ([]byte, error) {
	ebg, err := json.Marshal(o)
	if err != nil {
		return ebg, fmt.Errorf("[FSObject->JSON] %w", err)
	}

	return ebg, nil
}
