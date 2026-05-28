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
	$(GO_BIN) build -o req42-tracer ./cmd/req42-tracer/

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
	rm -f req42-tracer

install-tools:
	$(GO_BIN) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO_BIN) install honnef.co/go/tools/cmd/staticcheck@latest
	$(GO_BIN) install github.com/securego/gosec/v2/cmd/gosec@latest

install-hooks:
	cp scripts/pre-commit .git/hooks/ 2>/dev/null || echo "No pre-commit script yet"
	chmod +x .git/hooks/pre-commit

run:
	$(GO_BIN) run ./cmd/req42-tracer/ $(ARGS)

help:
	@grep -E '^[a-zA-Z_-]+:' $(MAKEFILE_LIST) | sed 's/:.*//' | sort
