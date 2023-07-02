package server

import "fmt"

func (s *FSWServer) processSyncStart(c *SyncConnection) {
	fmt.Println("processSyncStart")
	c.sendEvents = true
}
