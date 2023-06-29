package v1

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lazybark/go-cloud-sync/pkg/fse"
)

func (c *FSWClient) DownloadObject(obj fse.FSObject) {
	pathUnescaped := filepath.Join(c.fp.UnescapePath(obj))
	pathFullUnescaped := filepath.Join(c.fp.GetPathUnescaped(obj))
	c.AddToActionBuffer(pathFullUnescaped)
	defer c.RemoveFromActionBuffer(pathFullUnescaped)
	//Create file in cache
	file, err := c.fp.CreateFileInCache()
	if err != nil {
		c.extErc <- fmt.Errorf("[DOWNLOAD TO CACHE]: %w", err)
		return
	}
	fmt.Println(file.Name())
	fmt.Println(pathUnescaped)
	fmt.Println(pathFullUnescaped)
	//Downloading object to cache
	err = c.link.DownloadObject(obj, file)
	if err != nil {
		c.extErc <- fmt.Errorf("[DOWNLOAD TO CACHE]: %w", err)
		file.Close()
		err = c.fp.DeleteFileInCache(file.Name())
		if err != nil {
			c.extErc <- fmt.Errorf("[DOWNLOAD TO CACHE]: %w", err)
		}
		return
	}
	file.Close()
	//Moving from cache to real place
	if err := os.MkdirAll(pathUnescaped, os.ModePerm); err != nil {
		c.extErc <- fmt.Errorf("[DOWNLOAD TO CACHE]: %w", err)
		err = c.fp.DeleteFileInCache(file.Name())
		if err != nil {
			c.extErc <- fmt.Errorf("[DOWNLOAD TO CACHE]: %w", err)
		}
		return
	}
	err = os.Rename(file.Name(), pathFullUnescaped)
	if err != nil {
		c.extErc <- fmt.Errorf("[DOWNLOAD TO CACHE]: %w", err)
		err = c.fp.DeleteFileInCache(file.Name())
		if err != nil {
			c.extErc <- fmt.Errorf("[DOWNLOAD TO CACHE]: %w", err)
		}
		return
	}
}

func (c *FSWClient) AddToActionBuffer(object string) {
	c.ActionsBufferMutex.Lock()
	c.ActionsBuffer[object] = true
	c.ActionsBufferMutex.Unlock()
}

func (c *FSWClient) RemoveFromActionBuffer(object string) {
	c.ActionsBufferMutex.Lock()
	delete(c.ActionsBuffer, object)
	c.ActionsBufferMutex.Unlock()
}
