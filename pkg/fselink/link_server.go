package fselink

import (
	"fmt"

	"github.com/lazybark/go-cloud-sync/pkg/fse"
	gts "github.com/lazybark/go-tls-server/v2/server"
)

func NewServer() FSEServerPool {
	s := &SyncServer{}
	return s
}

type SyncServer struct {
	extMessChan chan (*gts.Message)
	extConnChan chan (*gts.Connection)
	extErrChan  chan (error)
}

func (s *SyncServer) Init(extMessChan chan (*gts.Message), extConnChan chan (*gts.Connection), extErrChan chan (error)) error {
	s.extMessChan = extMessChan
	s.extConnChan = extConnChan
	s.extErrChan = extErrChan
	return nil
}

func (s *SyncServer) Listen(addr, port string) error {
	done := make(chan bool)

	conf := &gts.Config{KeepOldConnections: 1, NotifyAboutNewConnections: true}
	srv, err := gts.New("localhost", `certs/cert.pem`, `certs/key.pem`, conf)
	if err != nil {
		return fmt.Errorf("[SERVER][LISTEN] error notifying server: %w", err)
	}
	go srv.Listen(port)
	go func() {
		defer close(done)
		for {
			select {
			case err := <-srv.ErrChan:
				s.extErrChan <- err
			case conn := <-srv.ConnChan:
				s.extConnChan <- conn
			case m := <-srv.MessageChan:
				s.extMessChan <- m

			}
		}
	}()

	return nil
}
func (s *SyncServer) Stop() error {
	return nil

}
func (s *SyncServer) NotifyClients(e fse.FSEvent) error {
	return nil
}
