#!/usr/bin/env bash

GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -o main main.go

sudo docker build -t mythsman/mouban -f Dockerfile --platform=linux/amd64 .

rm -rf main

sudo docker push mythsman/mouban
