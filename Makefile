start:
	go run cmd/httpserver/main.go

build:
	go build -o aegix cmd/httpserver/main.go

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run

test:
	go test ./...