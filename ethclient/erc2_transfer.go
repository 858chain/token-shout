package ethclient

import (
	"context"
	//"github.com/ethereum/go-ethereum/common"
	//"github.com/ethereum/go-ethereum/types"
)

func (c *Client) erc20TranserWatcher(ctx context.Context, errCh chan error) {
	//contractAddress := common.HexToAddress("0x0d8775f648430679a709e98d2b0cb6250d2887ef")

	//query := ethereum.FilterQuery{
	//Addresses: []common.Address{contractAddress},
	//}

	//var ch = make(chan types.Log)
	//ctx := context.Background()

	//sub, err := c.rpcClient.SubscribeFilterLogs(ctx, query, ch)

	//if err != nil {
	//log.Fatal(err)
	//}

	//tokenAbi, err := abi.JSON(strings.NewReader(string(token.TokenABI)))

	//if err != nil {
	//log.Fatal(err)
	//}

	//for {
	//select {
	//case err := <-sub.Err():
	//log.Fatal(err)
	//case eventLog := <-ch:
	//var transferEvent struct {
	//From  common.Address
	//To    common.Address
	//Value *big.Int
	//}

	//err = tokenAbi.Unpack(&transferEvent, "Transfer", eventLog.Data)

	//if err != nil {
	//log.Println("Failed to unpack")
	//continue
	//}

	//transferEvent.From = common.BytesToAddress(eventLog.Topics[1].Bytes())
	//transferEvent.To = common.BytesToAddress(eventLog.Topics[2].Bytes())

	//log.Println("From", transferEvent.From.Hex())
	//log.Println("To", transferEvent.To.Hex())
	//log.Println("Value", transferEvent.Value)
	//}
	//}
}
