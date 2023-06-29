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

	var maa proto.ExchangeMessage
	err = AwaitAnswer(sc, &maa)
	if err != nil {
		return fmt.Errorf("[ConnectAndAuth]%w", err)
	}
	if maa.Type == proto.MessageTypeError {
		var se proto.MessageError
		err := UnpackMessage(maa, proto.MessageTypeError, &se)
		if err != nil {
			return fmt.Errorf("[ConnectAndAuth]%w", err)
		}
		return fmt.Errorf("sync error #%d: %s", se.ErrorCode, se.Error)
	} else if maa.Type == proto.MessageTypeAuthAns {
		var sm proto.MessageAuthAnswer
		err := UnpackMessage(maa, proto.MessageTypeAuthAns, &sm)
		if err != nil {
			return fmt.Errorf("[ConnectAndAuth]%w", err)
		}
		if !sm.Success {
			return fmt.Errorf("[ConnectAndAuth]%w", err)
		}
		sc.akey = sm.AuthKey
		return nil
	}

	return fmt.Errorf("[GetObjList] unexpected answer type '%s'", maa.Type)
}
