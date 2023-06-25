package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/lazybark/go-cloud-sync/configs"
	"github.com/lazybark/go-cloud-sync/pkg/client"
	"github.com/lazybark/go-cloud-sync/pkg/fse"
	"github.com/lazybark/go-cloud-sync/pkg/storage/sqlitestorage"
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

	evc := make(chan (fse.FSEvent))
	erc := make(chan error)

	sqstor, err := sqlitestorage.NewSQLite("", ",")
	if err != nil {
		log.Fatal(err)
	}

	w := client.NewClientV1(sqstor, filepath.Join(filepath.Split(cfg.FS.Root)))
	w.Init(evc, erc)
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

	//Just endless cycle for now
	for {

	}
}