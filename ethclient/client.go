package ethclient

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"sync"
	"time"

	"github.com/858chain/token-shout/notifier"
	"github.com/858chain/token-shout/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/fsnotify/fsnotify"
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
	go c.balanceCacheSyncer(ctx, errCh)
	go c.balanceChecker(ctx, errCh)

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

// syncer of balance, detecting balance change.
func (c *Client) balanceCacheSyncer(ctx context.Context, errCh chan error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		errCh <- err
		return
	}
	defer watcher.Close()

	watcher.Add(c.config.WalletDir)

	for {
		select {
		case <-ctx.Done():
			return
		// watch for events
		case event := <-watcher.Events:
			if event.Op&fsnotify.Create == fsnotify.Create {
				utils.L.Infof("created file: %s", event.Name)
				address, err := c.addressFromFilename(event.Name)
				if err == nil {
					balance, err := c.getBalance(address)
					if err == nil {
						c.lock.Lock()
						c.balanceCache[address] = balance
						c.lock.Unlock()
					}
				}
			}

		case err := <-watcher.Errors:
			utils.L.Error(err)
			errCh <- err
		}
	}
}

// getBalance call eth rpc and return balance in float64 format
func (c *Client) getBalance(address string) (float64, error) {
	utils.L.Debugf("get balance for address %s", address)

	var balance hexutil.Big
	err := c.rpcClient.CallContext(context.Background(), &balance,
		"eth_getBalance", common.HexToAddress(address), "latest")
	if err != nil {
		return 0, err
	}

	float64Value, _ := weiToEther(balance.ToInt()).Float64()
	return float64Value, nil
}

// balanceChecker periodically check balance of all known addresses.
func (c *Client) balanceChecker(ctx context.Context, errCh chan error) {
	ticker := time.NewTicker(c.config.WatchInterval)

	checkFunc := func() {
		utils.L.Infof("checking balance of total %d address", len(c.balanceCache))
		for address, balance := range c.balanceCache {
			newBalance, err := c.getBalance(address)
			if err != nil {
				utils.L.Error(err)
			}

			// balance updated
			// 1, need to emit event to notify any receiver may care about.
			// 2, update cached addressBalance
			if balance != newBalance {
				event := notifier.NewEthBalanceChangeEvent(map[string]interface{}{
					"address":    address,
					"newBalance": newBalance,
					"balance":    balance,
					"to":         address,
				})
				c.noti.EventChan() <- event

				c.lock.Lock()
				c.balanceCache[address] = newBalance
				c.lock.Unlock()
			}
		}
	}

	// initial check
	checkFunc()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			checkFunc()
		}
	}
}

// turn wei in big.Int into ether in big.Float
func weiToEther(wei *big.Int) *big.Float {
	weiFloat := new(big.Float)
	weiFloat.SetString(wei.String())
	return new(big.Float).Quo(weiFloat, big.NewFloat(math.Pow10(18)))
}
