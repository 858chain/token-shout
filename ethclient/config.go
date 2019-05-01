package ethclient

type Config struct {
	// rpc addr, should be one of https://, ws://, ipc
	RpcAddr          string
	DefaultReceivers []ReceiverConfig
}

type ReceiverConfig struct {
	RetryCount int      `json:"retrycount"`
	Endpoint   string   `json:"endpoint"`
	EvnetTypes []string `json:"eventTypes"`
}
