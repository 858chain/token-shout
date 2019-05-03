package ethclient

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/858chain/token-shout/notifier"
	"github.com/858chain/token-shout/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/google/uuid"
)

type Client struct {
	config *Config

	// notifier
	noti *notifier.Notifier

	// rpc client
	rpcClient *rpc.Client

	// eth address watched, event emit if balance changed
	balanceCache map[string]float64
	lock         sync.Mutex
}

func New(config *Config) (*Client, error) {
	client := &Client{
		config: config,
		noti:   notifier.New(),

		lock:         sync.Mutex{},
		balanceCache: make(map[string]float64),
	}

	// install all default receivers
	for _, rc := range client.config.DefaultReceivers {
		receiver := notifier.NewReceiver(rc)

		uuidIns, _ := uuid.NewUUID()
		client.noti.InstallReceiver(uuidIns.String(), receiver)
	}

	// connect to geth rpc
	err := client.connect()
	if err != nil {
		return nil, err
	}

	// parse all account address from wallet directory
	addresses, err := client.loadAddressesFromWallet()
	if err != nil {
		return nil, err
	}
	for _, addr := range addresses {
		utils.L.Debugf("start watch address %s", addr)

		client.balanceCache[addr], err = client.getBalance(addr)
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}

func (c *Client) Start() error {
	errCh := make(chan error, 1)
	ctx := context.Background()
	go c.noti.Start(ctx)

	// if asked to watch eth
	if c.config.Watch.Contains("eth") {
		go c.balanceCacheSyncer(ctx, errCh)
		go c.balanceChecker(ctx, errCh)
	}

	if c.config.Watch.Contains("erc20") {
		go c.erc20TranserWatcher(ctx, errCh)
	}

	return <-errCh
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
	utils.L.Debugf("ethClient connect to %s", c.config.RpcAddr)
	c.rpcClient, err = rpc.Dial(c.config.RpcAddr)
	if err != nil {
		return err
	}

	return nil
}

// retrieve address list base on keystore file within keystore directoy
func (c *Client) loadAddressesFromWallet() ([]string, error) {
	addresses := make([]string, 0)

	files, err := ioutil.ReadDir(c.config.WalletDir)
	if err != nil {
		return addresses, err
	}

	for _, fileInfo := range files {
		address, err := c.addressFromFilename(fileInfo.Name())
		if err != nil {
			return []string{}, err
		}
		addresses = append(addresses, address)
	}

	return addresses, nil
}

// parse address from wallet filename
func (c *Client) addressFromFilename(filename string) (string, error) {
	address := fmt.Sprintf("0x%s", utils.LastSplit(filename, "-"))
	if common.IsHexAddress(address) {
		return address, nil
	}
	return "", errors.New("not a valid address")
}
