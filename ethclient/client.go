package ethclient

import (
	"github.com/858chain/token-shout/notifier"
)

type Client struct {
	config *Config

	noti *notifier.Notifier
}

func New(config *Config) (*Client, error) {
	return &Client{
		config: config,
		noti:   notifier.New(),
	}, nil
}

// connect to rpc endpoint
func (c *Client) connect() error {
	return nil
}

func (c *Client) Ping() error {
	return nil
}
