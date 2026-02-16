package client

import (
	"context"

	"go.uber.org/fx"
	"google.golang.org/grpc"
)

// Module is the FX module for the Product gRPC client (Order "API" integration).
var Module = fx.Module("client",
	fx.Provide(NewConn),
	fx.Provide(NewClient),
	fx.Invoke(RegisterClientLifecycle),
)

// RegisterClientLifecycle closes the gRPC connection on app shutdown.
func RegisterClientLifecycle(lc fx.Lifecycle, conn *grpc.ClientConn) {
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return conn.Close()
		},
	})
}
