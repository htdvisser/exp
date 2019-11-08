package server

import (
	"htdvisser.dev/exp/backbone/server/grpc"
	"htdvisser.dev/exp/backbone/server/http"
)

type options struct {
	HTTPOptions         []http.Option
	GRPCOptions         []grpc.Option
	InternalHTTPOptions []http.Option
	InternalGRPCOptions []grpc.Option
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt.apply(o)
	}
}

// Option is an option for the server.
type Option interface {
	apply(*options)
}

type option func(*options)

func (f option) apply(opts *options) {
	f(opts)
}

// WithHTTPOptions returns an Option that adds the given HTTP options.
func WithHTTPOptions(opts ...http.Option) Option {
	return option(func(o *options) {
		o.HTTPOptions = append(o.HTTPOptions, opts...)
	})
}

// WithGRPCOptions returns an Option that adds the given GRPC options.
func WithGRPCOptions(opts ...grpc.Option) Option {
	return option(func(o *options) {
		o.GRPCOptions = append(o.GRPCOptions, opts...)
	})
}

// WithInternalHTTPOptions returns an Option that adds the given HTTP options for the internal server.
func WithInternalHTTPOptions(opts ...http.Option) Option {
	return option(func(o *options) {
		o.InternalHTTPOptions = append(o.InternalHTTPOptions, opts...)
	})
}

// WithInternalGRPCOptions returns an Option that adds the given GRPC options for the internal server.
func WithInternalGRPCOptions(opts ...grpc.Option) Option {
	return option(func(o *options) {
		o.InternalGRPCOptions = append(o.InternalGRPCOptions, opts...)
	})
}
