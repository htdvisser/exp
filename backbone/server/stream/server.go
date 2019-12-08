// Package stream provides the backbone for a TCP stream server with some opinionated defaults.
package stream

import (
	"context"
	"net"
	"sync"
	"time"
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

// Server implements a TCP stream server.
type Server struct {
	baseContext                 context.Context
	extendContextWithListener   func(context.Context, net.Listener) context.Context
	extendContextWithConn       func(context.Context, net.Conn) context.Context
	extendContextWithRemoteAddr func(context.Context, net.Addr) context.Context
	middleware                  []Middleware
	handler                     Handler

	mu        sync.Mutex
	listeners map[net.Listener]struct{}
}

// NewServer instantiates a new stream server with the given options.
func NewServer(handler Handler, opts ...Option) *Server {
	options := &options{}
	options.apply(opts...)
	return &Server{
		middleware: options.middleware,
		handler:    handler,
		listeners:  make(map[net.Listener]struct{}),
	}
}

const (
	initialBackoff = time.Millisecond
	maxBackoff     = time.Second
)

// Serve serves the stream server on lis.
func (s *Server) Serve(lis net.Listener) error {
	s.mu.Lock()
	s.listeners[lis] = struct{}{}
	s.mu.Unlock()
	lisCtx := context.Background()
	if s.baseContext != nil {
		lisCtx = s.baseContext
	}
	if s.extendContextWithListener != nil {
		lisCtx = s.extendContextWithListener(lisCtx, lis)
		if lisCtx == nil {
			panic("extendContextWithListener returned a nil context")
		}
	}
	handler := chain(s.handler, s.middleware...)
	var backoff time.Duration
	for {
		connCtx := lisCtx
		conn, err := lis.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if backoff == 0 {
					backoff = initialBackoff
				} else {
					backoff *= 2
				}
				if backoff > maxBackoff {
					backoff = maxBackoff
				}
				time.Sleep(backoff)
				continue
			}
			return err
		}
		backoff = 0
		if s.extendContextWithConn != nil {
			connCtx = s.extendContextWithConn(connCtx, conn)
			if connCtx == nil {
				panic("extendContextWithConn returned a nil context")
			}
		}
		if s.extendContextWithRemoteAddr != nil {
			connCtx = s.extendContextWithRemoteAddr(connCtx, conn.RemoteAddr())
			if connCtx == nil {
				panic("extendContextWithRemoteAddr returned a nil context")
			}
		}
		go func(ctx context.Context) {
			defer conn.Close()
			handler.HandleStream(ctx, conn)
		}(connCtx)
	}
}

// GracefulStop stops the stream server gracefully.
func (s *Server) GracefulStop() error {
	s.mu.Lock()
	for lis := range s.listeners {
		lis.Close()
	}
	s.mu.Unlock()
	return nil
}
