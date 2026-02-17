# Build each API in one Dockerfile. Use --target to pick the image.
#   docker build --target server -t grpc-server .
#   docker build --target client -t grpc-client .
ARG GO_VERSION=1.21

FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /app/server ./cmd/server && \
    CGO_ENABLED=0 go build -o /app/client ./cmd/client

# Server API (listens on 50051). Override: docker run -p 50052:50052 grpc-server -addr=:50052
FROM scratch AS server
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/server /bin/api
ENTRYPOINT ["/bin/api"]
EXPOSE 50051
CMD ["-addr=:50051"]

# Client API (connects to server:50051). Override: docker run grpc-client -addr=host.docker.internal:50051
FROM scratch AS client
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/client /bin/api
ENTRYPOINT ["/bin/api"]
CMD ["-addr=server:50051"]
