package v1

import (
	"sync"

	"github.com/lazybark/go-cloud-sync/pkg/fp"
	"github.com/lazybark/go-cloud-sync/pkg/fselink"
	"github.com/lazybark/go-cloud-sync/pkg/fselink/v1/proto"
	"github.com/lazybark/go-cloud-sync/pkg/storage"
	"github.com/lazybark/go-cloud-sync/pkg/watcher"
	gts "github.com/lazybark/go-tls-server/v3/server"
)

type FSWServer struct {
	conf ServerConfig
	w    watcher.IFilesystemWatcher

	fp fp.Fileprocessor

	connPool map[string]*syncConnection
	//connPoolMutex controls connPool
	connPoolMutex sync.RWMutex

	extEvc chan (proto.FSEvent)
	evc    chan (proto.FSEvent)
	extErc chan (error)
	erc    chan (error)

	srvMessChan chan *gts.Message
	srvErrChan  chan error
	srvConnChan chan *gts.Connection

	stor storage.IServerStorage

	htsrv fselink.ISyncLinkServer

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
