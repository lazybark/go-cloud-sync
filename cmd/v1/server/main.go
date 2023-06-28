package main

import (
	"fmt"
	"log"

	"github.com/lazybark/go-cloud-sync/pkg/fse"
	"github.com/lazybark/go-cloud-sync/pkg/server"
	"github.com/lazybark/go-cloud-sync/pkg/storage/sqlitestorage"
)

func main() {
	/*cfg, err := configs.MakeConfig()
	if err != nil {
		log.Fatal(err)
	}*/

	evc := make(chan (fse.FSEvent))
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

	sqstor, err := sqlitestorage.NewSQLite("", ",")
	if err != nil {
		log.Fatal(err)
	}

	w := server.NewServerV1(sqstor)
	err = w.Init(`D:\filesystem_root_server`, `D:\server_cache`, `localhost`, `5555`, ",", evc, erc)
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
