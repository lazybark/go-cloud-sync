package client

import "fmt"

func (c *FSWClient) IsInActionBuffer(object string) bool {
	c.ActionsBufferMutex.Lock()
	_, yes := c.ActionsBuffer[object]
	c.ActionsBufferMutex.Unlock()
	return yes
}

func (c *FSWClient) AddToActionBuffer(object string) {
	fmt.Println("buffer", object)
	c.ActionsBufferMutex.Lock()
	c.ActionsBuffer[object] = true
	c.ActionsBufferMutex.Unlock()
}

func (c *FSWClient) RemoveFromActionBuffer(object string) {
	c.ActionsBufferMutex.Lock()
	delete(c.ActionsBuffer, object)
	c.ActionsBufferMutex.Unlock()
}
