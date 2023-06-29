package v1

import (
	"time"
)

// rescanRoutine checks filesystem once per freq to find differences with current info in DB.
// Each difference is treated exactly as a FS event: processed & sent ot server.
func (c *FSWClient) rescanRoutine(freq time.Duration) {
	time.Sleep(freq)
	for {
		c.rescanOnce()
		time.Sleep(freq)
	}
}