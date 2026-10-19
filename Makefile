.DEFAULT_GOAL := help

MODULE  := auth302
CMD     := ./cmd/httpserver
BINARY  := bin/auth302
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -X '$(MODULE)/internal/infrastructure/httpserver.Version=$(VERSION)'

.PHONY: help build fmt fmt-ci vet test start kill

help:
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  %-18s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build server binary into bin/auth302
	@mkdir -p $(dir $(BINARY))
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) $(CMD)

fmt: ## Run go fmt
	go fmt ./...

fmt-ci: ## Fail if any file is not gofmt-clean
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "Code formatting check failed. Run 'make fmt' to fix formatting issues."; \
		exit 1; \
	fi

vet: ## Run go vet
	go vet ./...

test: ## Run unit tests
	go test ./internal/... ./pkg/...

start: ## Start the local Docker stack
	docker compose -f compose.yml up --build -d

kill: ## Stop the local Docker stack and remove volumes
	docker compose -f compose.yml down -v
