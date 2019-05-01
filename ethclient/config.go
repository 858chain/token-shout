package ethclient

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/858chain/token-shout/notifier"
	"github.com/858chain/token-shout/utils"

	"github.com/pkg/errors"
)

type Config struct {
	// rpc addr, should be one of http://, ws://, ipc
	RpcAddr          string
	WalletDir        string
	LogDir           string
	DefaultReceivers []ReceiverConfig
}

// Check config is valid.
func (c *Config) ValidCheck() error {
	if len(c.RpcAddr) == 0 {
		return errors.New("RpcAddr should not empty")
	}

	if len(c.WalletDir) == 0 {
		return errors.New("WalletDir should not empty")
	}

	stat, err := os.Stat(c.WalletDir)
	if err != nil {
		return errors.Wrap(err, "WalletDir: ")
	}

	if !stat.IsDir() {
		return errors.New("walletDir is not a directory")
	}

	stat, err = os.Stat(c.LogDir)
	if err != nil {
		return errors.Wrap(err, "logDir: ")
	}

	if !stat.IsDir() {
		return errors.New("logDir is not a directory")
	}

	// rpcaddr format check
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
