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

	sqstor, err := sqlitestorage.NewSQLite("", ",")
	if err != nil {
		log.Fatal(err)
	}

	w := server.NewServerV1(sqstor)
	w.Init(`D:\filesystem_root_server`, `localhost`, `5555`, ",", evc, erc)
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
