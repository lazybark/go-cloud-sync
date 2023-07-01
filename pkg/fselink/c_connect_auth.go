package fselink

import (
	"fmt"

	proto "github.com/lazybark/go-cloud-sync/pkg/fselink/proto/v1"
)

func (sc *SyncClient) ConnectAndAuth() error {
	err := sc.c.DialTo(sc.serverAddr, sc.serverPort, `certs/cert.pem`)
	if err != nil {
		return fmt.Errorf("[ConnectAndAuth]%w", err)
	}

	//At first we send credentials to get auth key
	err = sc.SendSyncMessage(proto.Credentials{Login: sc.login, Password: sc.pwd}, proto.MessageTypeAuthReq)
	if err != nil {
		return fmt.Errorf("[ConnectAndAuth]%w", err)
	}

	maa, err := sc.Await()
	if err != nil {
		return fmt.Errorf("[ConnectAndAuth]%w", err)
	}
	if maa.Type == proto.MessageTypeError {
		a, err := maa.ReadError()
		if err != nil {
			return fmt.Errorf("[ConnectAndAuth]%w", err)
		}
		return fmt.Errorf("sync error #%d: %s", a.ErrorCode, a.Error)
	} else if maa.Type == proto.MessageTypeAuthAns {
		a, err := maa.ReadAuthAnswer()
		if err != nil {
			return fmt.Errorf("[ConnectAndAuth]%w", err)
		}
		if !a.Success {
			return fmt.Errorf("[ConnectAndAuth]%w", err)
		}
		sc.akey = a.AuthKey
		return nil
	}

	return fmt.Errorf("[GetObjList] unexpected answer type '%s'", maa.Type)
}
