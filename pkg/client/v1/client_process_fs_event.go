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
	//Process event with storage
	if event.Action == fse.Create || event.Action == fse.Write {
		obj, err := c.fp.ProcessObject(event.Object, true)
		if err != nil {
			c.extErc <- fmt.Errorf("[PROCESS EVENT][%s] processing error: %w", event.Action.String(), err)
		}
		//We skip objects that have no access to avoid syncing temp junks of apps
		if obj.Hash == "" && !obj.IsDir {
			return
		}
		if event.Action == fse.Create && obj.IsDir {
			err = c.w.Add(c.fp.GetPathUnescaped(event.Object))
			if err != nil {
				c.extErc <- fmt.Errorf("[PROCESS EVENT][%s] adding watcher failed: %w", event.Action.String(), err)
			}
		}
	} else if event.Action == fse.Remove {
		c.w.RemoveIfExists(event.Object.Name)
	}
	//Send event to external code
	c.extEvChannel <- event

	//Now send to the server
	err := c.link.SendEvent(event)
	if err != nil {
		c.extErc <- fmt.Errorf("[PROCESS EVENT][%s] error notifying server: %w", event.Action.String(), err)
	}
}
