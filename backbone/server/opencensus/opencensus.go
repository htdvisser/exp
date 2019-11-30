// Package opencensus can be used to add opencensus support to the server.
package opencensus

import (
	"context"
	"fmt"
	"net/http"

	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/zpages"
	"google.golang.org/grpc"
	grpcclient "htdvisser.dev/exp/backbone/client/grpc"
	httpclient "htdvisser.dev/exp/backbone/client/http"
	_ "htdvisser.dev/exp/backbone/client/opencensus" // Import to register views.
	"htdvisser.dev/exp/backbone/server"
)

func init() {
	if err := view.Register(ocgrpc.DefaultServerViews...); err != nil {
		panic(fmt.Errorf("Failed to register server views for gRPC metrics: %v", err))
	}
	if err := view.Register(ochttp.DefaultServerViews...); err != nil {
		panic(fmt.Errorf("Failed to register server views for HTTP metrics: %v", err))
	}
}

type options struct {
	isPublicEndpoint bool
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt.apply(o)
	}
}

// Option is an option for the opencensus metrics and tracing.
type Option interface {
	apply(*options)
}

type option func(*options)

func (f option) apply(opts *options) {
	f(opts)
}

// WithIsPublicEndpoint returns an option that configures opencensus to treat
// incoming requests as external.
func WithIsPublicEndpoint(isPublicEndpoint bool) Option {
	return option(func(opts *options) {
		opts.isPublicEndpoint = isPublicEndpoint
	})
}

func contextExtender(ctx context.Context) context.Context {
	ctx = grpcclient.NewContextWithDialOptions(ctx, grpc.WithStatsHandler(&ocgrpc.ClientHandler{}))
	ctx = httpclient.NewContextWithRoundTripper(ctx, &ochttp.Transport{
		Base: httpclient.RoundTripperFromContext(ctx),
	})
	return ctx
}

// Register adds opencensus metrics and tracing to the server.
func Register(s *server.Server, opts ...Option) error {
	options := &options{
		isPublicEndpoint: true,
	}
	options.apply(opts...)
	s.GRPC.AddStatsHandler(&ocgrpc.ServerHandler{
		IsPublicEndpoint: options.isPublicEndpoint,
	})
	s.GRPC.AddContextExtender(contextExtender)
	s.HTTP.AddContextExtender(contextExtender)
	s.HTTP.AddMiddleware(func(next http.Handler) http.Handler {
		return &ochttp.Handler{
			Handler:          next,
			IsPublicEndpoint: options.isPublicEndpoint,
		}
	})
	zpages.Handle(s.InternalHTTP.ServeMux, "/debug")
	return nil
}
