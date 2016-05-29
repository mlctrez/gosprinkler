#!/usr/bin/env bash

# for testing, build the project and push it to the BBB in the basement

set -e

REMOTE=bb1


GOOS=linux GOARM=7 GOARCH=arm go build -o /tmp/sprinkler sprinkler.go

scp /tmp/sprinkler root@$REMOTE:/root/sprinkler

ssh root@$REMOTE /root/sprinkler