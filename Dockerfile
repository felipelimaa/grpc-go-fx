# Build the Product API image.
#   docker build -t product-api .
ARG GO_VERSION=1.21

FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /app/api ./cmd/api

# Product API (listens on 50051). Override: docker run -p 50052:50052 product-api -addr=:50052
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/api /bin/api
ENTRYPOINT ["/bin/api"]
EXPOSE 50051
CMD ["-addr=:50051"]
