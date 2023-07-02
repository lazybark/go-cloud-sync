package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/lazybark/go-cloud-sync/configs"
	"github.com/lazybark/go-cloud-sync/pkg/cloud"
	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/proto"
)

func main() {
	cfg, err := configs.MakeConfig()
	if err != nil {
		log.Fatal(err)
	}

	evc := make(chan (proto.FSEvent))
	erc := make(chan error)

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

	cfg.FS.Root = `D:\filesystem_root2`
	w := cloud.NewClientV1(`D:\client_cache2`, filepath.Join(filepath.Split(cfg.FS.Root)))
	err = w.Init(evc, erc, "login", "pwd")
	if err != nil {
		log.Fatal(err)
	}
	err = w.Start()
	if err != nil {
		log.Fatal(err)
	}

	//Just endless cycle for now
	for {

	}
}
