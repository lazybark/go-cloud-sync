package client

import (
	"sync"

	"github.com/lazybark/go-cloud-sync/pkg/cloud/v1/fp"
	"github.com/lazybark/go-cloud-sync/pkg/cloud/v1/watcher"
	"github.com/lazybark/go-cloud-sync/pkg/synclink"
	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/proto"
)

type FSWClient struct {
	w                  watcher.IFilesystemWatcherV1
	extEvChannel       chan (proto.FSEvent)
	evc                chan (proto.FSEvent)
	extErc             chan (error)
	erc                chan (error)
	fp                 fp.FileprocessorV1
	ActionsBuffer      map[string]bool
	ActionsBufferMutex sync.RWMutex

	link synclink.ISyncLinkClientV1

	cfg ClientConfig
}

type ClientConfig struct {
	//root is the full path to directory where Watcher will watch for events (subdirs included)
	Root     string
	CacheDir string
}
