package client

import (
	v1 "github.com/lazybark/go-cloud-sync/pkg/client/v1"
	"github.com/lazybark/go-cloud-sync/pkg/fselink/v1/proto"
)

type IFSWClient interface {
	Init(evc chan (proto.FSEvent), erc chan (error), login, pwd string) error
	Start() error
	Stop() error
}

func NewClientV1(cacheDir, root string) IFSWClient {
	return v1.NewClient(cacheDir, root)
}
