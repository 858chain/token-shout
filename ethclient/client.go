package ethclient

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/858chain/token-shout/notifier"
	"github.com/858chain/token-shout/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	config *Config
	l      *log.Logger

	// notifier
	noti *notifier.Notifier

	rpcClient *rpc.Client
}

func New(config *Config) (*Client, error) {
	client := &Client{
		config: config,
		noti:   notifier.New(),
	}

	// log initialization
	logPath := filepath.Join(config.LogDir, "notifier.log")
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return nil, err
	}

	client.l = &log.Logger{
		Out:       logFile,
		Formatter: new(log.JSONFormatter),
		Level:     log.DebugLevel,
	}

	// connect to geth rpc
	err = client.connect()
	if err != nil {
		return nil, err
	}

	// parse all account address from wallet directory
	addresses, err := client.addressesFromWalletDir()
	if err != nil {
		return nil, err
	}
	for _, addr := range addresses {
		client.l.Debugf("start watch address %s", addr)
	}

	return client, nil
}

func (c *Client) Ping() error {
	var hexHeight string
	err := c.rpcClient.CallContext(context.Background(), &hexHeight, "eth_blockNumber")
	if err != nil {
		return err
	}

	return nil
}

// connect to rpc endpoint
func (c *Client) connect() (err error) {
	c.l.Debugf("ethClient connect to %s", c.config.RpcAddr)
	c.rpcClient, err = rpc.Dial(c.config.RpcAddr)
	if err != nil {
		return err
	}

	return nil
}

// retrieve address list base on keystore file within keystore directoy
func (c *Client) addressesFromWalletDir() ([]string, error) {
	addresses := make([]string, 0)

	files, err := ioutil.ReadDir(c.config.WalletDir)
	if err != nil {
		return addresses, err
	}

	for _, fileInfo := range files {
		address := fmt.Sprintf("0x%s", utils.LastSplit(fileInfo.Name(), "-"))
		if common.IsHexAddress(address) {
			addresses = append(addresses, address)
		}
	}

	return addresses, nil
}
