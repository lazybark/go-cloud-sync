package client

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/proto"
)

func (sc *LinkClient) PushObject(obj proto.FSObject, fileData *os.File) (err error) {
	link, err := NewClient()
	if err != nil {
		return fmt.Errorf("[PushObject]%w", err)
	}

	err = link.Init(sc.serverPort, sc.serverAddr, sc.login, sc.pwd)
	if err != nil {
		return fmt.Errorf("[PushObject]%w", err)
	}
	defer link.c.Close()

	link.SetAuthKey(sc.akey)

	err = link.SendSyncMessage(proto.MessageObject{Object: obj}, proto.MessageTypePushFile)
	if err != nil {
		return fmt.Errorf("[PushObject]%w", err)
	}

	maa, err := link.Await()
	if err != nil {
		return fmt.Errorf("[PushObject]%w", err)
	}
	if maa.Type == proto.MessageTypeError {
		a, err := maa.ReadError()
		if err != nil {
			return fmt.Errorf("[PushObject]%w", err)
		}
		return fmt.Errorf("sync error #%d: %s", a.ErrorCode, a.Error)
	} else if maa.Type == proto.MessageTypePeerReady {
		// TLS record size can be up to 16KB but some extra bytes may apply
		// https://hpbn.co/transport-layer-security-tls/#optimize-tls-record-size
		buf := make([]byte, 15360)
		n := 0

		r := bufio.NewReader(fileData)

		for {
			n, err = r.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("[PushObject]%w", err)
			}

			err = link.SendSyncMessage(proto.MessageFilePart{Payload: buf[:n]}, proto.MessageTypeFileParts)
			if err != nil {
				return fmt.Errorf("[PushObject]%w", err)
			}
		}
		fileData.Close()

		err = link.SendSyncMessage(nil, proto.MessageTypeFileEnd)
		if err != nil {
			return fmt.Errorf("[PushObject]%w", err)
		}
		fmt.Println("CLIENT: SENT FILE")

		return
	} else if maa.Type == proto.MessageTypeClose {
		//We try to close only after server says that we can do it (no more answers expected)
		return
	} else {
		return fmt.Errorf("[PushObject] unexpected answer type '%s'", maa.Type)
	}
}
