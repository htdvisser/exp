// Package packet provides the backbone for a UDP packet server with some opinionated defaults.
package packet

import (
	"context"
	"net"
	"sync"
)

// Server implements a UDP packet server.
type Server struct {
	baseContext                 context.Context
	extendContextWithPacketConn func(context.Context, net.PacketConn) context.Context
	extendContextWithRemoteAddr func(context.Context, net.Addr) context.Context

	middleware []Middleware
	handler    Handler

	mu       sync.Mutex
	bindings map[net.PacketConn]struct{}
}

// NewServer instantiates a new packet server with the given options.
func NewServer(handler Handler, opts ...Option) *Server {
	options := &options{}
	options.apply(opts...)
	return &Server{
		middleware: options.middleware,
		handler:    handler,
		bindings:   make(map[net.PacketConn]struct{}),
	}
}

func (s *Server) Serve(conn net.PacketConn) error {
	s.mu.Lock()
	s.bindings[conn] = struct{}{}
	s.mu.Unlock()
	connCtx := context.Background()
	if s.baseContext != nil {
		connCtx = s.baseContext
	}
	if s.extendContextWithPacketConn != nil {
		connCtx = s.extendContextWithPacketConn(connCtx, conn)
		if connCtx == nil {
			panic("extendContextWithPacketConn returned a nil context")
		}
	}
	handler := chain(s.handler, s.middleware...)
	var buf [0xffff]byte
	for {
		n, addr, err := conn.ReadFrom(buf[:])
		if n > 0 {
			pkt := make([]byte, n)
			copy(pkt, buf[:n])
			pktCtx := connCtx
			if s.extendContextWithRemoteAddr != nil {
				pktCtx = s.extendContextWithRemoteAddr(pktCtx, addr)
				if pktCtx == nil {
					panic("extendContextWithRemoteAddr returned a nil context")
				}
			}
			handler.HandlePacket(pktCtx, pkt, addr, func(res []byte) error {
				_, err := conn.WriteTo(res, addr)
				return err
			})
		}
		if err != nil {
			return err
		}
	}
}

// GracefulStop stops the stream server gracefully.
func (s *Server) GracefulStop() error {
	s.mu.Lock()
	for conn := range s.bindings {
		conn.Close()
	}
	s.mu.Unlock()
	return nil
}
