package server

import (
	"context"
	"log"
	"net"
	stdhttp "net/http"
	_ "net/http/pprof" // Registers pprof endpoints to DefaultServeMux (the internal HTTP server).

	"golang.org/x/sync/errgroup"
	"htdvisser.dev/exp/backbone/server/grpc"
	"htdvisser.dev/exp/backbone/server/http"
	"htdvisser.dev/exp/backbone/server/internal/channelz"
)

type options struct {
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

// Config is the required configuration for the server.
type Config struct {
	ListenHTTP         string
	ListenGRPC         string
	ListenInternalHTTP string
	ListenInternalGRPC string
}

// Server wraps gRPC and HTTP servers.
type Server struct {
	config Config

	GRPC *grpc.Server
	HTTP *http.Server

	InternalGRPC *grpc.Server
	InternalHTTP *http.Server

	runGroup   *errgroup.Group
	runContext context.Context
}

// New instantiates a new server that uses the config and options.
func New(config Config, opts ...Option) *Server {
	options := &options{}
	options.apply(opts...)
	s := &Server{
		config:       config,
		GRPC:         grpc.NewServer(),
		HTTP:         http.NewServer(),
		InternalGRPC: grpc.NewServer(),
		InternalHTTP: http.NewServer(http.WithServeMux(stdhttp.DefaultServeMux)),
	}
	channelz.Register(s.InternalGRPC)
	return s
}

// Run runs the server until the Done channel of ctx is closed.
func (s *Server) Run(ctx context.Context) (err error) {
	if s.runGroup != nil {
		panic("server is already running")
	}
	s.runGroup, s.runContext = errgroup.WithContext(ctx)
	defer func() {
		gErr := s.runGroup.Wait()
		if ctx.Err() == nil {
			err = gErr
		}
	}()

	s.runServer(ctx, "gRPC", s.config.ListenGRPC, s.GRPC)
	s.runGroup.Go(s.GRPC.ServeLoopback)
	s.runServer(ctx, "internal gRPC", s.config.ListenInternalGRPC, s.InternalGRPC)
	s.runGroup.Go(s.InternalGRPC.ServeLoopback)
	s.runServer(ctx, "HTTP", s.config.ListenHTTP, s.HTTP)
	s.runServer(ctx, "internal HTTP", s.config.ListenInternalHTTP, s.InternalHTTP)

	<-s.runContext.Done()
	return s.runContext.Err()
}

func (s *Server) runServer(ctx context.Context, name, address string, server interface {
	Serve(lis net.Listener) error
	GracefulStop() error
}) error {
	go func() {
		<-ctx.Done()
		s.runGroup.Go(func() error {
			log.Printf("Gracefully stopping %s server...", name)
			server.GracefulStop()
			return nil
		})
	}()
	if address != "" {
		lis, err := net.Listen("tcp", address)
		if err != nil {
			return err
		}
		log.Printf("Serving %s on %s...", name, lis.Addr().String())
		s.runGroup.Go(func() error {
			return server.Serve(lis)
		})
	}
	return nil
}