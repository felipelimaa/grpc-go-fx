package client

import (
	"context"
	"testing"

	"grpc-go-fx/internal/config"
	"grpc-go-fx/internal/generated/product"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// fakeProductServiceClient is a test double for product.ProductServiceClient.
type fakeProductServiceClient struct {
	lastGetProductReq      *product.GetProductRequest
	getProductResp         *product.Product
	getProductErr          error
	lastListProductsReq    *product.ListProductsRequest
	listProductsResp       *product.ListProductsResponse
	listProductsErr        error
}

func (f *fakeProductServiceClient) GetProduct(ctx context.Context, in *product.GetProductRequest, opts ...grpc.CallOption) (*product.Product, error) {
	f.lastGetProductReq = in
	return f.getProductResp, f.getProductErr
}

func (f *fakeProductServiceClient) ListProducts(ctx context.Context, in *product.ListProductsRequest, opts ...grpc.CallOption) (*product.ListProductsResponse, error) {
	f.lastListProductsReq = in
	return f.listProductsResp, f.listProductsErr
}

type stubLifecycle struct {
	hooks []fx.Hook
}

func (s *stubLifecycle) Append(h fx.Hook) {
	s.hooks = append(s.hooks, h)
}

// newTestConn creates a lightweight gRPC client connection that can be closed
// without needing a running server. It uses a custom dialer that never dials.
func newTestConn(t *testing.T) *grpc.ClientConn {
	t.Helper()

	conn, err := grpc.Dial(
		"passthrough:///ignored",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to create test gRPC connection: %v", err)
	}
	return conn
}

func TestClientGetProductDelegatesToUnderlyingClient(t *testing.T) {
	fake := &fakeProductServiceClient{
		getProductResp: &product.Product{Id: "prod-1", Name: "Widget A"},
	}
	c := &Client{ProductServiceClient: fake}

	ctx := context.Background()
	got, err := c.GetProduct(ctx, "prod-1")
	if err != nil {
		t.Fatalf("GetProduct returned error: %v", err)
	}
	if got == nil {
		t.Fatal("GetProduct returned nil product")
	}
	if got.GetId() != "prod-1" {
		t.Fatalf("GetProduct returned wrong id: got %q, want %q", got.GetId(), "prod-1")
	}

	if fake.lastGetProductReq == nil {
		t.Fatal("underlying client did not receive request")
	}
	if fake.lastGetProductReq.GetId() != "prod-1" {
		t.Fatalf("underlying client received wrong id: got %q, want %q", fake.lastGetProductReq.GetId(), "prod-1")
	}
}

func TestClientListProductsDelegatesToUnderlyingClient(t *testing.T) {
	fake := &fakeProductServiceClient{
		listProductsResp: &product.ListProductsResponse{
			Products: []*product.Product{
				{Id: "prod-1"},
				{Id: "prod-2"},
			},
		},
	}
	c := &Client{ProductServiceClient: fake}

	ctx := context.Background()
	resp, err := c.ListProducts(ctx, 2)
	if err != nil {
		t.Fatalf("ListProducts returned error: %v", err)
	}
	if resp == nil {
		t.Fatal("ListProducts returned nil response")
	}

	if got, want := int32(len(resp.GetProducts())), int32(2); got != want {
		t.Fatalf("ListProducts returned wrong number of products: got %d, want %d", got, want)
	}

	if fake.lastListProductsReq == nil {
		t.Fatal("underlying client did not receive request")
	}
	if fake.lastListProductsReq.GetLimit() != 2 {
		t.Fatalf("underlying client received wrong limit: got %d, want %d", fake.lastListProductsReq.GetLimit(), 2)
	}
}

func TestNewConnAndRegisterClientLifecycle(t *testing.T) {
	cfg := &config.Config{
		ClientTarget: "localhost:0",
	}

	conn, err := NewConn(cfg)
	if err != nil {
		t.Fatalf("NewConn returned error: %v", err)
	}
	if conn == nil {
		t.Fatal("NewConn returned nil connection")
	}

	lc := &stubLifecycle{}
	RegisterClientLifecycle(lc, conn)
	if len(lc.hooks) != 1 {
		t.Fatalf("expected 1 lifecycle hook, got %d", len(lc.hooks))
	}

	// Ensure the OnStop hook closes the connection without error.
	if err := lc.hooks[0].OnStop(context.Background()); err != nil {
		t.Fatalf("OnStop returned error: %v", err)
	}
}

func TestNewClientWrapsGeneratedClient(t *testing.T) {
	// We only verify that NewClient returns a non-nil wrapper and does not panic
	// when constructing the generated ProductServiceClient.
	conn := newTestConn(t)
	defer conn.Close()

	c := NewClient(conn)
	if c == nil {
		t.Fatal("expected non-nil Client from NewClient")
	}
}

