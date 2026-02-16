# gRPC Go API Integration

A Go project that simulates integration between two APIs over gRPC: an **Order service** (client) calls a **Product service** (server) to fetch product details. Uses **Uber FX** for dependency injection and lifecycle management.

## Prerequisites

- Go 1.21+
- For code generation: **protoc**, **protoc-gen-go**, **protoc-gen-go-grpc**
  - [Protocol Buffers release](https://github.com/protocolbuffers/protobuf/releases) (or `brew install protobuf` on macOS)
  - `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
  - `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`

## Code generation

After cloning or editing `.proto` files, regenerate Go code:

```bash
make generate
```

This runs `scripts/gen.sh`, which uses `protoc` to generate `internal/generated/product/*.pb.go`. **Run this before using the client** so request/response marshaling works correctly.

## Build and run

**Terminal 1 – Product gRPC server**

```bash
go build -o bin/server ./cmd/server
./bin/server -addr=:50051
```

**Terminal 2 – Order client (calls Product API)**

```bash
go build -o bin/client ./cmd/client
./bin/client -addr=localhost:50051
```

The client simulates the Order service: it calls `GetProduct("prod-1")` and `ListProducts(2)` and prints the results.

## Project layout

- `api/product/product.proto` – Product service and messages
- `internal/config` – Server/client address config (supplied via FX)
- `internal/generated/product` – Generated Go from proto (run `make generate`)
- `internal/server` – Product gRPC server + FX module
- `internal/client` – gRPC client wrapper + FX module
- `cmd/server` – Server entrypoint (FX app)
- `cmd/client` – Client entrypoint (FX app, runs demo then exits)

## Documentation

See [docs/INTEGRATION.md](docs/INTEGRATION.md) for architecture, gRPC/FX overview, and how to extend the integration.
