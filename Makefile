.PHONY: help
.DEFAULT_GOAL := help

CURRENT_REVISION = $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS = "-X main.revision=$(CURRENT_REVISION)"

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build All
	make build-proto
	go build -o satelit -ldflags $(BUILD_LDFLAGS) .

build-linux: ## Build for Linux
	make build-proto
	GOOS=linux GOARCH=amd64 go build -o satelit-linux-amd64 -ldflags $(BUILD_LDFLAGS) .

build-proto: ## Build proto file
	mkdir -p ./api/satelit
	protoc -I ./api/satelit --go_out=plugins=grpc:./api/satelit ./api/satelit/*.proto
	protoc -I ./api/satelit_datastore --experimental_allow_proto3_optional --go_out=plugins=grpc:./api/satelit_datastore ./api/satelit_datastore/*.proto

test: ## Exec test
	go test -v ./...

up-dev: ## Run application for development
	GOOS=linux GOARCH=amd64 make build
	docker-compose up
