#!/usr/bin/env bash
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build server.go
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build agent.go
