package v1

import (
	"fmt"
	"path/filepath"

	"github.com/lazybark/go-cloud-sync/pkg/fse"
)

func (c *FSWClient) processFilesystemEvent(event fse.FSEvent) {
	if _, buffered := c.ActionsBuffer[filepath.Join(event.Object.Path, event.Object.Name)]; buffered {
		return
	}
	//c.extEvChannel <- event
	//Process event with storage
	if event.Action == fse.Create || event.Action == fse.Write {
		obj, err := c.fp.ProcessObject(event.Object, true)
		if err != nil {
			c.extErc <- fmt.Errorf("[PROCESS EVENT][%s] processing error: %w", event.Action.String(), err)
			return
		}
		go c.PushObject(obj)
		//We skip objects that have no access to avoid syncing temp junks of apps
		if obj.Hash == "" && !obj.IsDir {
			return
		}
		if event.Action == fse.Create && obj.IsDir {
			err = c.w.Add(c.fp.GetPathUnescaped(event.Object))
			if err != nil {
				c.extErc <- fmt.Errorf("[PROCESS EVENT][%s] adding watcher failed: %w", event.Action.String(), err)
				return
			}
		}
	} else if event.Action == fse.Remove {
		c.w.RemoveIfExists(event.Object.Name)
		dir, name, err := c.fp.ConvertPathName(event.Object)
		if err != nil {
			c.extErc <- fmt.Errorf("[PROCESS EVENT][%s]: %w", event.Action.String(), err)
			return
		}
		event.Object.Path = dir
		event.Object.Name = name
		err = c.link.DeleteObject(event.Object)
		if err != nil {
			c.extErc <- fmt.Errorf("[PROCESS EVENT][%s]: %w", event.Action.String(), err)
			return
		}
	}

}
