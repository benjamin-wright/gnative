#!/bin/bash

alias go='docker run --rm -it -v go-cache:/go/pkg/mod -v $(pwd):/var/apps/src --network host -w /var/apps/src golang:latest go'