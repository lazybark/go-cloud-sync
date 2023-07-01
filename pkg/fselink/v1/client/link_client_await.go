package client

import (
	"encoding/json"
	"fmt"

	"github.com/lazybark/go-cloud-sync/pkg/fselink/v1/proto"
)

func (sc *LinkClient) Await() (proto.ExchangeMessage, error) {
	var m proto.ExchangeMessage
	ans := <-sc.c.MessageChan
	err := json.Unmarshal(ans.Bytes(), &m)
	if err != nil {
		return m, fmt.Errorf("[AwaitAnswer] %w", err)
	}
	return m, nil
}
