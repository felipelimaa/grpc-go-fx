.PHONY: install-tools generate run-server run-client run-all test test-cover
# Server listens on SERVER_PORT; client connects to CLIENT_TARGET (default: same port).
SERVER_PORT ?= 50051
SERVER_ADDR ?= :$(SERVER_PORT)
CLIENT_TARGET ?= localhost:$(SERVER_PORT)
# Ensure protoc can find Go-installed plugins (protoc-gen-go, protoc-gen-go-grpc)
export PATH := $(shell go env GOPATH)/bin:$(PATH)

# Install protoc Go plugins (required for make generate). Need protoc on PATH (e.g. brew install protobuf).
install-tools:
	@command -v protoc-gen-go >/dev/null 2>&1 || go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@command -v protoc-gen-go-grpc >/dev/null 2>&1 || go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@command -v protoc-gen-grpc-gateway >/dev/null 2>&1 || go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	@echo "protoc-gen-go, protoc-gen-go-grpc and protoc-gen-grpc-gateway installed. Ensure GOPATH/bin (or GOBIN) is in your PATH."

generate: install-tools
	@./scripts/gen.sh

# Start gRPC server (listens on SERVER_PORT). Override: make run-server SERVER_PORT=50052
run-server:
	@go run ./cmd/server -addr=$(SERVER_ADDR)

# Run gRPC client (one-shot). Connects to CLIENT_TARGET. Start server first: make run-server (or use make run-all).
run-client:
	@go run ./cmd/client -addr=$(CLIENT_TARGET)

# Run both APIs: server in background on SERVER_PORT, then client. Stop server with Ctrl+C or kill the server process.
run-all:
	@go run ./cmd/server -addr=$(SERVER_ADDR) & \
	pid=$$!; sleep 1; go run ./cmd/client -addr=$(CLIENT_TARGET); kill $$pid 2>/dev/null || true

# Run unit tests for core handwritten packages with coverage enabled.
test:
	@go test ./internal/server ./internal/client ./internal/gateway ./internal/config -cover

# Run unit tests with coverage profile and print per-function coverage.
test-cover:
	@go test ./internal/server ./internal/client ./internal/gateway ./internal/config -coverprofile=coverage.out
	@go tool cover -func=coverage.out
