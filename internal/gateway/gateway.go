package gateway

import (
	"context"
	"net/http"

	"grpc-go-fx/internal/config"
	"grpc-go-fx/internal/generated/product"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
)

// Module wires the HTTP/JSON gateway for the ProductService into the FX app lifecycle.
//
// It exposes gRPC methods over HTTP using the default grpc-gateway mappings with
// generate_unbound_methods=true, which results in POST endpoints like:
//   - POST /product.v1.ProductService/GetProduct
//   - POST /product.v1.ProductService/ListProducts
var Module = fx.Module("gateway",
	fx.Provide(NewServeMux),
	fx.Invoke(RegisterGatewayLifecycle),
)

// NewServeMux builds a grpc-gateway ServeMux and registers the ProductService handlers.
func NewServeMux(svc product.ProductServiceServer) (*runtime.ServeMux, error) {
	mux := runtime.NewServeMux()
	ctx := context.Background()

	// Register handlers that translate HTTP/JSON requests into gRPC calls
	// handled by the in-process ProductServiceServer implementation.
	if err := product.RegisterProductServiceHandlerServer(ctx, mux, svc); err != nil {
		return nil, err
	}

	return mux, nil
}

// RegisterGatewayLifecycle starts and stops the HTTP gateway with the FX lifecycle.
func RegisterGatewayLifecycle(lc fx.Lifecycle, cfg *config.Config, mux *runtime.ServeMux) {
	var srv *http.Server

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if cfg.HTTPGatewayAddr == "" {
				// If no HTTP gateway address is configured, do not start the HTTP server.
				return nil
			}

			srv = &http.Server{
				Addr:    cfg.HTTPGatewayAddr,
				Handler: mux,
			}

			go func() {
				// http.ErrServerClosed is expected on graceful shutdown.
				_ = srv.ListenAndServe()
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			if srv == nil {
				return nil
			}
			return srv.Shutdown(ctx)
		},
	})
}

