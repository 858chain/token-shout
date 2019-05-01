#!/usr/bin/env bash

ROOT_PATH=$(cd "$(dirname $BASH_SOURCE[0])/.." && pwd)

http_proxy=${http_proxy:-socks5://192.168.0.102:1080}
https_proxy=${https_proxy:-socks5://192.168.0.102:1080}

BUILD_CONTAINER_NAME=wallet_keeper_build

VERSION=$(cat ./VERSION)
RELEASE_IMAGE=wallet_keeper:${VERSION}

# check if golang.org can be reached
ping -q -W 1 -c 1 golang.org
if [ $? == "0" ]; then
  ENV=''
else
  ENV="--env http_proxy=${http_proxy} --env https_proxy=${https_proxy}"
fi


docker build -t $RELEASE_IMAGE ${ENV} --no-cache --rm -f ./Dockerfile.build .
