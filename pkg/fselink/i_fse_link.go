package fselink

import (
	"github.com/lazybark/go-cloud-sync/pkg/fse"
	gts "github.com/lazybark/go-tls-server/v2/server"
)

// FSEClientLink is used by clients to connect to a server
type FSEClientLink interface {
	Init(port int, addr, login, pwd string) error
	//Break() error
	SendEvent(e fse.FSEvent) error
	/*
		//RequestObjectList returns list of objects on the server that changed since
		RequestObjectList(t time.Time) ([]fse.FSObject, error)

		//RequestObject requests current state of specific object
		RequestObject(o fse.FSObject) error*/
}

// FSEServerPool is used by server to keep
type FSEServerPool interface {
	Init(chan (*gts.Message), chan (*gts.Connection), chan (error)) error
	Listen(string, string) error
	//Stop() error
	//NotifyClients(e fse.FSEvent) error
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
	Success   bool
	AuthKey   string
	ErrorCode ErrorCode
	Error     string
}

type ExchangeMessageType int

const (
	MessageTypeAuthReq ExchangeMessageType = iota + 1
	MessageTypeAuthAns
	MessageTypeEvent
)

type ErrorCode int

const (
	err_codes_start ErrorCode = iota
	ErrorCodeWrongCreds
	err_codes_end
)

func (ec ErrorCode) String() string {
	codes := [...]string{"wrong login or password"}
	if ec <= err_codes_start || ec >= err_codes_end {
		return "unknown error"
	}
	return codes[ec-1]
}
