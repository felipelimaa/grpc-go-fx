package config

// Config holds server and client gRPC addresses for the Product API integration.
type Config struct {
	// ServerAddr is the listen address for the gRPC server (e.g. ":50051").
	ServerAddr string
	// ClientTarget is the target address for the gRPC client to dial (e.g. "localhost:50051").
	ClientTarget string
}
