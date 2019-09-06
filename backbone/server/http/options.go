package http

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

type options struct {
	serveMux         *http.ServeMux
	router           *mux.Router
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

// WithRouter returns an option that sets the Router of the server.
func WithRouter(router *mux.Router) Option {
	return option(func(opts *options) {
		opts.router = router
	})
}
