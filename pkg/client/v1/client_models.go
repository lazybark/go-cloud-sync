package v1

import (
	"sync"

	"github.com/lazybark/go-cloud-sync/pkg/fp"
	"github.com/lazybark/go-cloud-sync/pkg/fse"
	"github.com/lazybark/go-cloud-sync/pkg/fselink"
	"github.com/lazybark/go-cloud-sync/pkg/watcher"
)

type FSWClient struct {
	w                  watcher.IFilesystemWatcher
	extEvChannel       chan (fse.FSEvent)
	evc                chan (fse.FSEvent)
	extErc             chan (error)
	erc                chan (error)
	fp                 fp.Fileprocessor
	ActionsBuffer      map[string]bool
	ActionsBufferMutex sync.RWMutex

	link fselink.SyncLinkClient

	cfg ClientConfig
}

type ClientConfig struct {
	//root is the full path to directory where Watcher will watch for events (subdirs included)
	Root     string
	CacheDir string
}
