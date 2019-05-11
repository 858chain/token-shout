**Token shout** notify balance changes for both **ETH/ERC20** tokens, save registered upstream servics 
from repeatly querying for any event/balance change from ethreum chains.

## Run

```bash
$ ./bin/token-shout --log-level debug \
  --log-dir /tmp \
  start \
  --http-listen-addr 0.0.0.0:8081 \
  --rpc-addr ws://192.168.0.101:9546 \
  --receivers-conf-path /tmp/receivers \
  --watch-list eth,usdt,dusd \
  --eth-wallet-dir /tmp/wallets \
  --eth-watch-interval 20s \
  --erc20-contracts-dir /tmp/contracts
```



## Run In Docker(docker-compose style)

```
version: "3"

services:
  token-shout:
    image: 858chain/token-shout:latest
    volumes:
      - /data/token-shout-data:/data
      - /data/wallet-keeper-data/eth-wallet:/eth-wallet
    network_mode: host
    environment:
      - LOG_LEVEL=debug
      - HTTP_LISTEN_ADDR=127.0.0.1:8001
      - LOG_DIR=/data/log
      - RPCADDR=ws://localhost:9546
      - RECEIVERS_CONF_PATH=/data/receivers/receivers.json
      - WATCH_LIST=eth
      - ETH_WALLET_DIR=/eth-wallet
      - ETH_WATCH_INTERVAL=20s
      - ERC20_CONTRACTS_DIR=/data/contracts
    restart: always
```

**Notice**:

1、RECEIVERS_CONF_PATH， it should be the receivers.json full filepath , not path.

2、ETH_WALLET_DIR， it should be geth wallet keystone folder.

3、ERC20_CONTRACTS_DIR, it shoud be contract setting folder. the config js can generate by run **scripts/contract_abi.rb** file.


## How to config

before start token-shout, you shoud setting a **receivers.json** config. below is sample config:

```json
[
  {"endpoint": "http://exmaple.com/callback",
    "eventTypes": [
      "eth_balance_change_event",
      "erc20_log_event"
    ],
    "retryCount": 1,
    "newBalanceRemaining": 0.01,
    "from": ["*"],
    "to": ["all-in-wallet"]
  },
  {"endpoint": "http://exmaple2.com/callback",
    "eventTypes": [
      "eth_balance_change_event"
    ],
    "retryCount": 2,
    "newBalanceRemaining": 0.01,
    "from": ["*"],
    "to": ["all-in-wallet", "0x9d8ec6a2cfe2a6abdf5a6041e752ba851abb4dbb"]
  }
]
```

**Notice:**

1、newBalanceRemaining： filter event by > 0.01 can trigger event  notification.

2、retryCount: send notification msg retry count, when the consumer  return http status 200 code, the token shout service will confirm the msg is successful sending.  

