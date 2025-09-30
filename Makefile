BINARY_NAME=topn
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

.PHONY: build clean test install

build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/topn

install:
	go install $(LDFLAGS) ./cmd/topn

test:
	go test -v ./...

clean:
	rm -rf bin/

release-dry-run:
	goreleaser release --snapshot --clean

release:
	goreleaser release --clean