package client

import (
	"fmt"
)

// NewConnectionInSession returns a new connection with same auth key. Server will treat
// this connection as the same session
func (sc *LinkClient) NewConnectionInSession() (*LinkClient, error) {
	link, err := NewClient()
	if err != nil {
		return nil, fmt.Errorf("[DownloadObject]%w", err)
	}

	err = link.Init(sc.serverPort, sc.serverAddr, sc.login, sc.pwd)
	if err != nil {
		return nil, fmt.Errorf("[DownloadObject]%w", err)
	}

	link.SetAuthKey(sc.akey)

	return link, nil
}
