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
	engine         *gin.Engine

	client *ethclient.Client
}

func (api *ApiServer) InitEthClient(host, receiverConfPath, logDir string) (err error) {
	cfg := &ethclient.Config{RpcAddr: host}
	if _, err := os.Stat(receiverConfPath); err != nil && os.IsNotExist(err) {
		return err
	}

	file, err := os.OpenFile(receiverConfPath, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&cfg.DefaultReceivers)
	if err != nil {
		return err
	}

	api.client, err = ethclient.New(cfg)
	return err
}

//Check
func (api *ApiServer) HealthCheck() (err error) {

	err = api.client.Ping()
	if err != nil {
		err = errors.Wrap(err, "eth: ")
	}

	return err
}

func NewApiServer(addr string) *ApiServer {
	apiServer := &ApiServer{
		httpListenAddr: addr,
	}

	// build gin.Engine and register routers
	apiServer.buildEngine()

	return apiServer
}

func (api *ApiServer) buildEngine() {
	r := gin.Default()

	// misc API
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		err := api.HealthCheck()
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

	api.engine = r
}

func (api *ApiServer) HttpListen() error {
	return api.engine.Run(api.httpListenAddr)
}
