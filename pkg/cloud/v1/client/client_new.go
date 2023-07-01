package client

import (
	"log"

	"github.com/lazybark/go-cloud-sync/pkg/cloud/v1/fp"
	"github.com/lazybark/go-cloud-sync/pkg/cloud/v1/watcher"
	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/client"
	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/proto"
)

func NewClient(cacheDir, root string) *FSWClient {
	c := &FSWClient{
		cfg: ClientConfig{Root: root},
	}
	c.evc = make(chan (proto.FSEvent))
	c.erc = make(chan (error))

	c.w = watcher.NewWatcherV1()
	c.fp = fp.NewFPv1(",", root, cacheDir)
	link, err := client.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	c.link = link

	return c
}
