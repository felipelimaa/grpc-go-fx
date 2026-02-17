#!/usr/bin/env bash
# Generate Go code from product.proto. Requires protoc, protoc-gen-go, protoc-gen-go-grpc.
# Install: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
#          go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
#          protoc: https://protobuf.dev/downloads/ or brew install protobuf
set -e
cd "$(dirname "$0")/.."
mkdir -p internal/generated/product
protoc --go_out=internal/generated/product --go_opt=paths=source_relative \
  --go-grpc_out=internal/generated/product --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=internal/generated/product --grpc-gateway_opt=paths=source_relative,generate_unbound_methods=true \
  -I api/product \
  api/product/product.proto
echo "Generated internal/generated/product/*.pb.go and product.pb.gw.go"
