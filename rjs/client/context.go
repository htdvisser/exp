package client

import "context"

type contextKeyType struct{}

var contextKey contextKeyType

// NewContextWithClient returns a context derived from parent that contains the Client.
func NewContextWithClient(parent context.Context, client *Client) context.Context {
	return context.WithValue(parent, contextKey, client)
}

// FromContext returns the Client if the context contained one.
func FromContext(ctx context.Context) *Client {
	if client, ok := ctx.Value(contextKey).(*Client); ok {
		return client
	}
	return nil
}
