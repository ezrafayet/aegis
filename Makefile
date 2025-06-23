start:
	go run cmd/httpserver/main.go

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

lint:
	golangci-lint run

test:
	go test ./...