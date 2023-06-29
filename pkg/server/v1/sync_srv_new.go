package v1

import (
	"github.com/lazybark/go-cloud-sync/pkg/fselink"
	"github.com/lazybark/go-cloud-sync/pkg/storage"
	"github.com/lazybark/go-cloud-sync/pkg/watcher"
	gts "github.com/lazybark/go-tls-server/v3/server"
)

func NewServer(stor storage.IServerStorage) *FSWServer {
	s := &FSWServer{}
	s.srvConnChan = make(chan *gts.Connection)
	s.srvErrChan = make(chan error)
	s.srvMessChan = make(chan *gts.Message)

	s.w = watcher.NewWatcher()
	s.stor = stor
	s.htsrv = fselink.NewServer()

	return s
}
