package watcher

import "github.com/lazybark/go-cloud-sync/pkg/synclink/v1/proto"

//IFilesystemWatcher represents watcher that uses event (evc) and error (erc) channels to report all
//changes in specified dir and its subdirs.
type IFilesystemWatcherV1 interface {
	//Init sets initial parameters for fs watcher
	Init(root string, evc chan (proto.FSEvent), erc chan (error)) error

	//Start launches the watcher routine or returns an error
	Start() error

	//Stop stops the watcher routine or returns an error. It also should close event & error channels,
	//which means new Start() will need new Init() with new channels
	Stop() error

	Add(dir string) error

	RemoveIfExists(dir string)
}
