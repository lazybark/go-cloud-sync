package fselink

import (
	"os"

	"github.com/lazybark/go-cloud-sync/pkg/fse"
	gts "github.com/lazybark/go-tls-server/v2/server"
)

// FSEClientLink is used by clients to connect to a server
type FSEClientLink interface {
	Init(port int, addr, login, pwd string) error
	//Break() error
	GetObjList() ([]fse.FSObject, error)
	DownloadObject(obj fse.FSObject, writeTo *os.File) error
	SendEvent(e fse.FSEvent) error
}

// FSEServerPool is used by server to keep
type FSEServerPool interface {
	Init(chan (*gts.Message), chan (*gts.Connection), chan (error)) error
	Listen(string, string) error
	//Stop() error
	//NotifyClients(e fse.FSEvent) error
}
