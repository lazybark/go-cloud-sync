package v1

import (
	"fmt"
	"time"

	"github.com/lazybark/go-cloud-sync/pkg/fse"
)

// rescanRoutine checks filesystem once per freq to find differences with current info in DB.
// Each difference is treated exactly as a FS event: processed & sent ot server.
func (c *FSWClient) rescanRoutine(freq time.Duration) {
	for {
		c.rescanOnce()
		time.Sleep(freq)
	}
}

func (c *FSWClient) rescanOnce() {
	local, err := c.GetLocalObjects()
	if err != nil {
		c.extErc <- fmt.Errorf("[FSWATCHER][START] local objects: %w", err)
	}
	objsOnServer, err := c.GetServerObjList()
	if err != nil {
		c.extErc <- fmt.Errorf("[FSWATCHER][START][DiffListWithServer]: %w", err)
	}
	fmt.Println("++++++++++++++++++")
	fmt.Println(objsOnServer)
	fmt.Println("++++++++++++++++++")
	download, created, updated, err := c.GetDiffListWithServer(local, objsOnServer)
	if err != nil {
		c.extErc <- fmt.Errorf("[FSWATCHER][START] diff with server: %w", err)
	}
	fmt.Println("TO DOWNLOAD")
	for _, o := range download {
		fmt.Println(o.Path, o.Name)
		if o.IsDir {
			//We do not download whole dirs, only file by file
			continue
		}
		go c.DownloadObject(o)
	}
	fmt.Println("TO created")
	//Each obj in created & updated is treated as a new FS event
	for _, o := range created {
		fmt.Println(o.Path, o.Name)
		go c.addFsEvent(fse.FSEvent{Action: fse.Create, Object: o})
	}
	fmt.Println("TO updated")
	for _, o := range updated {
		fmt.Println(o.Path, o.Name)
		go c.addFsEvent(fse.FSEvent{Action: fse.Write, Object: o})
	}
}
