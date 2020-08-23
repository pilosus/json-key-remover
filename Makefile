.PHONY: build clean

VERSION=$(shell git describe --tags --long --dirty)
NOW=$(shell date +'%Y-%m-%dT%I:%M:%SZ')

clean:
	rm -f cmd/json-key-remover/json-key-remover
	rm -rf ./json-key-remover

build:
	go build -ldflags "-X main.version=$(VERSION) -X main.timestamp=$(NOW)" -o json-key-remover cmd/json-key-remover/*.go
