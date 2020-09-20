SHELL := /bin/bash

# Go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

# Name of the executable (default is current directory name)
TARGET := $(shell echo $${PWD\#\#*/})

# Build version and time
VERSION=$(shell git describe --tags --long --dirty --always)
NOW=$(shell date +'%Y-%m-%dT%I:%M:%SZ')

# Build flags
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.timestamp=$(NOW)"

# Entrypoint files
ENTRYPOINT=cmd/json-key-remover/*.go

.PHONY: clean build test test-all lint fmt strict-check

clean:
	rm -rf ./$(TARGET)

build:
	go build $(LDFLAGS) -o $(TARGET) $(ENTRYPOINT)

test:
	go test -v -cover ./...

test-all:
	go test -race ./...

lint:
	go vet ./...

fmt:
	gofmt -l -w $(SRC)

strict-check:
	@test -z $(shell gofmt -l . | tee /dev/stderr) || echo "[WARN] Fix formatting issues with 'make fmt'"
	@for d in $$(go list ./... | grep -v /vendor/); do golint $${d}; done
