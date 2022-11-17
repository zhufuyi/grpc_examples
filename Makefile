SHELL := /bin/bash

PROJECT_NAME := "github.com/zhufuyi/grpc_examples"
PKG := "$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/ | grep -v /api/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)


.PHONY: install
# installation of dependent tools
install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install github.com/srikrsna/protoc-gen-gotag@latest
	go install github.com/zhufuyi/sponge/cmd/protoc-gen-go-gin@latest
	go install github.com/zhufuyi/sponge/cmd/protoc-gen-go-rpc-tmpl@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/ofabry/go-callvis@latest


.PHONY: mod
# add missing and remove unused modules
mod:
	go mod tidy


.PHONY: fmt
# go format *.go files
fmt:
	gofmt -s -w .


.PHONY: ci-lint
# check the code specification against the rules in the .golangci.yml file
ci-lint: fmt
	golangci-lint run ./...


.PHONY: test
# go test *_test.go files, the parameter -count=1 means that caching is disabled
test:
	go test -count=1 -short ${PKG_LIST}


.PHONY: cover
# generate test coverage
cover:
	go test -short -coverprofile=cover.out -covermode=atomic ${PKG_LIST}
	go tool cover -html=cover.out


# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m  %-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := all
