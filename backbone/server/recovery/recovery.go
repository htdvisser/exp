// Package recovery can be used to add panic recovery middleware to the server.
package recovery

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"htdvisser.dev/exp/backbone/server"
)

type options struct {
	panicToError        func(ctx context.Context, p interface{}) error
	errorToHTTPResponse func(w http.ResponseWriter, r *http.Request, err error)
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt.apply(o)
	}
}

// Option is an option for the panic recovery.
type Option interface {
	apply(*options)
}

type option func(*options)

func (f option) apply(opts *options) {
	f(opts)
}

// WithPanicToError returns an option that sets the function to convert panics to errors.
func WithPanicToError(f func(ctx context.Context, p interface{}) error) Option {
	return option(func(opts *options) {
		opts.panicToError = f
	})
}

// WithErrorToHTTPResponse returns an option that sets the function to write errors to HTTP responses.
func WithErrorToHTTPResponse(f func(w http.ResponseWriter, r *http.Request, err error)) Option {
	return option(func(opts *options) {
		opts.errorToHTTPResponse = f
	})
}

// Register registers the panic recovery to the server.
func Register(s *server.Server, opts ...Option) error {
	options := &options{
		panicToError: func(_ context.Context, p interface{}) error {
			if err, ok := p.(error); ok {
				return err
			}
			return status.Errorf(codes.Internal, "%s", p)
		},
		errorToHTTPResponse: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		},
	}
	options.apply(opts...)
	s.GRPC.AddUnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if p := recover(); p != nil {
				err = options.panicToError(ctx, p)
			}
		}()
		return handler(ctx, req)
	})
	s.GRPC.AddStreamInterceptor(func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if p := recover(); p != nil {
				err = options.panicToError(ss.Context(), p)
			}
		}()
		return handler(srv, ss)
	})
	s.HTTP.AddMiddleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if p := recover(); p != nil {
					err := options.panicToError(r.Context(), p)
					options.errorToHTTPResponse(w, r, err)
				}
			}()
			next.ServeHTTP(w, r)
		})
	})
	return nil
}
