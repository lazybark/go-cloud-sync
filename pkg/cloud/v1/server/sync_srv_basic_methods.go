package server

import (
	"fmt"
	"strings"

	"github.com/lazybark/go-cloud-sync/pkg/cloud/v1/fp"
)

func (s *FSWServer) Init(root, cacheRoot, host, port, escSymbol string, erc chan (error)) error {
	s.conf.root = root
	s.conf.cacheRoot = cacheRoot
	s.conf.host = host
	s.conf.port = port
	s.conf.escSymbol = escSymbol
	s.extErc = erc

	s.fp = fp.NewFPv1(escSymbol, root, cacheRoot)

	err := s.htsrv.Init(s.srvMessChan, s.srvConnChan, s.srvErrChan)
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

	return nil
}

func (s *FSWServer) Stop() error {
	close(s.extErc)
	close(s.erc)

	s.isActive = false
	fmt.Println("Server stopped")

	return nil
}

func (s *FSWServer) ExtractOwnerFromPath(p string, o string) string {
	dirs := strings.Split(p, s.conf.escSymbol)
	if len(dirs) < 2 {
		return p
	} else {
		return strings.ReplaceAll(p, "?ROOT_DIR?,"+o, "?ROOT_DIR?")
	}
}

func (s *FSWServer) GetOwnerFromPath(p string) string {
	dirs := strings.Split(p, s.conf.escSymbol)
	if len(dirs) < 2 {
		return ""
	} else {
		return dirs[1]
	}
}
