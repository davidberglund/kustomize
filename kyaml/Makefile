# Copyright 2019 The Kubernetes Authors.
# SPDX-License-Identifier: Apache-2.0

.PHONY: generate fix vet fmt test build tidy

GOPATH := $(shell go env GOPATH)

build:
	go build -v -o $GOPATH/bin/kyaml .

all: generate fix vet fmt test lint lint tidy

fix:
	go fix ./...

fmt:
	go fmt ./...

generate:
	go generate ./...

tidy:
	go mod tidy

lint:
	(which $(GOPATH)/bin/golangci-lint || go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.19.1)
	$(GOPATH)/bin/golangci-lint run ./...

test:
	go test -cover ./...

vet:
	go vet ./...

