package server

import (
	"github.com/lazybark/go-cloud-sync/pkg/fse"
	v1 "github.com/lazybark/go-cloud-sync/pkg/server/v1"
	"github.com/lazybark/go-cloud-sync/pkg/storage"
)

type IFSWServer interface {
	Init(root, cacheRoot, host, port, escSymbol string, evc chan (fse.FSEvent), erc chan (error)) error
	Start() error
	Stop() error
}

func NewServerV1(stor storage.IServerStorage) IFSWServer {
	return v1.NewServer(stor)
}
