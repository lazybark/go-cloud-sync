package main

import (
	"fmt"
	"time"

	"github.com/lazybark/go-cloud-sync/configs"
	"github.com/lazybark/go-cloud-sync/pkg/fse"
	"github.com/lazybark/go-cloud-sync/pkg/watcher"
)

func main() {
	cfg, err := configs.MakeConfig()
	if err != nil {
		fmt.Println(err)
	}

	evc := make(chan (fse.FSEvent))
	erc := make(chan error)

	w := watcher.NewV1()
	w.Init(cfg.FS.Root, evc, erc)
	w.Start()

	go func() {
		for event := range evc {
			fmt.Println(event)
		}
	}()

	go func() {
		for err := range erc {
			fmt.Println(err)
		}
	}()

	time.Sleep(time.Second * 10)

	w.Stop()

	fmt.Println("end")
}
