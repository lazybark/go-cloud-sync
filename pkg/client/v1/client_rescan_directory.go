package v1

import (
	"fmt"
)

func (c *FSWClient) rescanOnce() {
	local, err := c.GetLocalObjects()
	if err != nil {
		c.extErc <- fmt.Errorf("[SCAN DIR] getting local objects: %w", err)
		return
	}
	objsOnServer, err := c.GetServerObjList()
	if err != nil {
		c.extErc <- fmt.Errorf("[SCAN DIR] Getting server objects: %w", err)
		return
	}
	/*	fmt.Println("++++++++++++++++++")
		fmt.Println(objsOnServer)
		fmt.Println("++++++++++++++++++")*/
	download, created, updated, err := c.GetDiffListWithServer(local, objsOnServer)
	if err != nil {
		c.extErc <- fmt.Errorf("[SCAN DIR] checking for differences: %w", err)
		return
	}
	//fmt.Println("TO DOWNLOAD")
	for _, o := range download {
		//fmt.Println(o.Path, o.Name)
		if o.IsDir {
			//We do not download whole dirs, only file by file
			continue
		}
		go c.DownloadObject(o)
	}
	//fmt.Println("TO created")
	//Each obj in created & updated is treated as a new FS event
	for _, o := range created {
		go c.PushObject(o)
		//fmt.Println(o.Path, o.Name)
		//go c.processFilesystemEvent(fse.FSEvent{Action: fse.Create, Object: o})
	}
	//fmt.Println("TO updated")
	for _, o := range updated {
		go c.PushObject(o)
		//fmt.Println(o.Path, o.Name)
		//go c.processFilesystemEvent(fse.FSEvent{Action: fse.Write, Object: o})
	}
}
