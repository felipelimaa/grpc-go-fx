package gateway

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"grpc-go-fx/internal/api"
	"grpc-go-fx/internal/config"
	"grpc-go-fx/internal/generated/product"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
)

func TestNewServeMux_RegistersHandlers(t *testing.T) {
	svc := api.NewProductService()

	mux, err := NewServeMux(svc)
	if err != nil {
		t.Fatalf("NewServeMux returned error: %v", err)
	}
	if mux == nil {
		t.Fatal("expected non-nil ServeMux")
	}
}

func TestGateway_GetProductViaHTTP(t *testing.T) {
	svc := api.NewProductService()
	mux, err := NewServeMux(svc)
	if err != nil {
		t.Fatalf("NewServeMux returned error: %v", err)
	}

	body := `{"id":"prod-1"}`
	req := httptest.NewRequest(http.MethodPost, "/product.v1.ProductService/GetProduct", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status code: got %d, want %d. body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}

	var p product.Product
	if err := json.Unmarshal(rr.Body.Bytes(), &p); err != nil {
		t.Fatalf("failed to unmarshal response body: %v (body=%s)", err, rr.Body.String())
	}
	if p.GetId() != "prod-1" {
		t.Fatalf("unexpected product id: got %q, want %q", p.GetId(), "prod-1")
	}
	if p.GetName() == "" {
		t.Fatalf("expected non-empty product name for id %q", p.GetId())
	}
}

func TestGateway_ListProductsViaHTTP(t *testing.T) {
	svc := api.NewProductService()
	mux, err := NewServeMux(svc)
	if err != nil {
		t.Fatalf("NewServeMux returned error: %v", err)
	}

	body := `{"limit":2}`
	req := httptest.NewRequest(http.MethodPost, "/product.v1.ProductService/ListProducts", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status code: got %d, want %d. body=%s", rr.Code, http.StatusOK, rr.Body.String())
	}

	var resp product.ListProductsResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response body: %v (body=%s)", err, rr.Body.String())
	}

	products := resp.GetProducts()
	if got, want := len(products), 2; got != want {
		t.Fatalf("unexpected number of products: got %d, want %d", got, want)
	}

	for _, p := range products {
		if p.GetId() == "" {
			t.Fatalf("expected non-empty product id in list response")
		}
	}
}

type stubLifecycle struct {
	hooks []fx.Hook
}

func (s *stubLifecycle) Append(h fx.Hook) {
	s.hooks = append(s.hooks, h)
}

func TestRegisterGatewayLifecycle_NoAddrConfigured(t *testing.T) {
	lc := &stubLifecycle{}
	cfg := &config.Config{HTTPGatewayAddr: ""}
	mux := runtime.NewServeMux()

	RegisterGatewayLifecycle(lc, cfg, mux)
	if len(lc.hooks) != 1 {
		t.Fatalf("expected 1 lifecycle hook, got %d", len(lc.hooks))
	}

	hook := lc.hooks[0]
	// With no address configured, OnStart should be a no-op and OnStop should not error.
	if err := hook.OnStart(context.Background()); err != nil {
		t.Fatalf("OnStart returned error for empty HTTPGatewayAddr: %v", err)
	}
	if err := hook.OnStop(context.Background()); err != nil {
		t.Fatalf("OnStop returned error for empty HTTPGatewayAddr: %v", err)
	}
}

func TestRegisterGatewayLifecycle_StartsAndStopsHTTPServer(t *testing.T) {
	lc := &stubLifecycle{}
	// Use an ephemeral port on localhost.
	cfg := &config.Config{HTTPGatewayAddr: "127.0.0.1:0"}
	mux := runtime.NewServeMux()

	RegisterGatewayLifecycle(lc, cfg, mux)
	if len(lc.hooks) != 1 {
		t.Fatalf("expected 1 lifecycle hook, got %d", len(lc.hooks))
	}
	hook := lc.hooks[0]

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := hook.OnStart(ctx); err != nil {
		t.Fatalf("OnStart returned error: %v", err)
	}

	// Stop the HTTP server gracefully; this should not error even if there are no requests.
	if err := hook.OnStop(context.Background()); err != nil {
		t.Fatalf("OnStop returned error: %v", err)
	}
}

