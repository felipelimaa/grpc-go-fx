package api

import (
	"context"
	"net"
	"testing"
	"time"

	"grpc-go-fx/internal/config"
	"grpc-go-fx/internal/generated/product"

	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func TestNewProductServiceSeedsStore(t *testing.T) {
	svc := NewProductService()
	if svc == nil {
		t.Fatal("expected non-nil ProductService")
	}

	// The service should be seeded with three known products.
	if got, want := len(svc.store), 3; got != want {
		t.Fatalf("unexpected number of seeded products: got %d, want %d", got, want)
	}

	for _, id := range []string{"prod-1", "prod-2", "prod-3"} {
		p, ok := svc.store[id]
		if !ok {
			t.Fatalf("expected product with id %q to be present", id)
		}
		if p.GetId() != id {
			t.Fatalf("product id mismatch: got %q, want %q", p.GetId(), id)
		}
		if p.GetName() == "" {
			t.Fatalf("expected non-empty name for product %q", id)
		}
	}
}

func TestProductServiceGetProduct_Found(t *testing.T) {
	svc := NewProductService()
	ctx := context.Background()

	got, err := svc.GetProduct(ctx, &product.GetProductRequest{Id: "prod-1"})
	if err != nil {
		t.Fatalf("GetProduct returned error: %v", err)
	}
	if got == nil {
		t.Fatal("GetProduct returned nil product for existing id")
	}
	if got.GetId() != "prod-1" {
		t.Fatalf("GetProduct returned wrong id: got %q, want %q", got.GetId(), "prod-1")
	}
}

func TestProductServiceGetProduct_NotFound(t *testing.T) {
	svc := NewProductService()
	ctx := context.Background()

	got, err := svc.GetProduct(ctx, &product.GetProductRequest{Id: "unknown"})
	if err != nil {
		t.Fatalf("GetProduct returned error: %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil product for unknown id, got %+v", got)
	}
}

func TestProductServiceListProducts_DefaultLimit(t *testing.T) {
	svc := NewProductService()
	ctx := context.Background()

	resp, err := svc.ListProducts(ctx, &product.ListProductsRequest{Limit: 0})
	if err != nil {
		t.Fatalf("ListProducts returned error: %v", err)
	}
	if resp == nil {
		t.Fatal("ListProducts returned nil response")
	}

	products := resp.GetProducts()
	if got, want := len(products), len(svc.store); got != want {
		t.Fatalf("unexpected number of products: got %d, want %d", got, want)
	}

	// Ensure all returned IDs exist in the store.
	for _, p := range products {
		if _, ok := svc.store[p.GetId()]; !ok {
			t.Fatalf("ListProducts returned unknown product id %q", p.GetId())
		}
	}
}

func TestProductServiceListProducts_WithLimit(t *testing.T) {
	svc := NewProductService()
	ctx := context.Background()

	resp, err := svc.ListProducts(ctx, &product.ListProductsRequest{Limit: 2})
	if err != nil {
		t.Fatalf("ListProducts returned error: %v", err)
	}
	if resp == nil {
		t.Fatal("ListProducts returned nil response")
	}

	products := resp.GetProducts()
	if got, want := len(products), 2; got != want {
		t.Fatalf("unexpected number of products with limit applied: got %d, want %d", got, want)
	}

	for _, p := range products {
		if _, ok := svc.store[p.GetId()]; !ok {
			t.Fatalf("ListProducts returned unknown product id %q", p.GetId())
		}
	}
}

func TestNewGRPCServerRegistersProductService(t *testing.T) {
	cfg := &config.Config{}
	svc := NewProductService()

	srv := NewGRPCServer(cfg, svc)
	if srv == nil {
		t.Fatal("expected non-nil gRPC server")
	}

	info := srv.GetServiceInfo()
	if _, ok := info["product.v1.ProductService"]; !ok {
		t.Fatalf("ProductService not registered on gRPC server; services: %v", info)
	}
}

type stubLifecycle struct {
	hooks []fx.Hook
}

func (s *stubLifecycle) Append(h fx.Hook) {
	s.hooks = append(s.hooks, h)
}

func TestRegisterGRPCLifecycle_StartsAndStopsServer(t *testing.T) {
	// Use an ephemeral port on localhost to avoid collisions.
	lc := &stubLifecycle{}
	cfg := &config.Config{ServerAddr: "127.0.0.1:0"}

	svc := NewProductService()
	srv := grpc.NewServer()
	product.RegisterProductServiceServer(srv, svc)

	RegisterGRPCLifecycle(lc, srv, cfg)
	if len(lc.hooks) != 1 {
		t.Fatalf("expected 1 lifecycle hook, got %d", len(lc.hooks))
	}

	hook := lc.hooks[0]

	// Start the gRPC server.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Verify that we can listen on the configured address before OnStart runs.
	// This also ensures the address string is valid.
	l, err := net.Listen("tcp", cfg.ServerAddr)
	if err != nil {
		t.Fatalf("failed to listen on %q: %v", cfg.ServerAddr, err)
	}
	l.Close()

	if err := hook.OnStart(ctx); err != nil {
		t.Fatalf("OnStart returned error: %v", err)
	}

	// Stop the server gracefully.
	if err := hook.OnStop(context.Background()); err != nil {
		t.Fatalf("OnStop returned error: %v", err)
	}
}

