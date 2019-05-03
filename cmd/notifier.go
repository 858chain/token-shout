package main

import (
	"fmt"
	"os"

	"github.com/858chain/token-shout/api"
	"github.com/858chain/token-shout/ethclient"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var startCmd = cli.Command{
	Name:    "start",
	Aliases: []string{"s"},
	Flags: []cli.Flag{
		httpAddrFlag,
		rpcAddrFlag,
		receiversConfPathFlag,
		ethWalletDirFlag,
		ethWatchIntervalFlag,
		ERC20ContractsDirFlag,
		watchListFlag,
	},

	Usage: "start eth/erc20 token notification service",
	Action: func(c *cli.Context) error {

		var err error
		apiServer := api.NewApiServer(c.String("http-listen-addr"))

		cfg := &ethclient.Config{
			RpcAddr:   c.String("rpc-addr"), // host
			WatchList: ethclient.WatchList(c.String("watch-list")),
			LogDir:    c.GlobalString("log-dir"),

			ReceiversConfPath: c.String("receivers-conf-path"),

			EthWalletDir:      c.String("eth-wallet-dir"),
			EthWatchInterval:  c.Duration("eth-watch-interval"),
			ERC20ContractsDir: c.String("erc20-contracts-dir"),
		}

		fmt.Fprintf(os.Stdout, "%#v\n", cfg)
		// Validation Check make sure cfg valid
		err = cfg.SanityAndValidCheck()
		if err != nil {
			return err
		}

		err = apiServer.InitAndStartEthClient(cfg)
		if err != nil {
			log.Error(err)
			return err
		}

		err = apiServer.HealthCheck()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}

		fmt.Fprintf(os.Stdout, "starting notification service at addr: %s", c.String("http-listen-addr"))
		return apiServer.HttpListen()
	},
}
