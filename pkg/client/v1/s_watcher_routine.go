package v1

import (
	"fmt"

	"github.com/lazybark/go-cloud-sync/pkg/fse"
)

func (c *FSWClient) watcherRoutine() {
	done := make(chan bool)
	go func() {
		defer close(done)

		for {
			select {
			case event, ok := <-c.evc:
				if !ok {
					return
				}
				c.addFsEvent(event)

			case err, ok := <-c.erc:
				if !ok {
					return
				}
				c.extErc <- fmt.Errorf("[FSWATCHER][WATCH] fs watcher error: %w", err)
			}
		}
	}()

	<-done
}

func (c *FSWClient) addFsEvent(event fse.FSEvent) {
	//Process event with storage
	if event.Action == fse.Create || event.Action == fse.Write {
		obj, err := c.fp.ProcessObject(event.Object, true)
		if err != nil {
			c.extErc <- fmt.Errorf("[FSWATCHER][WATCH][%s] processing error: %w", event.Action.String(), err)
		}
		//We skip objects that have no access to avoid syncing temp junks of apps
		if obj.Hash == "" && !obj.IsDir {
			return
		}
		/*err = c.db.AddOrUpdateObject(obj)
		if err != nil {
			c.extErc <- fmt.Errorf("[FSWATCHER][WATCH][%s] processing error: %w", event.Action.String(), err)
		}*/
		if event.Action == fse.Create && obj.IsDir {
			err = c.w.Add(c.fp.GetPathUnescaped(event.Object))
			if err != nil {
				c.extErc <- fmt.Errorf("[FSWATCHER][WATCH][%s] adding watcher failed: %w", event.Action.String(), err)
			}
		}
	} else if event.Action == fse.Remove {
		c.w.RemoveIfExists(event.Object.Name)
		/*obj := event.Object
		path, name, err := c.fp.ConvertPathName(obj)
		if err != nil {
			c.extErc <- fmt.Errorf("[FSWATCHER][WATCH][ConvertPathName] processing error: %w", err)
		}
		obj.Path = path
		obj.Name = name

		err = c.db.RemoveObject(obj, true)
		if err != nil {
			c.extErc <- fmt.Errorf("[FSWATCHER][WATCH][RemoveObject] processing error: %w", err)
		}*/
	}
	//Send event to external code
	c.extEvChannel <- event

	//Now send to the server
	err := c.link.SendEvent(event)
	if err != nil {
		c.extErc <- fmt.Errorf("[FSWATCHER][WATCH][%s] error notifying server: %w", event.Action.String(), err)
	}
}
