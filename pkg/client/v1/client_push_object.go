package v1

import (
	"fmt"
	"path/filepath"

	"github.com/lazybark/go-cloud-sync/pkg/fse"
)

func (c *FSWClient) PushObject(obj fse.FSObject) {
	pathFullUnescaped := filepath.Join(c.fp.GetPathUnescaped(obj))

	if c.IsInActionBuffer(pathFullUnescaped) {
		return
	}

	c.AddToActionBuffer(pathFullUnescaped)
	defer c.RemoveFromActionBuffer(pathFullUnescaped)

	fmt.Println("UPLOADING:", pathFullUnescaped)

	file, err := c.fp.OpenToRead(pathFullUnescaped)
	if err != nil {
		c.extErc <- fmt.Errorf("[PUSH TO SERVER]%w", err)
		return
	}
	defer file.Close()
	err = c.link.PushObject(obj, file)
	if err != nil {
		c.extErc <- fmt.Errorf("[PUSH TO SERVER]%w", err)
		return
	}
}
