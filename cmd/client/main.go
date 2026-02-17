package main

import (
	"context"
	"flag"
	"log"
	"strings"
	"time"

	"grpc-go-fx/internal/client"
	"grpc-go-fx/internal/config"
	"grpc-go-fx/internal/generated/product"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("cannot initialize zap logger: %v", err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	if product.File_product_proto == nil {
		logger.Fatal("Generated proto descriptor not loaded. Run: make generate (requires protoc, protoc-gen-go, protoc-gen-go-grpc). See README.")
	}
	target := flag.String("addr", "localhost:50051", "gRPC server address to dial")
	flag.Parse()

	cfg := &config.Config{
		ClientTarget: *target,
	}

	app := fx.New(
		fx.Supply(cfg),
		fx.Supply(logger),
		client.Module,
		fx.Invoke(RunOrderDemo),
	)
	app.Run()
}

// RunOrderDemo simulates the Order service integrating with the Product API via gRPC.
func RunOrderDemo(c *client.Client, shutdowner fx.Shutdowner, logger *zap.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Info("Order service: fetching product from Product API", zap.String("product_id", "prod-1"))
	p, err := c.GetProduct(ctx, "prod-1")
	if err != nil {
		if status.Code(err) == codes.Unavailable && strings.Contains(err.Error(), "connection refused") {
			logger.Error("GetProduct failed; gRPC server unavailable", zap.Error(err))
			logger.Info("Hint: start the Product gRPC server first (e.g. make run-server in another terminal, or use make run-all).")
		} else {
			logger.Error("GetProduct failed", zap.Error(err))
		}
		_ = shutdowner.Shutdown()
		return
	}
	logger.Info("Product fetched",
		zap.String("id", p.GetId()),
		zap.String("name", p.GetName()),
		zap.Float64("price", p.GetPrice()),
	)

	logger.Info("Order service: listing products", zap.Int32("limit", 2))
	resp, err := c.ListProducts(ctx, 2)
	if err != nil {
		logger.Error("ListProducts failed", zap.Error(err))
		_ = shutdowner.Shutdown()
		return
	}
	for _, p := range resp.GetProducts() {
		logger.Info("Product in list",
			zap.String("id", p.GetId()),
			zap.String("name", p.GetName()),
			zap.Float64("price", p.GetPrice()),
		)
	}
	_ = shutdowner.Shutdown() // exit cleanly after demo
}
