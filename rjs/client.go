package rjs

import "context"

type Client interface {
	Call(ctx context.Context, uri string, request interface{}, response interface{}) error
}
