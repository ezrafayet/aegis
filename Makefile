.ONESHELL:
SHELL := /bin/bash

# Set working directory to src
.PHONY: build fmt fmt-ci vet test test-unit test-integration

VERSION := $(shell git describe --tags)

build:
	cd src
	go build -ldflags="-X 'aegis/internal/infrastructure/httpserver.Version=$(VERSION)'" -o main cmd/httpserver/main.go

fmt:
	cd src
	go fmt ./...

fmt-ci:
	cd src
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "Code formatting check failed. Run 'make fmt' to fix formatting issues."; \
		exit 1; \
	fi

vet:
	cd src
	go vet ./...

test:
	make test-unit
	make test-integration

test-integration:
	cd src
	go test -tags=integration -timeout=240s ./integration/...

test-unit:
	cd src
	go test -timeout=15s ./internal/... ./pkg/...

# For dev

start:
	sudo docker-compose -f docker-compose-dev.yml up --build

kill:
	sudo docker-compose -f docker-compose-dev.yml down -v