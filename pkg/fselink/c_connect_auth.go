package fselink

import "fmt"

func (sc *SyncClient) ConnectAndAuth() error {
	err := sc.c.DialTo(sc.serverAddr, sc.serverPort, `certs/cert.pem`)
	if err != nil {
		return fmt.Errorf("[ConnectAndAuth]%w", err)
	}

	//At first we send credentials to get auth key
	err = SendSyncMessage(sc, Credentials{Login: sc.login, Password: sc.pwd}, MessageTypeAuthReq)
	if err != nil {
		return fmt.Errorf("[ConnectAndAuth]%w", err)
	}

	var maa ExchangeMessage
	err = AwaitAnswer(sc, &maa)
	if err != nil {
		return fmt.Errorf("[ConnectAndAuth]%w", err)
	}
	if maa.Type == MessageTypeError {
		var se MessageError
		err := UnpackMessage(maa, MessageTypeError, &se)
		if err != nil {
			return fmt.Errorf("[ConnectAndAuth]%w", err)
		}
		return fmt.Errorf("sync error #%d: %s", se.ErrorCode, se.Error)
	} else if maa.Type == MessageTypeAuthAns {
		var sm MessageAuthAnswer
		err := UnpackMessage(maa, MessageTypeAuthAns, &sm)
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
