package server

import (
	"fmt"

	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/proto"
)

func (s *FSWServer) watcherRoutine() {
	fmt.Println("Waiting for connections")
	for {
		select {
		case err, ok := <-s.srvErrChan:
			if !ok {
				return
			}
			s.extErc <- err
		case connection, ok := <-s.srvConnChan:
			if !ok {
				return
			}
			//Add connection to pool to be able to control it after
			c := SyncConnection{tlsConnection: connection}
			s.addToPool(&c)

			go func() {
				for !c.IsClosed() {
					m, err := c.Await()
					if err != nil {
						s.extErc <- err
					}
					//SECURITY CHECKS
					if m.Type != proto.MessageTypeAuthReq {
						if m.AuthKey == "" {
							c.SendError(proto.ErrForbidden)
							continue
						}
						ok, err := s.checkToken(c.clientTokenHash, m.AuthKey)
						if err != nil {
							c.SendError(proto.ErrInternalServerError)
							s.extErc <- err
							continue
						}
						if !ok {
							c.SendError(proto.ErrForbidden)
							continue
						}
					}

					//MESSAGE PROCESSING
					if m.Type == proto.MessageTypeAuthReq {

						s.processAuth("", "", &c)

					} else if m.Type == proto.MessageTypeFullSyncRequest {

						s.processFullSyncRequest(&c)

					} else if m.Type == proto.MessageTypeGetFile {

						s.processGetFile(&c, m)

					} else if m.Type == proto.MessageTypePushFile {

						s.processPushFile(&c, m)

					} else if m.Type == proto.MessageTypeDeleteObject {

						s.processDelete(&c, m)

					} else {
						c.SendError(proto.ErrUnexpectedMessageType)
					}
				}
				s.remFromPool(&c)

				/*for mess := range connection.MessageChan {
					//GET MESSAGE HEAD
					err := json.Unmarshal(mess.Bytes(), &m)
					if err != nil {
						s.extErc <- err
						continue
					}

					//SECURITY CHECKS
					if m.Type != proto.MessageTypeAuthReq {
						if m.AuthKey == "" {
							c.SendError(proto.ErrForbidden)
							continue
						}
						ok, err := s.checkToken(c.clientTokenHash, m.AuthKey)
						if err != nil {
							c.SendError(proto.ErrInternalServerError)
							s.extErc <- err
							continue
						}
						if !ok {
							c.SendError(proto.ErrForbidden)
							continue
						}
					}

					//MESSAGE PROCESSING
					if m.Type == proto.MessageTypeAuthReq {

						s.processAuth("", "", &c)

					} else if m.Type == proto.MessageTypeFullSyncRequest {

						s.processFullSyncRequest(&c)

					} else if m.Type == proto.MessageTypeGetFile {

						s.processGetFile(&c, m)

					} else if m.Type == proto.MessageTypePushFile {

						s.processPushFile(&c, m)

					} else if m.Type == proto.MessageTypeDeleteObject {

						s.processDelete(&c, m)

					} else {
						c.SendError(proto.ErrUnexpectedMessageType)
						continue
					}
				}*/
			}()
		}
	}

}
