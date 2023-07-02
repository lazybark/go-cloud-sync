package server

import (
	"fmt"

	gts "github.com/lazybark/go-tls-server/v3/server"
)

func (s *LinkServer) Init(extMessChan chan (*gts.Message), extConnChan chan (*gts.Connection), extErrChan chan (error)) error {
	s.extMessChan = extMessChan
	s.extConnChan = extConnChan
	s.extErrChan = extErrChan
	return nil
}

func (s *LinkServer) Listen(addr, port string) error {
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
			}
		}
	}()

	return nil
}
func (s *LinkServer) Stop() error {
	return nil

}
