package main

import (
	"time"

	"github.com/urfave/cli"
)

var httpAddrFlag = cli.StringFlag{
	Name:   "http-listen-addr",
	Value:  "0.0.0.0:8001",
	Usage:  "http address of web application",
	EnvVar: "HTTP_LISTEN_ADDR",
}

var logLevelFlag = cli.StringFlag{
	Name:   "log-level",
	Value:  "info",
	Usage:  "default log level",
	EnvVar: "LOG_LEVEL",
}

var logDirFlag = cli.StringFlag{
	Name:   "log-dir",
	EnvVar: "LOG_DIR",
	Value:  "/var/log/",
}

var rpcAddrFlag = cli.StringFlag{
	Name:   "rpc-addr",
	Value:  "http://192.168.0.101:8545",
	EnvVar: "RPCADDR",
}

var receiversConfPathFlag = cli.StringFlag{
	Name:   "receivers-conf-path",
	EnvVar: "RECEIVERS_CONF_PATH",
}

var ethWalletDirFlag = cli.StringFlag{
	Name:   "eth-wallet-dir",
	EnvVar: "ETH_WALLET_DIR",
}

var ethWatchIntervalFlag = cli.DurationFlag{
	Name:   "eth-watch-interval",
	Value:  time.Duration(20 * time.Second),
	EnvVar: "ETH_WATCH_INTERVAL",
}

var watchListFlag = cli.StringFlag{
	Name:   "watch-list",
	Value:  "usdt,eth",
	EnvVar: "WATCH_LIST",
}

var ERC20ContractsDirFlag = cli.StringFlag{
	Name:   "erc20-contracts-dir",
	Value:  "",
	EnvVar: "ERC20_CONTRACTS_DIR",
}
