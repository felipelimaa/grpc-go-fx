.PHONY: install-tools generate run-api test test-cover
# API listens on API_PORT.
API_PORT ?= 50051
API_ADDR ?= :$(API_PORT)
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

# Start Product API gRPC server (listens on API_PORT). Override: make run-api API_PORT=50052
run-api:
	@go run ./cmd/api -addr=$(API_ADDR)

# Run unit tests for core handwritten packages with coverage enabled.
test:
	@go test ./internal/api ./internal/gateway ./internal/config -cover

# Run unit tests with coverage profile and print per-function coverage.
test-cover:
	@go test ./internal/api ./internal/gateway ./internal/config -coverprofile=coverage.out
	@go tool cover -func=coverage.out
