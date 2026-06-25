#!/usr/bin/env bash

GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -o main .

