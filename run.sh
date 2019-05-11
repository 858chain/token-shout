#!/usr/bin/env bash

./bin/token-shout --log-level debug \
  --log-dir /tmp \
  start \
  --http-listen-addr 0.0.0.0:8081 \
  --rpc-addr ws://192.168.0.101:9546 \
  --receivers-conf-path /tmp/receivers \
  --watch-list eth,usdt,dusd \
  --eth-wallet-dir /tmp/wallets \
  --eth-watch-interval 20s \
  --erc20-contracts-dir /tmp/contracts
