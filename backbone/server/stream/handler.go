package stream

import (
	"context"
	"net"
)

// Handler is the interface for handling streams.
type Handler interface {
	HandleStream(context.Context, net.Conn) error
}

// HandlerFunc is the Handler func.
type HandlerFunc func(context.Context, net.Conn) error

// HandleStream implements the Handler interface.
func (f HandlerFunc) HandleStream(ctx context.Context, conn net.Conn) error {
	return f(ctx, conn)
}
