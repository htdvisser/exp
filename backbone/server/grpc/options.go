package grpc

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
)

type options struct {
	contextExtenders       []func(context.Context) context.Context
	gRPCUnaryInterceptors  []grpc.UnaryServerInterceptor
	gRPCStreamInterceptors []grpc.StreamServerInterceptor
	gRPCServerOptions      []grpc.ServerOption
	gRPCStatsHandlers      []stats.Handler
	runtimeServeMuxOptions []runtime.ServeMuxOption
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt.apply(o)
	}
}

// Option is an option for the gRPC server.
type Option interface {
	apply(*options)
}

type option func(*options)

func (f option) apply(opts *options) {
	f(opts)
}

// WithGRPCServerOption adds serverOptions. The options UnaryInterceptor,
// StreamInterceptor and StatsHandler should not be used.
func WithGRPCServerOption(serverOptions ...grpc.ServerOption) Option {
	return option(func(o *options) {
		o.gRPCServerOptions = append(o.gRPCServerOptions, serverOptions...)
	})
}

// WithRuntimeServeMuxOption adds serveMuxOptions.
func WithRuntimeServeMuxOption(serveMuxOptions ...runtime.ServeMuxOption) Option {
	return option(func(o *options) {
		o.runtimeServeMuxOptions = append(o.runtimeServeMuxOptions, serveMuxOptions...)
	})
}
