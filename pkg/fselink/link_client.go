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

type SyncClient struct {
	akey string
	c    *client.Client
}

func (sc *SyncClient) Init(port int, addr, login, pwd string) error {
	conf := client.Config{SuppressErrors: false, MessageTerminator: '\n'}
	c := client.New(&conf)
	sc.c = c

	err := sc.c.DialTo(addr, port, `certs/cert.pem`)
	if err != nil {
		return fmt.Errorf("[Init] %w", err)
	}

	//At first we send credentials to get auth key
	creds, err := json.Marshal(Credentials{Login: login, Password: pwd})
	if err != nil {
		return fmt.Errorf("[Init] %w", err)
	}
	mess, err := json.Marshal(ExchangeMessage{Type: MessageTypeAuthReq, Payload: creds})
	if err != nil {
		return fmt.Errorf("[Init] %w", err)
	}
	err = sc.c.SendByte(mess)
	if err != nil {
		return fmt.Errorf("[Init] %w", err)
	}
	//And await for answer with auth key
	ans := <-sc.c.MessageChan
	var ms ExchangeMessage
	err = json.Unmarshal(ans.Bytes(), &ms)
	if err != nil {
		return fmt.Errorf("[Init] %w", err)
	}
	fmt.Println(ms.Type)
	fmt.Println(string(ms.Payload))

	return nil
}

func (sc *SyncClient) SendEvent(e fse.FSEvent) error {
	ej, err := e.JSON()
	if err != nil {
		return fmt.Errorf("[SendEvent][MarshalEvent] %w", err)
	}

	m := ExchangeMessage{Type: MessageTypeEvent, AuthKey: sc.akey, Payload: ej}
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
