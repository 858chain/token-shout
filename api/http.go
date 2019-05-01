package api

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/858chain/token-shout/ethclient"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// http api list
var METHODS_SUPPORTED = map[string]string{
	// misc
	"/ping":   "check if api service valid and backend bitcoin service healthy",
	"/health": "check system status",
	"/help":   "display this message",

	// useful APIs here
	"/install":   "install receiver",
	"/uninstall": "uninstall/remove receiver",
	"/list":      "list all avaliable receivers",
}

type ApiServer struct {
	httpListenAddr string
	// http gin engine instance
	engine *gin.Engine

	// eth rpc client
	client *ethclient.Client
}

// InitEthClient do the config  validation for make initial call to eth backend.
// Error return if malformat config or rpc server unreachable.
func (api *ApiServer) InitEthClient(host, receiverConfPath, walletDir, logDir string) (err error) {
	cfg := &ethclient.Config{
		RpcAddr:   host,
		WalletDir: walletDir,
		LogDir:    logDir,
	}

	// fail fast if receiver config not exists.
	if _, err := os.Stat(receiverConfPath); err != nil && os.IsNotExist(err) {
		return err
	}

	file, err := os.OpenFile(receiverConfPath, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	// return error if malformat receiver config file.
	err = json.NewDecoder(file).Decode(&cfg.DefaultReceivers)
	if err != nil {
		return err
	}

	// Validation Check make sure cfg valid
	err = cfg.ValidCheck()
	if err != nil {
		return err
	}

	api.client, err = ethclient.New(cfg)
	return err
}

// Check eth rpc server connectivity.
func (api *ApiServer) HealthCheck() (err error) {
	err = api.client.Ping()
	if err != nil {
		err = errors.Wrap(err, "eth: ")
	}

	return err
}

// Hook all HTTP routes and start listen on `addr`
func NewApiServer(addr string) *ApiServer {
	apiServer := &ApiServer{
		httpListenAddr: addr,
	}

	r := gin.Default()

	// misc API
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		err := apiServer.HealthCheck()
		if err != nil {
			c.JSON(500, gin.H{
				"message": fmt.Sprint(err),
			})
		} else {
			c.JSON(200, gin.H{
				"message": "healthy",
			})
		}
	})

	r.GET("/help", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"methods": METHODS_SUPPORTED,
		})
	})

	apiServer.engine = r
	return apiServer
}

func (api *ApiServer) HttpListen() error {
	return api.engine.Run(api.httpListenAddr)
}
