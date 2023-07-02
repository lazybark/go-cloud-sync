package synclink

import (
	"os"

	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/client"
	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/proto"
	gts "github.com/lazybark/go-tls-server/v3/server"
)

type ISyncLinkClientV1 interface {
	Init(port int, addr, login, pwd string) error
	NewConnectionInSession() (*client.LinkClient, error)
	Await() (proto.ExchangeMessage, error)
	SendSyncMessage(payload any, mType proto.ExchangeMessageType) error
	//DownloadObject(obj proto.FSObject, writeTo *os.File) error
	PushObject(obj proto.FSObject, readFrom *os.File) error
}

type ISyncLinkServerV1 interface {
	Init(chan (*gts.Message), chan (*gts.Connection), chan (error)) error
	Listen(string, string) error
}
