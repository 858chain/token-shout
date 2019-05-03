package ethclient

import (
	"os"
	"strings"
	"time"

	"github.com/858chain/token-shout/notifier"

	"github.com/pkg/errors"
)

type Watch string

func (w Watch) Contains(target string) bool {
	return strings.Contains(strings.ToLower(string(w)), strings.ToLower(target))
}

type Config struct {
	// rpc addr, should be one of http://, ws://, ipc
	RpcAddr          string
	WalletDir        string
	LogDir           string
	DefaultReceivers []notifier.ReceiverConfig
	WatchInterval    time.Duration
	Watch            Watch
}

// Check config is valid.
func (c *Config) ValidCheck() error {
	if len(c.RpcAddr) == 0 {
		return errors.New("RpcAddr should not empty")
	}

	if len(c.WalletDir) == 0 {
		return errors.New("WalletDir should not empty")
	}

	if len(c.Watch) == 0 {
		return errors.New("Watch should not empty")
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
