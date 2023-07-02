package client

import (
	"bufio"
	"fmt"
	"io"
	"path/filepath"

	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/proto"
)

func (c *FSWClient) PushObject(obj proto.FSObject) {
	pathFullUnescaped := filepath.Join(c.fp.GetPathUnescaped(obj))

	if c.IsInActionBuffer(pathFullUnescaped) {
		return
	}

	c.AddToActionBuffer(pathFullUnescaped)
	defer c.RemoveFromActionBuffer(pathFullUnescaped)

	fmt.Println("UPLOADING:", pathFullUnescaped)

	file, err := c.fp.OpenToRead(pathFullUnescaped)
	if err != nil {
		c.extErc <- fmt.Errorf("[PushObject]%w", err)
		return
	}
	defer file.Close()

	newConn, err := c.link.NewConnectionInSession()
	if err != nil {
		c.extErc <- fmt.Errorf("[PushObject]: %w", err)
		return
	}

	err = newConn.SendSyncMessage(proto.MessageObject{Object: obj}, proto.MessageTypePushFile)
	if err != nil {
		c.extErc <- fmt.Errorf("[PushObject]: %w", err)
		return
	}
	maa, err := newConn.Await()
	if err != nil {
		c.extErc <- fmt.Errorf("[PushObject]: %w", err)
		return
	}
	if maa.Type == proto.MessageTypeError {
		a, err := maa.ReadError()
		if err != nil {
			c.extErc <- fmt.Errorf("[PushObject]: %w", err)
			return
		}
		c.extErc <- fmt.Errorf("sync error #%d: %s", a.ErrorCode, a.Error)
	} else if maa.Type == proto.MessageTypePeerReady {
		// TLS record size can be up to 16KB but some extra bytes may apply
		// https://hpbn.co/transport-layer-security-tls/#optimize-tls-record-size
		buf := make([]byte, 15360)
		n := 0

		r := bufio.NewReader(file)

		for {
			n, err = r.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				c.extErc <- fmt.Errorf("[PushObject]: %w", err)
				return
			}

			err = newConn.SendSyncMessage(proto.MessageFilePart{Payload: buf[:n]}, proto.MessageTypeFileParts)
			if err != nil {
				c.extErc <- fmt.Errorf("[PushObject]: %w", err)
				return
			}
		}

		err = newConn.SendSyncMessage(nil, proto.MessageTypeFileEnd)
		if err != nil {
			c.extErc <- fmt.Errorf("[PushObject]: %w", err)
			return
		}
		fmt.Println("CLIENT: SENT FILE")

		return
	} else if maa.Type == proto.MessageTypeClose {
		//We try to close only after server says that we can do it (no more answers expected)
		err = newConn.Close()
		if err != nil {
			c.extErc <- fmt.Errorf("[PushObject]: %w", err)
			return
		}
		return
	} else {
		c.extErc <- fmt.Errorf("[PushObject] unexpected answer type '%s'", maa.Type)
		return
	}
}
