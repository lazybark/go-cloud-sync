package server

import (
	"sync"

	"github.com/lazybark/go-cloud-sync/pkg/storage"
	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/server"
	gts "github.com/lazybark/go-tls-server/v3/server"
)

func NewServer(stor storage.IServerStorageV1) *FSWServer {
	s := &FSWServer{}
	s.srvConnChan = make(chan *gts.Connection)
	s.srvErrChan = make(chan error)
	s.srvMessChan = make(chan *gts.Message)
	s.connPool = make(map[string]*syncConnection)
	s.connPoolMutex = sync.RWMutex{}

	s.stor = stor
	s.htsrv = server.NewServer()

	return s
}
