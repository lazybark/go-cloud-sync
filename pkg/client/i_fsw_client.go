package client

import (
	v1 "github.com/lazybark/go-cloud-sync/pkg/client/v1"
	"github.com/lazybark/go-cloud-sync/pkg/fse"
	"github.com/lazybark/go-cloud-sync/pkg/storage"
)

type IFSWClient interface {
	Init(evc chan (fse.FSEvent), erc chan (error)) error
	Start() error
	Stop() error
}

func NewClientV1(stor storage.IStorage, cacheDir, root string) IFSWClient {
	return v1.NewClient(stor, cacheDir, root)
}
