package fselink

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/lazybark/go-cloud-sync/pkg/fse"
	"github.com/lazybark/go-tls-server/v2/client"
)

func NewClient() (*SyncClient, error) {
	c := &SyncClient{}

	return c, nil
}

func (sc *SyncClient) GetObjList() (l []fse.FSObject, err error) {
	err = SendSyncMessage(sc, nil, MessageTypeFullSyncRequest)
	if err != nil {
		err = fmt.Errorf("[GetObjList] %w", err)
		return
	}
	var maa ExchangeMessage
	err = AwaitAnswer(sc, &maa)
	if err != nil {
		err = fmt.Errorf("[GetObjList] %w", err)
		return
	}
	if maa.Type == MessageTypeError {
		var se MessageError
		err = UnpackMessage(maa, MessageTypeError, &se)
		if err != nil {
			return l, fmt.Errorf("[GetObjList] %w", err)
		}
		return l, fmt.Errorf("sync error #%d: %s", se.ErrorCode, se.Error)
	} else if maa.Type == MessageTypeFullSyncReply {
		var sm MessageFullSyncReply
		err = UnpackMessage(maa, MessageTypeFullSyncReply, &sm)
		if err != nil {
			return l, fmt.Errorf("[GetObjList] %w", err)
		}
		l = sm.Objects

		return
	}

	return l, fmt.Errorf("[GetObjList] unexpected answer type '%s'", maa.Type)
}

func (sc *SyncClient) PushObject(obj fse.FSObject, fileData *os.File) (err error) {
	link, err := NewClient()
	if err != nil {
		return fmt.Errorf("[PushObject]%w", err)
	}

	err = link.Init(sc.serverPort, sc.serverAddr, sc.login, sc.pwd)
	if err != nil {
		return fmt.Errorf("[PushObject]%w", err)
	}
	defer link.c.Close()

	err = link.ConnectAndAuth()
	if err != nil {
		return fmt.Errorf("[PushObject]%w", err)
	}

	err = SendSyncMessage(link, MessageGetFile{Object: obj}, MessageTypePushFile)
	if err != nil {
		return fmt.Errorf("[PushObject]%w", err)
	}

	var maa ExchangeMessage
	err = AwaitAnswer(link, &maa)
	if err != nil {
		return fmt.Errorf("[PushObject]%w", err)
	}
	if maa.Type == MessageTypeError {
		var se MessageError
		err := UnpackMessage(maa, MessageTypeError, &se)
		if err != nil {
			return fmt.Errorf("[PushObject]%w", err)
		}
		return fmt.Errorf("sync error #%d: %s", se.ErrorCode, se.Error)
	} else if maa.Type == MessageTypePeerReady {
		fmt.Println("PEER READY")

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

			err = SendSyncMessage(link, MessageFilePart{Payload: buf[:n]}, MessageTypeFileParts)
			if err != nil {
				return fmt.Errorf("[PushObject]%w", err)
			}
		}
		fileData.Close()

		err = SendSyncMessage(link, nil, MessageTypeFileEnd)
		if err != nil {
			return fmt.Errorf("[PushObject]%w", err)
		}
		fmt.Println("CLIENT: SENT FILE")

		return
	}
	return fmt.Errorf("[PushObject] unexpected answer type '%s'", maa.Type)
}

func (sc *SyncClient) DownloadObject(obj fse.FSObject, destFile *os.File) (err error) {
	link, err := NewClient()
	if err != nil {
		return fmt.Errorf("[DownloadObject]%w", err)
	}

	err = link.Init(sc.serverPort, sc.serverAddr, sc.login, sc.pwd)
	if err != nil {
		return fmt.Errorf("[DownloadObject]%w", err)
	}
	defer link.c.Close()

	err = link.ConnectAndAuth()
	if err != nil {
		return fmt.Errorf("[DownloadObject]%w", err)
	}

	err = SendSyncMessage(link, MessageGetFile{Object: obj}, MessageTypeGetFile)
	if err != nil {
		return fmt.Errorf("[DownloadObject]%w", err)
	}

	var maa ExchangeMessage
	for {
		err = AwaitAnswer(link, &maa)
		if err != nil {
			return fmt.Errorf("[DownloadObject]%w", err)
		}
		if maa.Type == MessageTypeError {
			var se MessageError
			err := UnpackMessage(maa, MessageTypeError, &se)
			if err != nil {
				return fmt.Errorf("[ConnectAndAuth]%w", err)
			}
			return fmt.Errorf("sync error #%d: %s", se.ErrorCode, se.Error)
		} else if maa.Type == MessageTypeFileParts {
			var m MessageFilePart
			err := UnpackMessage(maa, MessageTypeFileParts, &m)
			if err != nil {
				return fmt.Errorf("[ConnectAndAuth]%w", err)
			}
			_, err = destFile.Write(m.Payload)
			if err != nil {
				return fmt.Errorf("[DownloadObject]%w", err)
			}

		} else if maa.Type == MessageTypeFileEnd {
			return nil
		} else {
			return fmt.Errorf("[DownloadObject] unexpected answer type '%s'", maa.Type)
		}
	}
}

func (sc *SyncClient) Init(port int, addr, login, pwd string) error {

	conf := client.Config{SuppressErrors: false, MessageTerminator: '\n'}
	c := client.New(&conf)
	sc.c = c
	sc.serverAddr = addr
	sc.serverPort = port
	sc.login = login
	sc.pwd = pwd

	fmt.Println("Connecting to server")
	err := sc.ConnectAndAuth()
	if err != nil {
		return fmt.Errorf("[SyncClient][Init]%w", err)
	}
	fmt.Println("Got auth key:", sc.akey)

	return nil
}

func (sc *SyncClient) AwaitAnswer() (*client.Message, error) {
	ans := <-sc.c.MessageChan
	//pissble error will be added here later to prevent endless loops
	return ans, nil
}

func (sc *SyncClient) compileMessageBody(t ExchangeMessageType) ExchangeMessage {
	return ExchangeMessage{Type: t, AuthKey: sc.akey, Payload: []byte{}}
}
