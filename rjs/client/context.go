package client

import "context"

type contextKeyType struct{}

var contextKey contextKeyType

func NewContextWithClient(parent context.Context, client *Client) context.Context {
	return context.WithValue(parent, contextKey, client)
}

func FromContext(ctx context.Context) *Client {
	if client, ok := ctx.Value(contextKey).(*Client); ok {
		return client
	}
	return nil
}
