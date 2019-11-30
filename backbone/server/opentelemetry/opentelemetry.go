// Package opentelemetry can be used to add opentelemetry support to the server.
package opentelemetry

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/distributedcontext"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/propagation"
	"go.opentelemetry.io/otel/api/trace"
	sdkpropagation "go.opentelemetry.io/otel/propagation"
	exporttrace "go.opentelemetry.io/otel/sdk/export/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	grpcclient "htdvisser.dev/exp/backbone/client/grpc"
	"htdvisser.dev/exp/backbone/server"
)

type options struct {
	propagator        propagation.TextFormatPropagator
	sampleProbability float64
	syncers           []exporttrace.SpanSyncer
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt.apply(o)
	}
}

// Option is an option for the opentelemetry metrics and tracing.
type Option interface {
	apply(*options)
}

type option func(*options)

func (f option) apply(opts *options) {
	f(opts)
}

// WithSampleProbability sets the default trace sampling probability.
func WithSampleProbability(probability float64) Option {
	return option(func(opts *options) {
		opts.sampleProbability = probability
	})
}

// WithSyncer adds a span syncer to the tracer.
func WithSyncer(syncer exporttrace.SpanSyncer) Option {
	return option(func(opts *options) {
		opts.syncers = append(opts.syncers, syncer)
	})
}

// Register adds opentelemetry metrics and tracing to the server.
func Register(s *server.Server, opts ...Option) error {
	options := &options{
		propagator: &sdkpropagation.B3Propagator{},
	}
	options.apply(opts...)
	traceOpts := []sdktrace.ProviderOption{
		sdktrace.WithConfig(sdktrace.Config{
			DefaultSampler: sdktrace.ProbabilitySampler(options.sampleProbability),
		}),
	}
	for _, syncer := range options.syncers {
		traceOpts = append(traceOpts, sdktrace.WithSyncer(syncer))
	}
	tp, err := sdktrace.NewProvider(traceOpts...)
	if err != nil {
		return err
	}
	s.TraceProvider = tp
	contextExtender := func(ctx context.Context) context.Context {
		return grpcclient.NewContextWithDialOptions(ctx,
			grpc.WithChainUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
				md, _ := metadata.FromOutgoingContext(ctx)
				tr := s.TraceProvider.Tracer("backbone/grpc")
				ctx, span := tr.Start(ctx, method)
				defer span.End()
				md = md.Copy()
				options.propagator.Inject(ctx, mdAttrs(md))
				ctx = metadata.NewOutgoingContext(ctx, md)
				err := invoker(ctx, method, req, reply, cc, opts...)
				if status, ok := status.FromError(err); ok {
					trace.CurrentSpan(ctx).SetStatus(status.Code())
				}
				return err
			}),
			// TODO: Streaming Client Interceptor.
			// grpc.WithChainStreamInterceptor(func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (ClientStream, error) {
			// 	...
			// }),
		)
	}
	s.GRPC.AddContextExtender(contextExtender)
	s.HTTP.AddContextExtender(contextExtender)
	s.GRPC.AddUnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, _ := metadata.FromIncomingContext(ctx)
		spanCtx, correlationCtx := options.propagator.Extract(ctx, mdAttrs(md))
		correlationCtxKVs := extractKVs(correlationCtx)
		ctx = distributedcontext.WithMap(ctx, distributedcontext.NewMap(distributedcontext.MapUpdate{
			MultiKV: correlationCtxKVs,
		}))
		attrs := grpcAttrsFromContext(ctx)
		tr := s.TraceProvider.Tracer("backbone/grpc")
		ctx, span := tr.Start(
			ctx,
			strings.TrimPrefix(info.FullMethod, "/"),
			trace.WithAttributes(attrs...),
			trace.ChildOf(spanCtx),
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()
		resp, err = handler(ctx, req)
		if status, ok := status.FromError(err); ok {
			trace.CurrentSpan(ctx).SetStatus(status.Code())
		}
		return resp, err
	})
	s.GRPC.AddStreamInterceptor(func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		md, _ := metadata.FromIncomingContext(ctx)
		spanCtx, correlationCtx := options.propagator.Extract(ctx, mdAttrs(md))
		correlationCtxKVs := extractKVs(correlationCtx)
		ctx = distributedcontext.WithMap(ctx, distributedcontext.NewMap(distributedcontext.MapUpdate{
			MultiKV: correlationCtxKVs,
		}))
		attrs := grpcAttrsFromContext(ctx)
		tr := s.TraceProvider.Tracer("backbone/grpc")
		ctx, span := tr.Start(
			ctx,
			strings.TrimPrefix(info.FullMethod, "/"),
			trace.WithAttributes(attrs...),
			trace.ChildOf(spanCtx),
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()
		wrappedStream := middleware.WrapServerStream(ss)
		wrappedStream.WrappedContext = ctx
		// TODO: Intercept stream messages.
		// TODO: Stream status.
		return handler(srv, wrappedStream)
	})
	s.HTTP.AddMiddleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			spanCtx, correlationCtx := options.propagator.Extract(ctx, r.Header)
			correlationCtxKVs := extractKVs(correlationCtx)
			ctx = distributedcontext.WithMap(ctx, distributedcontext.NewMap(distributedcontext.MapUpdate{
				MultiKV: correlationCtxKVs,
			}))
			attrs := httpAttrsFromRequest(r)
			tr := s.TraceProvider.Tracer("backbone/http")
			spanName := strings.TrimPrefix(r.URL.Path, "/")
			if route := mux.CurrentRoute(r); route != nil {
				if pathTemplate, err := route.GetPathTemplate(); err == nil {
					spanName = strings.TrimPrefix(pathTemplate, "/")
				}
			}
			ctx, span := tr.Start(
				ctx,
				spanName,
				trace.WithAttributes(attrs...),
				trace.ChildOf(spanCtx),
				trace.WithSpanKind(trace.SpanKindServer),
			)
			defer span.End()
			next.ServeHTTP(w, r.WithContext(ctx))
			// TODO: Get HTTP Status.
		})
	})
	return nil
}

