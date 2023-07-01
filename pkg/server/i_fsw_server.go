package server

import (
	"github.com/lazybark/go-cloud-sync/pkg/fselink/v1/proto"
	v1 "github.com/lazybark/go-cloud-sync/pkg/server/v1"
	"github.com/lazybark/go-cloud-sync/pkg/storage"
)

type IFSWServerV1 interface {
	Init(root, cacheRoot, host, port, escSymbol string, evc chan (proto.FSEvent), erc chan (error)) error
	Start() error
	Stop() error
}

func NewServerV1(stor storage.IServerStorageV1) IFSWServerV1 {
	return v1.NewServer(stor)
}
