package client

import (
	"github.com/lazybark/go-tls-server/v2/client"
)

type LinkClient struct {
	akey       string
	serverAddr string
	serverPort int
	login      string
	pwd        string
	c          *client.Client
}

func NewClient() (*LinkClient, error) {
	c := &LinkClient{}

	return c, nil
}
