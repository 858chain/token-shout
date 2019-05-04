#!/usr/bin/env bash

./bin/token-shout --log-level debug \
  --log-dir /tmp \
  start \
  --http-listen-addr 0.0.0.0:8081 \
  --rpc-addr ws://107.150.126.20:9546 \
  --receivers-conf-path /tmp/receivers \
  --watch-list eth,usdt,dusd \
  --eth-wallet-dir /tmp/wallets \
  --eth-watch-interval 20s \
  --erc20-contracts-dir /tmp/contracts

#firewall-cmd --permanent --zone=public --add-rich-rule=' rule family="ipv4" source address="27.189.223.168/32" port protocol="tcp" port="9546" accept'
