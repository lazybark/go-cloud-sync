package client

import (
	"fmt"

	tls "github.com/lazybark/go-tls-server/v2/client"
)

func (sc *LinkClient) SetAuthKey(k string) {
	sc.akey = k
}

func (sc *LinkClient) Init(port int, addr, login, pwd string) error {

	conf := tls.Config{SuppressErrors: false, MessageTerminator: '\n'}
	c := tls.New(&conf)
	sc.c = c
	sc.serverAddr = addr
	sc.serverPort = port
	sc.login = login
	sc.pwd = pwd

	err := sc.ConnectAndAuth()
	if err != nil {
		return fmt.Errorf("[SyncClient][Init]%w", err)
	}

	return nil
}

func (sc *LinkClient) Close() error {
	return sc.c.Close()
}
