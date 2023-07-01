package cloud

import (
	"github.com/lazybark/go-cloud-sync/pkg/cloud/v1/client"
	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/proto"
)

type IFSWClientV1 interface {
	Init(evc chan (proto.FSEvent), erc chan (error), login, pwd string) error
	Start() error
	Stop() error
}

func NewClientV1(cacheDir, root string) IFSWClientV1 {
	return client.NewClient(cacheDir, root)
}
