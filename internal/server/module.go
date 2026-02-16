package server

import (
	"context"
	"net"

	"grpc-go-fx/internal/config"
	"grpc-go-fx/internal/generated/product"

	"go.uber.org/fx"
	"google.golang.org/grpc"
)

// Module is the FX module for the Product gRPC server.
var Module = fx.Module("server",
	fx.Provide(fx.Annotate(NewProductService, fx.As(new(product.ProductServiceServer)))),
	fx.Provide(NewGRPCServer),
	fx.Invoke(RegisterGRPCLifecycle),
)

// RegisterGRPCLifecycle registers the gRPC server with FX lifecycle (OnStart listen/serve, OnStop GracefulStop).
func RegisterGRPCLifecycle(lc fx.Lifecycle, srv *grpc.Server, cfg *config.Config) {
	var lis net.Listener
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			var err error
			lis, err = net.Listen("tcp", cfg.ServerAddr)
			if err != nil {
				return err
			}
			go func() { _ = srv.Serve(lis) }()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			srv.GracefulStop()
			return nil
		},
	})
}