type mdAttrs metadata.MD

func (a mdAttrs) Get(key string) string {
	return strings.Join(metadata.MD(a).Get(key), ",")
}

func (a mdAttrs) Set(key string, value string) {
	metadata.MD(a).Set(key, value)
}

func extractKVs(correlationCtx distributedcontext.Map) []core.KeyValue {
	var kvs []core.KeyValue
	correlationCtx.Foreach(func(kv core.KeyValue) bool {
		kvs = append(kvs, kv)
		return true
	})
	return kvs
}

var (
	componentKey  = key.New("component")
	httpMethodKey = key.New("http.method")
	httpFlavorKey = key.New("http.flavor")
	httpURLKey    = key.New("http.url")
	httpTargetKey = key.New("http.target")
	httpHostKey   = key.New("http.host")
	hostPortKey   = key.New("host.port")
	httpSchemeKey = key.New("http.scheme")
	httpRouteKey  = key.New("http.route")
	peerIP4Key    = key.New("peer.ip4")
	peerIP6Key    = key.New("peer.ip6")
)

func grpcAttrsFromContext(ctx context.Context) []core.KeyValue {
	kvs := []core.KeyValue{
		componentKey.String("grpc"),
	}
	if peer, ok := peer.FromContext(ctx); ok {
		if host, _, err := net.SplitHostPort(peer.Addr.String()); err == nil {
			if ip := net.ParseIP(host); ip != nil {
				if ipv4 := ip.To4(); len(ipv4) == net.IPv4len {
					kvs = append(kvs, peerIP4Key.String(ip.String()))
				} else {
					kvs = append(kvs, peerIP6Key.String(ip.String()))
				}
			}
		}
	}
	return kvs
}

func httpAttrsFromRequest(r *http.Request) []core.KeyValue {
	kvs := []core.KeyValue{
		componentKey.String("http"),
		httpMethodKey.String(r.Method),
		httpFlavorKey.String(strings.TrimPrefix(r.Proto, "HTTP/")),
		httpTargetKey.String(r.URL.String()),
		httpHostKey.String(r.Host),
		httpSchemeKey.String(r.URL.Scheme),
	}
	if port, err := strconv.Atoi(r.URL.Port()); err == nil {
		kvs = append(kvs, hostPortKey.Int(port))
	}
	if route := mux.CurrentRoute(r); route != nil {
		if pathTemplate, err := route.GetPathTemplate(); err == nil {
			kvs = append(kvs, httpRouteKey.String(pathTemplate))
		}
	}
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		if ip := net.ParseIP(host); ip != nil {
			if ipv4 := ip.To4(); len(ipv4) == net.IPv4len {
				kvs = append(kvs, peerIP4Key.String(ip.String()))
			} else {
				kvs = append(kvs, peerIP6Key.String(ip.String()))
			}
		}
	}
	return kvs
}
