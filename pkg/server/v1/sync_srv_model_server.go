package v1

import (
	"sync"

	"github.com/lazybark/go-cloud-sync/pkg/fp"
	"github.com/lazybark/go-cloud-sync/pkg/storage"
	fselink "github.com/lazybark/go-cloud-sync/pkg/synclink"
	gts "github.com/lazybark/go-tls-server/v3/server"
)

type FSWServer struct {
	conf ServerConfig

	fp fp.FileprocessorV1

	connPool map[string]*syncConnection
	//connPoolMutex controls connPool
	connPoolMutex sync.RWMutex

	extErc chan (error)
	erc    chan (error)

	srvMessChan chan *gts.Message
	srvErrChan  chan error
	srvConnChan chan *gts.Connection

	stor storage.IServerStorageV1

	htsrv fselink.ISyncLinkServerV1

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
