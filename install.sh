#!/bin/bash

set -e

docker run -w /go/src/app --rm -v go-cache:/go/pkg -v $(pwd):/go/src/app -e GOOS=darwin -e GOARCH=amd64 golang:latest go build -o gnative
mv gnative /usr/local/bin/gnative

set +e