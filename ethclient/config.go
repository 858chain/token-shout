package ethclient

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/858chain/token-shout/notifier"
	"github.com/858chain/token-shout/utils"
)

type Config struct {
	// rpc addr, should be one of http://, ws://, ipc
	RpcAddr          string
	DefaultReceivers []ReceiverConfig
}

// Check config is valid.
func (c *Config) ValidCheck() error {
	if len(c.RpcAddr) == 0 {
		return errors.New("RpcAddr should not empty")
	}

	if !(strings.HasPrefix(c.RpcAddr, "http://") ||
		strings.HasPrefix(c.RpcAddr, "ws://") ||
		strings.HasSuffix(c.RpcAddr, ".ipc")) {
		return errors.New("rpcaddr should like http://, ws:// or /xxx/xx/foo.ipc")
	}

	if len(c.DefaultReceivers) == 0 {
		return errors.New("no receivers")
	}

	// Make sure every receiver is valid.
	for _, receiverConf := range c.DefaultReceivers {
		if err := receiverConf.ValidCheck(); err != nil {
			return err
		}
	}

	return nil
}

type ReceiverConfig struct {
	RetryCount int      `json:"retrycount"`
	Endpoint   string   `json:"endpoint"`
	EventTypes []string `json:"eventTypes"`
}

func (rc ReceiverConfig) ValidCheck() error {
	if rc.RetryCount <= 0 {
		return errors.New("retryCount should greater than 0")
	}

	if _, err := url.Parse(rc.Endpoint); err != nil {
		return err
	}

	if len(rc.EventTypes) == 0 {
		errors.New("eventTypes not provided")
	}

	// EventType should registered.
	for _, etype := range rc.EventTypes {
		if !utils.StringSliceContains(notifier.EventTypeRegistry, etype) {
			return errors.New(fmt.Sprintf("%s not a valid event type, make sure in %+v", etype, notifier.EventTypeRegistry))
		}
	}

	return nil
}
