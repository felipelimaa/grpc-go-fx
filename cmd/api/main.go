package main

import (
	"flag"

	"grpc-go-fx/internal/api"
	"grpc-go-fx/internal/config"
	"grpc-go-fx/internal/gateway"

	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func main() {
	addr := flag.String("addr", ":50051", "gRPC API listen address")
	httpAddr := flag.String("http-addr", ":8080", "HTTP/JSON gateway listen address (grpc-gateway)")
	flag.Parse()

	cfg := &config.Config{
		ServerAddr:      *addr,
		HTTPGatewayAddr: *httpAddr,
	}

	app := fx.New(
		fx.Supply(cfg),
		api.Module,
		gateway.Module,
		fx.Invoke(func(*grpc.Server) {}), // ensure API server is built and lifecycle runs
	)
	app.Run()
}

