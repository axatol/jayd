GO_MODULE_NAME = github.com/axatol/jayd
BUILD_COMMIT ?= $(shell git rev-parse HEAD)
GO_BUILD_LDFLAGS = -X '$(GO_MODULE_NAME)/pkg/config.BuildCommit=$(BUILD_COMMIT)'
GO_BUILD_LDFLAGS += -X '$(GO_MODULE_NAME)/pkg/config.BuildTime=$(shell date +"%Y-%m-%dT%H:%M:%S%z")'

vet:
	go vet ./...

deps:
	go mod download

build:
	go build -o ./bin/server -ldflags="$(GO_BUILD_LDFLAGS)" ./cmd/server/main.go
