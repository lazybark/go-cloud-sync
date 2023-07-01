package v1

import (
	"sync"

	"github.com/lazybark/go-cloud-sync/pkg/fp"
	"github.com/lazybark/go-cloud-sync/pkg/fselink"
	"github.com/lazybark/go-cloud-sync/pkg/fselink/v1/proto"
	"github.com/lazybark/go-cloud-sync/pkg/watcher"
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

	link fselink.ISyncLinkClientV1

	cfg ClientConfig
}

type ClientConfig struct {
	//root is the full path to directory where Watcher will watch for events (subdirs included)
	Root     string
	CacheDir string
}
