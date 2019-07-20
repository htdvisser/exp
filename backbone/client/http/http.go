package http

import (
	"context"
	"net/http"
)

type roundTripperCtxKeyType struct{}

var roundTripperCtxKey roundTripperCtxKeyType

// NewContextWithRoundTripper returns a context derived from parent that contains
// roundTripper.
func NewContextWithRoundTripper(parent context.Context, roundTripper http.RoundTripper) context.Context {
	return context.WithValue(parent, roundTripperCtxKey, roundTripper)
}

// RoundTripperFromContext returns the HTTP RoundTripper from the context. If there
// is no RoundTripper in the context, the default RoundTripper is returned.
func RoundTripperFromContext(ctx context.Context) http.RoundTripper {
	if roundTripper, ok := ctx.Value(roundTripperCtxKey).(http.RoundTripper); ok {
		return roundTripper
	}
	return http.DefaultTransport
}

type clientCtxKeyType struct{}

var clientCtxKey clientCtxKeyType

// NewContextWithClient returns a context derived from parent that contains client.
func NewContextWithClient(parent context.Context, client *http.Client) context.Context {
	return context.WithValue(parent, clientCtxKey, client)
}

// ClientFromContext returns the HTTP Client from the context. If there is no Client
// in the context, ClientFromContext checks if there is a RoundTripper in the context.
// If there is, ClientFromContext returns a new HTTP Client with that RoundTripper.
// Otherwise the default Client is returned.
func ClientFromContext(ctx context.Context) *http.Client {
	if client, ok := ctx.Value(clientCtxKey).(*http.Client); ok {
		return client
	}
	if roundTripper, ok := ctx.Value(roundTripperCtxKey).(http.RoundTripper); ok {
		return &http.Client{Transport: roundTripper}
	}
	return http.DefaultClient
}
