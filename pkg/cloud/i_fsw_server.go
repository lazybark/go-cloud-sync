package cloud

import (
	"github.com/lazybark/go-cloud-sync/pkg/cloud/v1/server"
	"github.com/lazybark/go-cloud-sync/pkg/storage"
)

type IFSWServerV1 interface {
	Init(root, cacheRoot, host, port, escSymbol string, erc chan (error)) error
	Start() error
	Stop() error
}

func NewServerV1(stor storage.IServerStorageV1) IFSWServerV1 {
	return server.NewServer(stor)
}
