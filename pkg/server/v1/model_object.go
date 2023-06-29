package v1

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/lazybark/go-cloud-sync/pkg/fse"
)

type Object struct {
	Path        Path
	IsDir       bool
	Hash        string
	Ext         string
	UserID      string
	Size        int64
	UpdatedAt   time.Time
	IsProcessed bool
}

func (o Object) JSON() ([]byte, error) {
	ebg, err := json.Marshal(o)
	if err != nil {
		return ebg, fmt.Errorf("[FSObject->JSON] %w", err)
	}

	return ebg, nil
}

func (o Object) ToProtoObject() fse.FSObject {
	return fse.FSObject{
		Path:      o.Path.PathEscaped,
		Name:      o.Path.Name,
		IsDir:     o.IsDir,
		Hash:      o.Hash,
		Ext:       o.Ext,
		Size:      o.Size,
		UpdatedAt: o.UpdatedAt,
	}
}

type Path struct {
	Name        string
	PathClean   string
	PathEscaped string
}