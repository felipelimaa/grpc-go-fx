package config

// Config holds addresses for the Product API.
type Config struct {
	// ServerAddr is the listen address for the gRPC API server (e.g. ":50051").
	ServerAddr string
	// HTTPGatewayAddr is the listen address for the HTTP/JSON gateway (e.g. ":8080").
	HTTPGatewayAddr string
}
