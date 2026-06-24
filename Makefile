.PHONY: run test test-compile build fmt

run:
	go run ./cmd/mouban

test:
	CGO_ENABLED=0 go test ./...

test-compile:
	CGO_ENABLED=0 go test ./... -run '^$'

build:
	CGO_ENABLED=0 go build -ldflags='-s -w' -o main ./cmd/mouban

fmt:
	gofmt -w $(shell find . -name '*.go' -type f)
