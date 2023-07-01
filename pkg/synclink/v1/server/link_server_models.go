package server

import (
	gts "github.com/lazybark/go-tls-server/v3/server"
)

func NewServer() *LinkServer {
	s := &LinkServer{}
	return s
}

type LinkServer struct {
	extMessChan chan (*gts.Message)
	extConnChan chan (*gts.Connection)
	extErrChan  chan (error)
}
