// Package grpc provides the backbone for a gRPC server with some opinionated defaults.
package grpc

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/stats"
)

// Server wraps the gRPC server, gRPC-gateway and a loopback connection.
type Server struct {
	Server  *grpc.Server
	Web     *grpcweb.WrappedGrpcServer
	Gateway *runtime.ServeMux
	Health  *health.Server

	loopbackListener *inProcessListener
	loopbackServing  bool
	loopbackConn     *grpc.ClientConn

	contextExtenders []func(context.Context) context.Context

	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor

	statsHandlers []stats.Handler
}

// NewServer instantiates a new gRPC server with the given options.
func NewServer(opts ...Option) *Server {
	options := &options{
		runtimeServeMuxOptions: []runtime.ServeMuxOption{
			runtime.WithProtoErrorHandler(handleProtoError),
			runtime.WithStreamErrorHandler(handleStreamError),
		},
		runtimeIncomingHeaders: make(runtimeHeaders),
		runtimeOutgoingHeaders: make(runtimeHeaders),
	}
	options.runtimeIncomingHeaders.add(defaultRequestHeaders...)
	options.apply(opts...)
	s := &Server{
		Health: health.NewServer(),

		loopbackListener: newInProcessListener(context.Background()),

		contextExtenders: options.contextExtenders,

		unaryInterceptors:  options.gRPCUnaryInterceptors,
		streamInterceptors: options.gRPCStreamInterceptors,
		statsHandlers:      options.gRPCStatsHandlers,
	}
	gRPCServerOptions := append(
		options.gRPCServerOptions,
		grpc.UnaryInterceptor(s.interceptUnary),
		grpc.StreamInterceptor(s.interceptStream),
		grpc.StatsHandler(&statsHandler{s}),
	)
	grpcWebOptions := append(
		options.grpcWebOptions,
	)
	runtimeServeMuxOptions := append(
		options.runtimeServeMuxOptions,
		runtime.WithIncomingHeaderMatcher(options.runtimeIncomingHeaders.match),
		runtime.WithOutgoingHeaderMatcher(options.runtimeOutgoingHeaders.match),
	)
	s.Server = grpc.NewServer(gRPCServerOptions...)
	s.Web = grpcweb.WrapServer(s.Server, grpcWebOptions...)
	s.Gateway = runtime.NewServeMux(runtimeServeMuxOptions...)
	healthpb.RegisterHealthServer(s.Server, s.Health)
	return s
}

// LoopbackConn returns an in-process gRPC connection to the server.
func (s *Server) LoopbackConn() *grpc.ClientConn {
	if s.loopbackConn == nil {
		s.loopbackConn, _ = grpc.Dial(
			s.loopbackListener.Addr().String(),
			grpc.WithDialer(inProcessDialer(s.loopbackListener)),
			grpc.WithTransportCredentials(&inProcessCredentials{}),
		)
	}
	return s.loopbackConn
}

// ServeLoopback serves the gRPC server for the in-process gRPC connection.
func (s *Server) ServeLoopback() error {
	if !s.loopbackServing {
		s.loopbackServing = true
		return s.Server.Serve(s.loopbackListener)
	}
	return nil
}

// Serve serves the gRPC server (but not the gRPC-gateway) on lis.
func (s *Server) Serve(lis net.Listener) error {
	return s.Server.Serve(lis)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	switch {
	case s.Web.IsGrpcWebRequest(r) || s.Web.IsAcceptableGrpcCorsRequest(r):
		s.Web.ServeHTTP(w, r)
	case strings.HasPrefix(contentType, "application/grpc"):
		s.Server.ServeHTTP(w, r)
	default:
		http.NotFound(w, r)
	}
}

// GracefulStop stops the gRPC server gracefully.
func (s *Server) GracefulStop() error {
	s.Health.Shutdown()
	s.Server.GracefulStop()
	return nil
}
