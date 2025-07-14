# Multi-stage build for smaller image
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/httpserver

# Final stage
FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y \
        ca-certificates \
        wget && \
    rm -rf /var/lib/apt/lists/*

RUN groupadd -r authuser && useradd -r -g authuser authuser

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main /app/main
RUN chmod +x /app/main

USER authuser

EXPOSE 5666

CMD ["/app/main"]