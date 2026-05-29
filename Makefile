.PHONY: build test check clean install-tools install-hooks

GO_VERSION := 1.25.4
GO_INSTALL_DIR := /usr/local
GO_BIN := $(GO_INSTALL_DIR)/go/bin/go

install-go:
	wget https://go.dev/dl/go$(GO_VERSION).linux-amd64.tar.gz
	sudo rm -rf $(GO_INSTALL_DIR)/go
	sudo tar -C $(GO_INSTALL_DIR) -xzf go$(GO_VERSION).linux-amd64.tar.gz
	$(GO_BIN) version

build:
	$(GO_BIN) build -o bin/req42-tracer ./src/cmd/req42-tracer/

build-all:
	CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 $(GO_BIN) build -o bin/linux/amd64/req42-tracer     ./src/cmd/req42-tracer/
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO_BIN) build -o bin/windows/amd64/req42-tracer.exe ./src/cmd/req42-tracer/
	CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 $(GO_BIN) build -o bin/darwin/arm64/req42-tracer     ./src/cmd/req42-tracer/

test:
	$(GO_BIN) test ./...

test-race:
	$(GO_BIN) test -race ./...

vet:
	$(GO_BIN) vet ./...

staticcheck:
	staticcheck ./...

gosec:
	gosec ./...

check: vet staticcheck test-race
	@echo "✓ All checks passed"

clean:
	rm -rf bin/ req42-tracer

install-tools:
	$(GO_BIN) install honnef.co/go/tools/cmd/staticcheck@v0.6.0
	$(GO_BIN) install github.com/securego/gosec/v2/cmd/gosec@v2.22.4

install-hooks:
	cp project/req42-tracer/scripts/pre-commit .git/hooks/ 2>/dev/null || true
	chmod +x .git/hooks/pre-commit

run:
	$(GO_BIN) run ./src/cmd/req42-tracer/ $(ARGS)

help:
	@grep -E '^[a-zA-Z_-]+:' $(MAKEFILE_LIST) | sed 's/:.*//' | sort
