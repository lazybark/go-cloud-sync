package main

import (
	"fmt"

	"github.com/lazybark/go-cloud-sync/configs"
	"github.com/lazybark/go-cloud-sync/pkg/fse"
	"github.com/lazybark/go-cloud-sync/pkg/server"
	"github.com/lazybark/go-cloud-sync/pkg/storage/sqlitestorage"
)

func main() {
	cfg, err := configs.MakeConfig()
	if err != nil {
		fmt.Println(err)
	}

	evc := make(chan (fse.FSEvent))
	erc := make(chan error)

	w := server.NewServerV1(sqlitestorage.NewSQLite(""))
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

	//Just endless cycle for now
	for {

	}
}
