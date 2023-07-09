.PHONY: build

build:
	go build -v ./cmd/gophermart

test:
	go test ./...

.DEFAULT_GOAL:= build