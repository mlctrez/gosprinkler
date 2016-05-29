#!/usr/bin/env bash

# for testing, build the project and push it to the BBB in the basement

set -e

GOOS=linux GOARM=7 GOARCH=arm go build -o /tmp/sprinkler sprinkler.go

scp /tmp/sprinkler root@bb2:/root/sprinkler

ssh root@bb2 /root/sprinkler