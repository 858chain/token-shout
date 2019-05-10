package ethclient

import (
	"context"
	"math"
	"math/big"
	"time"

	"github.com/858chain/token-shout/notifier"
	"github.com/858chain/token-shout/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/fsnotify/fsnotify"
)

// syncer of balance, detecting balance change.
func (c *Client) balanceCacheSyncer(ctx context.Context, errCh chan error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		errCh <- err
		return
	}
	defer watcher.Close()

	watcher.Add(c.config.EthWalletDir)

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
	ticker := time.NewTicker(c.config.EthWatchInterval)

	checkFunc := func() {
		utils.L.Infof("checking balance of total %d address", len(c.balanceCache))

		balanceChangedMap := make(map[string]float64)
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

				// temproary store changed balance
				balanceChangedMap[address] = newBalance
			}
		}
		// update balanceCache
		for address, newBalance := range balanceChangedMap {
			c.lock.Lock()
			c.balanceCache[address] = newBalance
			c.lock.Unlock()
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
