build:
	go build -o main cmd/httpserver/main.go

fmt:
	go fmt ./...

fmt-ci:
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "Code formatting check failed. Run 'make fmt' to fix formatting issues."; \
		exit 1; \
	fi

vet:
	go vet ./...

test:
	make test-unit
	make test-integration

test-integration:
	go test -tags=integration -timeout=120s ./integration-tests/...

test-unit:
	go test -timeout=15s ./internal/... ./pkg/...
