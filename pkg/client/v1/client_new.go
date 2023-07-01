package v1

import (
	"log"

	"github.com/lazybark/go-cloud-sync/pkg/fp"
	"github.com/lazybark/go-cloud-sync/pkg/fse"
	"github.com/lazybark/go-cloud-sync/pkg/fselink/v1/client"
	"github.com/lazybark/go-cloud-sync/pkg/watcher"
)

func NewClient(cacheDir, root string) *FSWClient {
	c := &FSWClient{
		cfg: ClientConfig{Root: root},
	}
	c.evc = make(chan (fse.FSEvent))
	c.erc = make(chan (error))

	c.w = watcher.NewWatcher()
	c.fp = fp.NewFPv1(",", root, cacheDir)
	link, err := client.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	c.link = link

	return c
}
