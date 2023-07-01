package client

import (
	"fmt"
	"os"

	"github.com/lazybark/go-cloud-sync/pkg/fselink/v1/proto"
)

func (sc *LinkClient) DownloadObject(obj proto.FSObject, destFile *os.File) (err error) {
	link, err := NewClient()
	if err != nil {
		return fmt.Errorf("[DownloadObject]%w", err)
	}

	err = link.Init(sc.serverPort, sc.serverAddr, sc.login, sc.pwd)
	if err != nil {
		return fmt.Errorf("[DownloadObject]%w", err)
	}
	defer link.c.Close()

	link.SetAuthKey(sc.akey)

	err = link.SendSyncMessage(proto.MessageObject{Object: obj}, proto.MessageTypeGetFile)
	if err != nil {
		return fmt.Errorf("[DownloadObject]%w", err)
	}

	var maa proto.ExchangeMessage
	for {
		maa, err = link.Await()
		if err != nil {
			return fmt.Errorf("[DownloadObject]%w", err)
		}
		if maa.Type == proto.MessageTypeError {
			a, err := maa.ReadError()
			if err != nil {
				return fmt.Errorf("[DownloadObject]%w", err)
			}
			return fmt.Errorf("sync error #%d: %s", a.ErrorCode, a.Error)
		} else if maa.Type == proto.MessageTypeFileParts {
			a, err := maa.ReadFilePart()
			if err != nil {
				return fmt.Errorf("[DownloadObject]%w", err)
			}
			_, err = destFile.Write(a.Payload)
			if err != nil {
				return fmt.Errorf("[DownloadObject]%w", err)
			}

		} else if maa.Type == proto.MessageTypeFileEnd {
			return nil
		} else {
			return fmt.Errorf("[DownloadObject] unexpected answer type '%s'", maa.Type)
		}
	}
}
