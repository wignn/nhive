package grpcauth

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const metaKey = "x-internal-key"


func UnaryServerInterceptor(apiKey string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if apiKey == "" {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		vals := md.Get(metaKey)
		if len(vals) == 0 || vals[0] != apiKey {
			return nil, status.Error(codes.Unauthenticated, "invalid internal key")
		}

		return handler(ctx, req)
	}
}

// credentials implements grpc.PerRPCCredentials.
// It injects the internal API key into every outgoing gRPC call's metadata.
type credentials struct {
	key string
}

// NewCredentials returns a PerRPCCredentials that injects x-internal-key into
// every outgoing gRPC call. Pass it to grpc.WithPerRPCCredentials().
// If key is empty, no metadata is injected.
func NewCredentials(key string) credentials {
	return credentials{key: key}
}

func (c credentials) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	if c.key == "" {
		return nil, nil
	}
	return map[string]string{metaKey: c.key}, nil
}

// RequireTransportSecurity returns false because we use insecure gRPC inside
// the Docker private network. If you add TLS, set this to true.
func (c credentials) RequireTransportSecurity() bool {
	return false
}
