package ethclient

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/858chain/token-shout/notifier"
	"github.com/858chain/token-shout/utils"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

func (c *Client) erc20TranserWatcher(ctx context.Context, errCh chan error) {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{},
	}

	addressABIMap := make(map[common.Address][]byte)
	for _, cc := range c.config.ContractConfigs {
		query.Addresses = append(query.Addresses, common.HexToAddress(cc.Address))
		addressABIMap[common.HexToAddress(cc.Address)] = cc.Abi
	}
	fmt.Println("xxxxxxxxxxxxxxxxxxxxxxx1111111111111111")

	var ch = make(chan types.Log)
	filterCtx := context.Background()
	sub, err := c.rpcClient.EthSubscribe(filterCtx, ch, "logs", toFilterArg(query))
	if err != nil {
		fmt.Println(err)
		utils.L.Error(err)
		errCh <- err
		return
	}
	fmt.Println("xxxxxxxxxxxxxxxxxxxxxxx")

	for {
		select {
		case <-ctx.Done():
			return
		case err := <-sub.Err():
			utils.L.Error(err)
			errCh <- err
			return
		case eventLog := <-ch:
			if abiBytes, found := addressABIMap[eventLog.Address]; found {
				tokenAbi, err := abi.JSON(strings.NewReader(string(abiBytes)))
				if err != nil {
					utils.L.Error(err)
					errCh <- err
					return
				}
				var transferEvent struct {
					From  common.Address
					To    common.Address
					Value *big.Int
				}

				err = tokenAbi.Unpack(&transferEvent, "Transfer", eventLog.Data)
				if err != nil {
					utils.L.Debugf("Failed to unpack transfer event, try next event")
					continue
				}

				transferEvent.From = common.BytesToAddress(eventLog.Topics[1].Bytes())
				transferEvent.To = common.BytesToAddress(eventLog.Topics[2].Bytes())

				utils.L.Info("From", transferEvent.From.Hex())
				utils.L.Info("To", transferEvent.To.Hex())
				utils.L.Info("Value", transferEvent.Value)

				event := notifier.NewERC20LogEvent(map[string]interface{}{
					"address": eventLog.Address.Hex(),
					"from":    transferEvent.From.Hex(),
					"to":      transferEvent.To.Hex(),
					"value":   transferEvent.Value.Int64(),
				})
				c.noti.EventChan() <- event
			}
		}
	}
}

func toFilterArg(q ethereum.FilterQuery) interface{} {
	arg := map[string]interface{}{
		"fromBlock": toBlockNumArg(q.FromBlock),
		"toBlock":   toBlockNumArg(q.ToBlock),
		"address":   q.Addresses,
		"topics":    q.Topics,
	}
	if q.FromBlock == nil {
		arg["fromBlock"] = "0x0"
	}
	return arg
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	return hexutil.EncodeBig(number)
}
