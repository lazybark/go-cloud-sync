package server

import (
	"github.com/lazybark/go-cloud-sync/pkg/fse"
	v1 "github.com/lazybark/go-cloud-sync/pkg/server/v1"
	"github.com/lazybark/go-cloud-sync/pkg/storage"
)

type IFSWServer interface {
	//Init sets initial parameters for fs watcher
	Init(root string, evc chan (fse.FSEvent), erc chan (error)) error

	//Start launches the watcher routine or returns an error
	Start() error

	//Stop stops the watcher routine or returns an error. It also should close event & error channels,
	//which means new Start() will need new Init() with new channels
	Stop() error
}

func NewServerV1(stor storage.IStorage) IFSWServer {
	return v1.NewServer(stor)
}
