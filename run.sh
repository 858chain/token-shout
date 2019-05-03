#!/usr/bin/env bash

./bin/token-shout --log-level debug \
  --log-dir /tmp \
  start \
  --http-listen-addr 0.0.0.0:8081 \
  --eth-rpc-addr http://154.8.201.160:8545 \
  --receiver-conf-path /tmp/receivers \
  --wallet-dir /tmp/wallets \
  --watch-interval 20s \
  --watch eth,erc20
