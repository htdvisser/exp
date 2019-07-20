package grpc

import (
	"context"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

func (s *Server) interceptUnary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	ctx = s.extendContext(ctx)
	return middleware.ChainUnaryServer(s.unaryInterceptors...)(ctx, req, info, handler)
}

// WithUnaryInterceptor adds unary interceptors.
func WithUnaryInterceptor(interceptor ...grpc.UnaryServerInterceptor) Option {
	return option(func(opts *options) {
		opts.gRPCUnaryInterceptors = append(opts.gRPCUnaryInterceptors, interceptor...)
	})
}

// AddUnaryInterceptor adds unary interceptors.
func (s *Server) AddUnaryInterceptor(interceptor ...grpc.UnaryServerInterceptor) {
	s.unaryInterceptors = append(s.unaryInterceptors, interceptor...)
}

func (s *Server) interceptStream(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	wrappedStream := middleware.WrapServerStream(ss)
	wrappedStream.WrappedContext = s.extendContext(ss.Context())
	return middleware.ChainStreamServer(s.streamInterceptors...)(srv, wrappedStream, info, handler)
}

// WithStreamInterceptor adds stream interceptors.
func WithStreamInterceptor(interceptor ...grpc.StreamServerInterceptor) Option {
	return option(func(opts *options) {
		opts.gRPCStreamInterceptors = append(opts.gRPCStreamInterceptors, interceptor...)
	})
}

// AddStreamInterceptor adds stream interceptors.
func (s *Server) AddStreamInterceptor(interceptor ...grpc.StreamServerInterceptor) {
	s.streamInterceptors = append(s.streamInterceptors, interceptor...)
}
