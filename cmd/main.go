package main

import (
	"fmt"
	"os"

	"github.com/858chain/token-shout/utils"

	"github.com/gin-gonic/gin"
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
		// default mode is release
		gin.SetMode(gin.ReleaseMode)
		if len(os.Getenv("DEV")) > 0 {
			gin.SetMode(gin.DebugMode)
		}
		return utils.InitLogger(c.String("log-dir"), c.String("log-level"), "json")
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
