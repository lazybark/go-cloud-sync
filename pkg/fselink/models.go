package fselink

import (
	"github.com/lazybark/go-tls-server/v2/client"
)

type SyncMessenger interface {
	SendByte([]byte) (int, error)
}

type SyncReciever interface {
	AwaitAnswer() (*client.Message, error)
}

type SyncClient struct {
	akey       string
	serverAddr string
	serverPort int
	login      string
	pwd        string
	c          *client.Client
}
