package fselink

import (
	"encoding/json"
	"fmt"

	"github.com/lazybark/go-cloud-sync/pkg/fse"
)

func (sc *SyncClient) SendEvent(e fse.FSEvent) error {
	ej, err := e.JSON()
	if err != nil {
		return fmt.Errorf("[SendEvent][MarshalEvent] %w", err)
	}

	m := sc.compileMessageBody(MessageTypeEvent)
	m.Payload = ej
	mj, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("[SendEvent][MarshalMessage] %w", err)
	}

	err = sc.c.SendByte(mj)
	if err != nil {
		return err
	}

	return nil
}

func (sc *SyncClient) SendByte(b []byte) (int, error) {
	err := sc.c.SendByte(b)
	if err != nil {
		err = fmt.Errorf("[SendByte]: %w", err)
	}
	return len(b), err
}
