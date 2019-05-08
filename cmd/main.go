package main

import (
	"fmt"
	"os"

	"github.com/858chain/token-shout/utils"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Token shout"
	app.Usage = "Eth / ERC20 token notification service"
	app.Version = Version
	app.Commands = []cli.Command{
		startCmd,
	}

	app.Flags = []cli.Flag{
		logLevelFlag,
		logDirFlag,
	}

	app.Before = func(c *cli.Context) error {
		//return utils.InitLogger(c.String("log-level"), c.String("log-dir"), "json")
		return utils.InitLogger(c.String("log-dir"))
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
