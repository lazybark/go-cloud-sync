package v1

import (
	"fmt"
	"time"
)

// rescanRoutine checks filesystem once per freq to find differences with current info in DB.
// Each difference is treated exactly as a FS event: processed & sent ot server.
func (c *FSWClient) resyncRoutine(freq time.Duration) {
	time.Sleep(freq)
	for {
		c.resyncOnce()
		fmt.Println("RESCANNED")
		time.Sleep(freq)
	}
}
