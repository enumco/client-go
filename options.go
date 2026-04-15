package client

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type options struct {
	addr     string
	creds    credentials.TransportCredentials
	dialOpts []grpc.DialOption
}

type ClientOption func(*options)

// WithAddr overrides the default API endpoint (api.enum.co:443).
func WithAddr(addr string) ClientOption {
	return func(o *options) {
		o.addr = addr
	}
}

// WithInsecure disables TLS. For local development only.
func WithInsecure() ClientOption {
	return func(o *options) {
		o.creds = insecure.NewCredentials()
	}
}

// WithTLS overrides the default TLS configuration (system roots).
func WithTLS(creds credentials.TransportCredentials) ClientOption {
	return func(o *options) {
		o.creds = creds
	}
}

func WithToken(token string) ClientOption {
	return func(o *options) {
		o.dialOpts = append(o.dialOpts, grpc.WithPerRPCCredentials(bearerToken(token)))
	}
}

type bearerToken string

func (t bearerToken) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	return map[string]string{"authorization": "Bearer " + string(t)}, nil
}

func (t bearerToken) RequireTransportSecurity() bool {
	return false
}
