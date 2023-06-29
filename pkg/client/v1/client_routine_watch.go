package v1

import (
	"fmt"
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
				c.processFilesystemEvent(event)

			case err, ok := <-c.erc:
				if !ok {
					return
				}
				c.extErc <- fmt.Errorf("[WATCH EVENTS]%w", err)
			}
		}
	}()

	<-done
}
