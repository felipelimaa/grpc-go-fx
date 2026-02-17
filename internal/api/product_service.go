package api

import (
	"context"
	"sync"

	"grpc-go-fx/internal/config"
	"grpc-go-fx/internal/generated/product"

	"google.golang.org/grpc"
)

// ProductService implements product.ProductServiceServer with in-memory storage.
type ProductService struct {
	product.UnimplementedProductServiceServer
	mu    sync.RWMutex
	store map[string]*product.Product
}

// NewProductService creates a ProductService with seeded product data.
func NewProductService() *ProductService {
	store := map[string]*product.Product{
		"prod-1": {Id: "prod-1", Name: "Widget A", Description: "A useful widget", Price: 9.99},
		"prod-2": {Id: "prod-2", Name: "Gadget B", Description: "A handy gadget", Price: 19.99},
		"prod-3": {Id: "prod-3", Name: "Gizmo C", Description: "A small gizmo", Price: 4.99},
	}
	return &ProductService{store: store}
}

// GetProduct returns a product by ID.
func (s *ProductService) GetProduct(ctx context.Context, req *product.GetProductRequest) (*product.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if p, ok := s.store[req.GetId()]; ok {
		return p, nil
	}
	return nil, nil // not found: return empty (or use status.NotFound in production)
}

// ListProducts returns products up to the given limit.
func (s *ProductService) ListProducts(ctx context.Context, req *product.ListProductsRequest) (*product.ListProductsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	limit := req.GetLimit()
	if limit <= 0 {
		limit = 10
	}
	var list []*product.Product
	for _, p := range s.store {
		list = append(list, p)
		if int32(len(list)) >= limit {
			break
		}
	}
	return &product.ListProductsResponse{Products: list}, nil
}

// NewGRPCServer creates a gRPC server with the Product service registered.
func NewGRPCServer(cfg *config.Config, svc product.ProductServiceServer) *grpc.Server {
	srv := grpc.NewServer()
	product.RegisterProductServiceServer(srv, svc)
	return srv
}

