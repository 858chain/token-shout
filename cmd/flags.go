package main

import (
	"github.com/urfave/cli"
)

var httpAddrFlag = cli.StringFlag{
	Name:   "http-listen-addr",
	Value:  "0.0.0.0:8000",
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

var ethRpcAddrFlag = cli.StringFlag{
	Name:   "eth-rpc-addr",
	Value:  "http://192.168.0.101:8545",
	EnvVar: "ETH_RPCADDR",
}

var receiverConfPathFlag = cli.StringFlag{
	Name:   "receiver-conf-path",
	EnvVar: "RECEIVER_CONF_PATH",
}

var walletDirFlag = cli.StringFlag{
	Name:   "wallet-dir",
	EnvVar: "WALLET_DIR",
}
