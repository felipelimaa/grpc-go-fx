package main

import (
	"flag"

	"grpc-go-fx/internal/config"
	"grpc-go-fx/internal/server"

	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func main() {
	addr := flag.String("addr", ":50051", "gRPC server listen address")
	flag.Parse()

	cfg := &config.Config{ServerAddr: *addr}

	app := fx.New(
		fx.Supply(cfg),
		server.Module,
		fx.Invoke(func(*grpc.Server) {}), // ensure server is built and lifecycle runs
	)
	app.Run()
}
