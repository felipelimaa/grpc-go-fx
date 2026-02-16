package client

import (
	"context"

	"grpc-go-fx/internal/config"
	"grpc-go-fx/internal/generated/product"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client wraps the Product gRPC client for use by the Order "API" integration.
type Client struct {
	product.ProductServiceClient
}

// NewConn creates a gRPC client connection to the Product service.
func NewConn(cfg *config.Config) (*grpc.ClientConn, error) {
	return grpc.Dial(cfg.ClientTarget, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

// NewClient creates a Client wrapper around the generated gRPC client.
func NewClient(conn *grpc.ClientConn) *Client {
	return &Client{ProductServiceClient: product.NewProductServiceClient(conn)}
}

// GetProduct delegates to the generated client.
func (c *Client) GetProduct(ctx context.Context, id string) (*product.Product, error) {
	return c.ProductServiceClient.GetProduct(ctx, &product.GetProductRequest{Id: id})
}

// ListProducts delegates to the generated client.
func (c *Client) ListProducts(ctx context.Context, limit int32) (*product.ListProductsResponse, error) {
	return c.ProductServiceClient.ListProducts(ctx, &product.ListProductsRequest{Limit: limit})
}
