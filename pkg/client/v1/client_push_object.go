package v1

import (
	"fmt"
	"path/filepath"

	"github.com/lazybark/go-cloud-sync/pkg/fse"
)

func (c *FSWClient) PushObject(obj fse.FSObject) {
	pathFullUnescaped := filepath.Join(c.fp.GetPathUnescaped(obj))
	c.AddToActionBuffer(pathFullUnescaped)
	defer c.RemoveFromActionBuffer(pathFullUnescaped)
	file, err := c.fp.OpenToRead(pathFullUnescaped)
	if err != nil {
		c.extErc <- fmt.Errorf("[DOWNLOAD TO CACHE]: %w", err)
		return
	}
	err = c.link.PushObject(obj, file)
	if err != nil {
		c.extErc <- fmt.Errorf("[DOWNLOAD TO CACHE]: %w", err)
		return
	}
}