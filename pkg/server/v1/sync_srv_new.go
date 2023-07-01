package v1

import (
	"sync"

	"github.com/lazybark/go-cloud-sync/pkg/fselink/v1/server"
	"github.com/lazybark/go-cloud-sync/pkg/storage"
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
