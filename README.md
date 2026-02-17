# gRPC Go Product API

A Go project that exposes a **Product API** over gRPC and HTTP/JSON (via grpc-gateway). Uses **Uber FX** for dependency injection and lifecycle management.

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

This runs `scripts/gen.sh`, which uses `protoc` to generate:

- `internal/generated/product/*.pb.go` – gRPC types and service
- `internal/generated/product/product.pb.gw.go` – grpc-gateway HTTP/JSON bindings

**Run this before using the client or HTTP gateway** so request/response marshaling works correctly.

## Build and run

**Start the Product API (gRPC server + HTTP gateway)**

```bash
go build -o bin/api ./cmd/api
./bin/api -addr=:50051 -http-addr=:8080
```

## Unit tests

You can run the unit tests against the core, handwritten packages (API, gateway, config) with:

```bash
make test
```

To run tests with coverage enabled and see a per-function breakdown:

```bash
make test-cover
```

This uses Go’s built-in coverage tooling and excludes generated code under `internal/generated/product` so that coverage reflects only the application logic you maintain.

## HTTP/JSON gateway and OpenAPI

In addition to pure gRPC, this project exposes the `ProductService` over HTTP/JSON using **grpc-gateway**.

- **Gateway address**: `http://localhost:8080` (configurable via `-http-addr`)
- **OpenAPI spec**: `api/product/openapi.yaml`
  - Documents the HTTP endpoints that proxy to the gRPC methods:
    - `POST /product.v1.ProductService/GetProduct`
    - `POST /product.v1.ProductService/ListProducts`

### Test the API via HTTP with curl

Make sure the Product API is running:

```bash
./bin/api -addr=:50051 -http-addr=:8080
```

**Get a single product by ID**:

```bash
curl -X POST http://localhost:8080/product.v1.ProductService/GetProduct \
  -H "Content-Type: application/json" \
  -d '{
    "id": "prod-1"
  }'
```

**List products (limit 2)**:

```bash
curl -X POST http://localhost:8080/product.v1.ProductService/ListProducts \
  -H "Content-Type: application/json" \
  -d '{
    "limit": 2
  }'
```

You can omit the body or send `{}` to use the server’s default limit.

### Test the API via OpenAPI (Postman / Insomnia)

1. Start the Product API with the gateway:

   ```bash
   ./bin/api -addr=:50051 -http-addr=:8080
   ```

2. In Postman (or Insomnia):
   - **Import** `api/product/openapi.yaml` as an OpenAPI definition.
   - Ensure the server URL is `http://localhost:8080`.

3. Use the generated requests:
   - **GetProduct**:
     - Endpoint: `POST /product.v1.ProductService/GetProduct`
     - Body (JSON):

       ```json
       {
         "id": "prod-1"
       }
       ```

   - **ListProducts**:
     - Endpoint: `POST /product.v1.ProductService/ListProducts`
     - Body (JSON), for example:

       ```json
       {
         "limit": 2
       }
       ```

   - Send the requests and inspect the JSON responses.

## Project layout

- `api/product/product.proto` – Product service and messages
- `internal/config` – Product API configuration (supplied via FX)
- `internal/generated/product` – Generated Go from proto (run `make generate`)
- `api/product/openapi.yaml` – OpenAPI 3 spec for the HTTP/JSON gateway
- `internal/gateway` – grpc-gateway HTTP/JSON server wired into FX
- `internal/api` – Product API implementation + gRPC server constructor + FX module
- `cmd/api` – Product API entrypoint (FX app)

## Documentation

See [docs/INTEGRATION.md](docs/INTEGRATION.md) for architecture, gRPC/FX overview, and how to extend the integration.
