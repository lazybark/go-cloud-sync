package server

import (
	"sync"

	"github.com/lazybark/go-cloud-sync/pkg/cloud/v1/fp"
	"github.com/lazybark/go-cloud-sync/pkg/storage"
	"github.com/lazybark/go-cloud-sync/pkg/synclink"
	srv "github.com/lazybark/go-tls-server/v3/server"
)

type FSWServer struct {
	conf ServerConfig

	fp fp.FileprocessorV1

	connPool map[string]*syncConnection
	//connPoolMutex controls connPool
	connPoolMutex sync.RWMutex

	extErc chan (error)
	erc    chan (error)

	srvMessChan chan *srv.Message
	srvErrChan  chan error
	srvConnChan chan *srv.Connection

	stor storage.IServerStorageV1

	htsrv synclink.ISyncLinkServerV1

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
