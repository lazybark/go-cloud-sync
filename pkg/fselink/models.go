package fselink

import (
	"github.com/lazybark/go-cloud-sync/pkg/fse"
	"github.com/lazybark/go-tls-server/v2/client"
)

type SyncMessenger interface {
	SendByte([]byte) (int, error)
}

type SyncReciever interface {
	AwaitAnswer() (*client.Message, error)
}

type SyncClient struct {
	akey string
	c    *client.Client
}

type Credentials struct {
	Login    string
	Password string
}

type ExchangeMessage struct {
	Type    ExchangeMessageType
	AuthKey string
	Payload []byte
}

type MessageAuthAnswer struct {
	Success bool
	AuthKey string
}

type MessageError struct {
	Success   bool
	Error     string
	ErrorCode int
}

type MessageFullSyncReply struct {
	Success bool
	Objects []fse.FSObject
}

type ExchangeMessageType int

const (
	MessageTypeAuthReq ExchangeMessageType = iota + 1
	MessageTypeAuthAns
	MessageTypeEvent
	MessageTypeFullSyncRequest
	MessageTypeFullSyncReply
	MessageTypeError
)
