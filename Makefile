.PHONY: build test check clean install-tools install-hooks

build:
	go build -o req42-tracer ./cmd/req42-tracer/

test:
	go test ./...

test-race:
	go test -race ./...

vet:
	go vet ./...

staticcheck:
	staticcheck ./...

gosec:
	gosec ./...

check: vet staticcheck test-race
	@echo "✓ All checks passed"

clean:
	rm -f req42-tracer

install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest

install-hooks:
	cp scripts/pre-commit .git/hooks/ 2>/dev/null || echo "No pre-commit script yet"
	chmod +x .git/hooks/pre-commit

run:
	go run ./cmd/req42-tracer/

help:
	@grep -E '^[a-zA-Z_-]+:' $(MAKEFILE_LIST) | sed 's/:.*//' | sort
