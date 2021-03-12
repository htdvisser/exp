// Package recovery can be used to add panic recovery middleware to the server.
package recovery

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"htdvisser.dev/exp/backbone/server"
	"htdvisser.dev/exp/backbone/server/packet"
	"htdvisser.dev/exp/backbone/server/stream"
)

// Middleware is middleware for panic recovery.
type Middleware struct {
	panicToError        func(ctx context.Context, p interface{}) error
	errorToHTTPResponse func(w http.ResponseWriter, r *http.Request, err error)
}

// Option is an option for the panic recovery middleware.
type Option interface {
	apply(*Middleware)
}

type option func(*Middleware)

func (f option) apply(opts *Middleware) {
	f(opts)
}

// WithPanicToError returns an option that sets the function to convert panics to errors.
func WithPanicToError(f func(ctx context.Context, p interface{}) error) Option {
	return option(func(opts *Middleware) {
		opts.panicToError = f
	})
}

// WithErrorToHTTPResponse returns an option that sets the function to write errors to HTTP responses.
func WithErrorToHTTPResponse(f func(w http.ResponseWriter, r *http.Request, err error)) Option {
	return option(func(opts *Middleware) {
		opts.errorToHTTPResponse = f
	})
}

// NewMiddleware returns new middleware for panic recovery.
func NewMiddleware(opts ...Option) (*Middleware, error) {
	m := &Middleware{
		panicToError: func(_ context.Context, p interface{}) error {
			if err, ok := p.(error); ok {
				return err
			}
			return status.Errorf(codes.Internal, "%s", p)
		},
		errorToHTTPResponse: func(w http.ResponseWriter, _ *http.Request, err error) {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		},
	}
	for _, opt := range opts {
		opt.apply(m)
	}
	return m, nil
}

// RecoverUnaryRPC recovers from panics in unary RPCs.
func (m *Middleware) RecoverUnaryRPC(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if p := recover(); p != nil {
			err = m.panicToError(ctx, p)
		}
	}()
	return handler(ctx, req)
}

// RecoverStreamingRPC recovers from panics in streaming RPCs.
func (m *Middleware) RecoverStreamingRPC(srv interface{}, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	defer func() {
		if p := recover(); p != nil {
			ctx := ss.Context()
			err = m.panicToError(ctx, p)
		}
	}()
	return handler(srv, ss)
}

// RecoverHTTP recovers from panics in HTTP handlers.
func (m *Middleware) RecoverHTTP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if p := recover(); p != nil {
				ctx := r.Context()
				err := m.panicToError(ctx, p)
				m.errorToHTTPResponse(w, r, err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// RecoverStream recovers from panics in stream handlers.
func (m *Middleware) RecoverStream(next stream.HandlerFunc) stream.HandlerFunc {
	return stream.HandlerFunc(func(ctx context.Context, conn net.Conn) (err error) {
		defer func() {
			if p := recover(); p != nil {
				err = m.panicToError(ctx, p)
			}
		}()
		return next.HandleStream(ctx, conn)
	})
}

// RecoverPacket recovers from panics in packet handlers.
func (m *Middleware) RecoverPacket(next packet.HandlerFunc) packet.HandlerFunc {
	return packet.HandlerFunc(func(ctx context.Context, pkt []byte, addr net.Addr, reply func([]byte) error) (err error) {
		defer func() {
			if p := recover(); p != nil {
				err = m.panicToError(ctx, p)
			}
		}()
		return next.HandlePacket(ctx, pkt, addr, reply)
	})
}

// Register registers the panic recovery to the server.
func (m *Middleware) Register(s *server.Server) error {
	s.GRPC.AddUnaryInterceptor(m.RecoverUnaryRPC)
	s.GRPC.AddStreamInterceptor(m.RecoverStreamingRPC)
	s.HTTP.AddMiddleware(m.RecoverHTTP)
	s.InternalGRPC.AddUnaryInterceptor(m.RecoverUnaryRPC)
	s.InternalGRPC.AddStreamInterceptor(m.RecoverStreamingRPC)
	s.InternalHTTP.AddMiddleware(m.RecoverHTTP)
	return nil
}

// Register registers the panic recovery to the server.
func Register(s *server.Server, opts ...Option) error {
	m, err := NewMiddleware(opts...)
	if err != nil {
		return err
	}
	return m.Register(s)
}
