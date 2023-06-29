package v1

import (
	"fmt"

	"github.com/lazybark/go-cloud-sync/pkg/fp"
	"github.com/lazybark/go-cloud-sync/pkg/fse"
)

func (s *FSWServer) Init(root, cacheRoot, host, port, escSymbol string, evc chan (fse.FSEvent), erc chan (error)) error {
	s.conf.root = root
	s.conf.cacheRoot = cacheRoot
	s.conf.host = host
	s.conf.port = port
	s.conf.escSymbol = escSymbol
	s.extEvc = evc
	s.extErc = erc

	s.fp = fp.NewFPv1(escSymbol, root, cacheRoot)

	err := s.w.Init(root, s.extEvc, s.extErc)
	if err != nil {
		return fmt.Errorf("[INIT]: %w", err)
	}

	err = s.htsrv.Init(s.srvMessChan, s.srvConnChan, s.srvErrChan)
	if err != nil {
		return fmt.Errorf("[INIT]: %w", err)
	}

	fmt.Printf("Server started on %s:%s\n", host, port)
	fmt.Printf("Root path:%s\n", root)
	fmt.Printf("Cache path:%s\n", cacheRoot)

	return nil
}

func (s *FSWServer) Start() error {
	s.rescanOnce()
	fmt.Println("Root directory processed")

	go s.htsrv.Listen(s.conf.host, s.conf.port)
	go s.watcherRoutine()
	s.isActive = true
	err := s.w.Start()
	if err != nil {
		return fmt.Errorf("[SERVER][START]%w", err)
	}
	return nil
}

func (s *FSWServer) Stop() error {
	err := s.w.Stop()
	if err != nil {
		return fmt.Errorf("[STOP] can not stop watcher: %w", err)
	}
	close(s.extEvc)
	close(s.extErc)
	close(s.evc)
	close(s.erc)

	s.isActive = false
	fmt.Println("Server stopped")

	return nil
}
