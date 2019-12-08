package packet

import (
	"context"
	"net"
)

// Handler is the interface for handling packets.
type Handler interface {
	HandlePacket(context.Context, []byte, net.Addr, func([]byte) error) error
}

// HandlerFunc is the Handler func.
type HandlerFunc func(context.Context, []byte, net.Addr, func([]byte) error) error

// HandlePacket implements the Handler interface.
func (f HandlerFunc) HandlePacket(ctx context.Context, pkt []byte, addr net.Addr, reply func([]byte) error) error {
	return f(ctx, pkt, addr, reply)
}
