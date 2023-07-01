package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/lazybark/go-cloud-sync/configs"
	"github.com/lazybark/go-cloud-sync/pkg/client"
	"github.com/lazybark/go-cloud-sync/pkg/fselink/v1/proto"
)

//Basic client algorithm:
//Init with cofing
//Connect to server
//Start watching for changes
//Rescan filesystem
//Find what's changed
//Send changed metadata to server
//Request changed metadata from server
//Download & upload all changes
//Watch and upload all changed objects
//Wait for server notifications about changes on other clients and request changed objects

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

	w := client.NewClientV1(`D:\client_cache`, filepath.Join(filepath.Split(cfg.FS.Root)))
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
