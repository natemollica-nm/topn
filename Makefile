BINARY_NAME=topn
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

.PHONY: build clean test install release-dry-run release tag-release

build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/topn

install:
	go install $(LDFLAGS) ./cmd/topn

test:
	go test -v ./...

clean:
	rm -rf bin/ dist/

release-dry-run:
	goreleaser release --snapshot --clean

release: check-git-clean test
	@if [ -z "$(shell git tag -l)" ]; then \
		echo "No tags found. Creating initial tag v0.1.0..."; \
		git tag v0.1.0; \
		git push origin v0.1.0; \
	fi
	goreleaser release --clean

tag-release: check-git-clean
	@echo "Current version: $(VERSION)"
	@read -p "Enter new version (e.g., v1.0.0): " version; \
	if [ -n "$$version" ]; then \
		git tag $$version && \
		git push origin $$version && \
		echo "Tagged and pushed $$version"; \
	else \
		echo "No version provided, skipping tag"; \
	fi

check-git-clean:
	@if [ -n "$(shell git status --porcelain)" ]; then \
		echo "Error: Working directory is not clean. Commit or stash changes first."; \
		exit 1; \
	fi