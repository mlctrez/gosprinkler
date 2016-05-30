#!/usr/bin/env bash

# build the project and push it to the BBB in the basement  b^4?

set -e

REMOTE=bb1

# cross compiling for BBB requires GOXX in the line below
GOOS=linux GOARM=7 GOARCH=arm go build -o /tmp/sprinkler sprinkler.go

scp /tmp/sprinkler root@$REMOTE:/root/sprinkler

