package server

import (
	"time"

	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/proto"
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

func (o Object) ToProtoObject() proto.FSObject {
	return proto.FSObject{
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
