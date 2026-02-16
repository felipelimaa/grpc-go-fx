package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"grpc-go-fx/internal/client"
	"grpc-go-fx/internal/config"
	"grpc-go-fx/internal/generated/product"

	"go.uber.org/fx"
)

func main() {
	if product.File_api_product_product_proto == nil {
		log.Fatal("Generated proto descriptor not loaded. Run: make generate (requires protoc, protoc-gen-go, protoc-gen-go-grpc). See README.")
	}
	target := flag.String("addr", "localhost:50051", "gRPC server address to dial")
	flag.Parse()

	cfg := &config.Config{ClientTarget: *target}

	app := fx.New(
		fx.Supply(cfg),
		client.Module,
		fx.Invoke(RunOrderDemo),
	)
	app.Run()
}

// RunOrderDemo simulates the Order service integrating with the Product API via gRPC.
func RunOrderDemo(c *client.Client, shutdowner fx.Shutdowner) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Fprintln(os.Stderr, "Order service: fetching product prod-1 from Product API...")
	p, err := c.GetProduct(ctx, "prod-1")
	if err != nil {
		log.Printf("GetProduct: %v", err)
		_ = shutdowner.Shutdown()
		return
	}
	fmt.Printf("Product: id=%s name=%s price=%.2f\n", p.GetId(), p.GetName(), p.GetPrice())

	fmt.Fprintln(os.Stderr, "Order service: listing products (limit=2)...")
	resp, err := c.ListProducts(ctx, 2)
	if err != nil {
		log.Printf("ListProducts: %v", err)
		_ = shutdowner.Shutdown()
		return
	}
	for _, p := range resp.GetProducts() {
		fmt.Printf("  - %s: %s (%.2f)\n", p.GetId(), p.GetName(), p.GetPrice())
	}
	_ = shutdowner.Shutdown() // exit cleanly after demo
}
