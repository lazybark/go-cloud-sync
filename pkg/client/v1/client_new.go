package v1

import (
	"log"

	"github.com/lazybark/go-cloud-sync/pkg/fp"
	"github.com/lazybark/go-cloud-sync/pkg/fse"
	"github.com/lazybark/go-cloud-sync/pkg/fselink"
	"github.com/lazybark/go-cloud-sync/pkg/storage"
	"github.com/lazybark/go-cloud-sync/pkg/watcher"
)

func NewClient(db storage.IStorage, cacheDir, root string) *FSWClient {
	c := &FSWClient{
		db:  db,
		cfg: ClientConfig{Root: root},
	}
	c.evc = make(chan (fse.FSEvent))
	c.erc = make(chan (error))

	c.w = watcher.NewWatcher()
	c.fp = fp.NewFPv1(",", root, cacheDir)
	link, err := fselink.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	c.link = link

	return c
}
