package main

import (
	"fmt"
	"log"

	"github.com/lazybark/go-cloud-sync/pkg/cloud"
	"github.com/lazybark/go-cloud-sync/pkg/storage/sqlitestorage"
)

func main() {
	/*cfg, err := configs.MakeConfig()
	if err != nil {
		log.Fatal(err)
	}*/

	erc := make(chan error)

	go func() {
		for err := range erc {
			fmt.Println(err)
		}
	}()

	sqstor, err := sqlitestorage.NewSQLite("", ",")
	if err != nil {
		log.Fatal(err)
	}

	w := cloud.NewServerV1(sqstor)
	err = w.Init(`D:\filesystem_root_server`, `D:\server_cache`, `localhost`, `5555`, ",", erc)
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
