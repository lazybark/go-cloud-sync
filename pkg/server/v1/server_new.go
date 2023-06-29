package v1

import (
	"github.com/lazybark/go-cloud-sync/pkg/fse"
	"github.com/lazybark/go-cloud-sync/pkg/fselink"
	"github.com/lazybark/go-cloud-sync/pkg/storage"
	"github.com/lazybark/go-cloud-sync/pkg/watcher"
	gts "github.com/lazybark/go-tls-server/v2/server"
)

func NewServer(stor storage.IStorage) *FSWServer {
	s := &FSWServer{}

	s.evc = make(chan (fse.FSEvent))
	s.erc = make(chan (error)) //Not used for now
	s.srvConnChan = make(chan *gts.Connection)
	s.srvErrChan = make(chan error)
	s.srvMessChan = make(chan *gts.Message)

	s.w = watcher.NewWatcher()
	s.stor = stor
	s.htsrv = fselink.NewServer()

	return s
}
