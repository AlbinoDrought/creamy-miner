#!/bin/bash

set -x
set -e

OUTPUT=creamy-miner_$(date +"%Y-%m-%d").zip

rm -f creamy-miner
CGO_ENABLED=0 $(go env GOPATH)/bin/garble -literals build -ldflags "-s -w" -o creamy-miner

rm -f $OUTPUT
7z a $OUTPUT creamy-miner README.md LICENSE