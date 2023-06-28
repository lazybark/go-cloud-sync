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
	akey           string
	cacheDir       string
	filesystemRoot string
	serverAddr     string
	serverPort     int
	login          string
	pwd            string
	c              *client.Client
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
	Error     string
	ErrorCode int
}

type MessageFullSyncReply struct {
	Success bool
	Objects []fse.FSObject
}

type MessageGetFile struct {
	Object fse.FSObject
}

type MessageFilePart struct {
	Payload []byte
}

type ExchangeMessageType int

const (
	message_type_start ExchangeMessageType = iota

	MessageTypeAuthReq
	MessageTypeAuthAns
	MessageTypeEvent
	MessageTypeFullSyncRequest
	MessageTypeFullSyncReply
	MessageTypeError
	MessageTypeGetFile

	message_type_end
)

func (t ExchangeMessageType) String() string {
	ts := [...]string{
		"MessageTypeAuthReq",
		"MessageTypeAuthAns",
		"MessageTypeEvent",
		"MessageTypeFullSyncRequest",
		"MessageTypeFullSyncReply",
		"MessageTypeError",
		"MessageTypeGetFile",
	}

	if t <= message_type_start || t >= message_type_start {
		return "unknown message type"
	}

	return ts[t-1]
}
