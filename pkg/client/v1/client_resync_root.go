package v1

import (
	"fmt"
	"os"
	"path/filepath"
)

func (c *FSWClient) resyncOnce() {
	local, err := c.GetLocalObjects()
	if err != nil {
		c.extErc <- fmt.Errorf("[SCAN DIR] getting local objects: %w", err)
		return
	}
	/*fmt.Println("_________________________")
	for _, oo := range local {
		fmt.Println(oo)
	}
	fmt.Println("_________________________")*/
	objsOnServer, err := c.GetServerObjList()
	if err != nil {
		c.extErc <- fmt.Errorf("[SCAN DIR] Getting server objects: %w", err)
		return
	}
	/*fmt.Println("++++++++++++++++++")
	for _, oo := range objsOnServer {
		fmt.Println(oo)
	}
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
			pathFullUnescaped := filepath.Join(c.fp.GetPathUnescaped(o))
			if err := os.MkdirAll(pathFullUnescaped, os.ModePerm); err != nil {
				if err != nil {
					c.extErc <- fmt.Errorf("[DOWNLOAD TO CACHE]: %w", err)
				}
				continue
			}
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
