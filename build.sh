#!/usr/bin/env bash

GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -installsuffix cgo -o main main.go

docker buildx build -t mythsman/mouban:amd64-latest -f  Dockerfile --platform=linux/amd64 .

rm -rf main

GOOS=linux CGO_ENABLED=0 GOARCH=arm64 go build -ldflags="-s -w" -installsuffix cgo -o main main.go

docker buildx build -t mythsman/mouban:arm64-latest -f  Dockerfile --platform=linux/arm64 .

rm -rf main

docker manifest annotate mythsman/mouban  mythsman/mouban:amd64-latest --os linux --arch amd64

docker manifest annotate mythsman/mouban  mythsman/mouban:arm64-latest --os linux --arch arm64 --variant v8

docker push mythsman/mouban:amd64-latest

docker push mythsman/mouban:arm64-latest

docker manifest rm mythsman/mouban

sleep 1

docker manifest create mythsman/mouban mythsman/mouban:arm64-latest mythsman/mouban:amd64-latest

sleep 1

docker manifest push mythsman/mouban
