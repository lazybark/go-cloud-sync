package server

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/proto"
)

func (s *FSWServer) processGetFile(c *syncConnection, m proto.ExchangeMessage) {
	a, err := m.ReadObjectData()
	if err != nil {
		s.extErc <- err
		return
	}

	//Do not SEND dirs
	//Client should create dir after full sync request
	if a.Object.IsDir {
		c.SendError(proto.ErrWrongObjectType)
		return
	}

	a.Object.Path = strings.ReplaceAll(a.Object.Path, "?ROOT_DIR?", "?ROOT_DIR?,"+c.uid)

	fileName := s.fp.GetPathUnescaped(a.Object)
	dbObj, err := s.stor.GetObject(a.Object.Path, a.Object.Name)
	if err != nil {
		s.extErc <- err
		return
	}
	//Do not SEND dirs (2) - if client's a smartass and still wants it somehow
	if dbObj.IsDir {
		c.SendError(proto.ErrWrongObjectType)
		return
	}

	fileData, err := os.Open(fileName)
	if err != nil {
		c.SendError(proto.ErrFileReadingFailed)
		return
	}

	// TLS record size can be up to 16KB but some extra bytes may apply
	// https://hpbn.co/transport-layer-security-tls/#optimize-tls-record-size
	buf := make([]byte, 15360)
	n := 0

	r := bufio.NewReader(fileData)

	for {
		n, err = r.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			c.SendError(proto.ErrFileReadingFailed)
			break
		}

		err = c.SendMessage(proto.MessageFilePart{Payload: buf[:n]}, proto.MessageTypeFileParts)
		if err != nil {
			s.extErc <- err
			break
		}
	}
	fileData.Close()

	err = c.SendMessage(nil, proto.MessageTypeFileEnd)
	if err != nil {
		s.extErc <- err
		return
	}

	fmt.Println("SENT FILE")
	c.Close()
}
