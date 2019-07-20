package http

import (
	"context"
	"net/http"
)

type options struct {
	serveMux         *http.ServeMux
	contextExtenders []func(context.Context) context.Context
	middleware       []Middleware
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt.apply(o)
	}
}

// Option is an option for the HTTP server.
type Option interface {
	apply(*options)
}

type option func(*options)

func (f option) apply(opts *options) {
	f(opts)
}

// WithServeMux returns an option that sets the ServeMux of the server.
func WithServeMux(serveMux *http.ServeMux) Option {
	return option(func(opts *options) {
		opts.serveMux = serveMux
	})
}
