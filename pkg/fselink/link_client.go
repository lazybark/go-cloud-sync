package fselink

import (
	"encoding/json"
	"fmt"

	"github.com/lazybark/go-cloud-sync/pkg/fse"
	"github.com/lazybark/go-tls-server/v2/client"
)

func NewClient() FSEClientLink {
	c := &SyncClient{}
	return c
}

func (sc *SyncClient) GetObjList() (l []fse.FSObject, err error) {
	err = SendSyncMessage(sc, nil, MessageTypeFullSyncRequest)
	if err != nil {
		err = fmt.Errorf("[GetObjList] %w", err)
		return
	}
	var maa MessageFullSyncReply
	err = AwaitAnswer(sc, &maa)
	if err != nil {
		err = fmt.Errorf("[GetObjList] %w", err)
		return
	}
	l = maa.Objects
	return
}

func (sc *SyncClient) DownloadObject(objs fse.FSObject, destPath string) (err error) {
	return
}

func (sc *SyncClient) Init(port int, addr, login, pwd string) error {
	conf := client.Config{SuppressErrors: false, MessageTerminator: '\n'}
	c := client.New(&conf)
	sc.c = c

	fmt.Println("Connecting to server")
	err := sc.c.DialTo(addr, port, `certs/cert.pem`)
	if err != nil {
		return fmt.Errorf("[SyncClient][Init] %w", err)
	}

	//At first we send credentials to get auth key
	err = SendSyncMessage(sc, Credentials{Login: login, Password: pwd}, MessageTypeAuthReq)
	if err != nil {
		return fmt.Errorf("[SyncClient][Init] %w", err)
	}
	//And await for answer with auth key
	var maa MessageAuthAnswer
	err = AwaitAnswer(sc, &maa)
	if err != nil {
		return fmt.Errorf("[SyncClient][Init] %w", err)
	}
	if !maa.Success {
		return fmt.Errorf("[SyncClient][Init] %w", err)
	}
	fmt.Println(maa)
	if maa.AuthKey == "" {
		return fmt.Errorf("[SyncClient][Init] auth key is empty")
	}
	sc.akey = maa.AuthKey
	fmt.Println("Got auth key from Server", maa.AuthKey)

	return nil
}

func (sc *SyncClient) SendByte(b []byte) (int, error) {
	err := sc.c.SendByte(b)
	if err != nil {
		err = fmt.Errorf("[SendByte]: %w", err)
	}
	return len(b), err
}

func (sc *SyncClient) AwaitAnswer() (*client.Message, error) {
	ans := <-sc.c.MessageChan
	//pissble error will be added here later to prevent endless loops
	return ans, nil
}

func (sc *SyncClient) SendEvent(e fse.FSEvent) error {
	ej, err := e.JSON()
	if err != nil {
		return fmt.Errorf("[SendEvent][MarshalEvent] %w", err)
	}

	m := sc.compileMessageBody(MessageTypeEvent)
	m.Payload = ej
	mj, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("[SendEvent][MarshalMessage] %w", err)
	}

	err = sc.c.SendByte(mj)
	if err != nil {
		return err
	}

	return nil
}

func (sc *SyncClient) compileMessageBody(t ExchangeMessageType) ExchangeMessage {
	return ExchangeMessage{Type: t, AuthKey: sc.akey, Payload: []byte{}}
}
