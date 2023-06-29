package v1

import (
	"github.com/lazybark/go-cloud-sync/pkg/fp"
	"github.com/lazybark/go-cloud-sync/pkg/fse"
	"github.com/lazybark/go-cloud-sync/pkg/fselink"
	"github.com/lazybark/go-cloud-sync/pkg/storage"
	"github.com/lazybark/go-cloud-sync/pkg/watcher"
	gts "github.com/lazybark/go-tls-server/v2/server"
)

type FSWServer struct {
	conf ServerConfig
	w    watcher.IFilesystemWatcher

	fp fp.Fileprocessor

	extEvc chan (fse.FSEvent)
	evc    chan (fse.FSEvent)
	extErc chan (error)
	erc    chan (error)

	srvMessChan chan *gts.Message
	srvErrChan  chan error
	srvConnChan chan *gts.Connection

	stor storage.IStorage

	htsrv fselink.FSEServerPool

	isActive bool
}

type ServerConfig struct {
	root      string
	cacheRoot string
	host      string
	port      string
	escSymbol string
}

type Client struct {
}