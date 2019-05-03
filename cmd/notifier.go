package main

import (
	"fmt"
	"os"

	"github.com/858chain/token-shout/api"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var startCmd = cli.Command{
	Name:    "start",
	Aliases: []string{"s"},
	Flags: []cli.Flag{
		httpAddrFlag,
		ethRpcAddrFlag,
		receiverConfPathFlag,
		walletDirFlag,
		watchIntervalFlag,
	},

	Usage: "start eth/erc20 token notification service",
	Action: func(c *cli.Context) error {

		var err error
		apiServer := api.NewApiServer(c.String("http-listen-addr"))

		log.Infof("eth rpc client with  addr: %s", c.String("eth-rpc-addr"))
		err = apiServer.InitEthClient(
			c.String("eth-rpc-addr"),       // host
			c.String("receiver-conf-path"), // receiver conf path
			c.String("wallet-dir"),
			c.GlobalString("log-dir"), // logDir
			c.Duration("watch-interval"),
		)
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
