package fselink

import (
	"os"

	"github.com/lazybark/go-cloud-sync/pkg/fselink/v1/proto"
	gts "github.com/lazybark/go-tls-server/v3/server"
)

// FSEClientLink is used by clients to connect to a server
type SyncLinkClient interface {
	Init(port int, addr, login, pwd string) error
	//Break() error
	GetObjList() ([]proto.FSObject, error)
	DownloadObject(obj proto.FSObject, writeTo *os.File) error
	PushObject(obj proto.FSObject, readFrom *os.File) error
	DeleteObject(obj proto.FSObject) (err error)
}

// FSEServerPool is used by server to keep
type FSEServerPool interface {
	Init(chan (*gts.Message), chan (*gts.Connection), chan (error)) error
	Listen(string, string) error
	//Stop() error
	//NotifyClients(e fse.FSEvent) error
}
