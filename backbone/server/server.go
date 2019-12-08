package server

import (
	"context"
	"fmt"
	"log"
	"net"
	stdhttp "net/http"
	_ "net/http/pprof" // Registers pprof endpoints to DefaultServeMux (the internal HTTP server).

	"go.opentelemetry.io/otel/api/trace"
	"golang.org/x/sync/errgroup"
	"htdvisser.dev/exp/backbone/server/grpc"
	"htdvisser.dev/exp/backbone/server/http"
	"htdvisser.dev/exp/backbone/server/internal/channelz"
)

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

	TraceProvider trace.Provider

	GRPC *grpc.Server
	HTTP *http.Server

	InternalGRPC *grpc.Server
	InternalHTTP *http.Server

	tcpServers []tcpServer

	runGroup   *errgroup.Group
	runContext context.Context
}

// New instantiates a new server that uses the config and options.
func New(config Config, opts ...Option) *Server {
	options := &options{
		InternalHTTPOptions: []http.Option{
			http.WithServeMux(stdhttp.DefaultServeMux),
		},
	}
	options.apply(opts...)
	s := &Server{
		config:        config,
		TraceProvider: &trace.NoopProvider{},
		GRPC:          grpc.NewServer(options.GRPCOptions...),
		HTTP:          http.NewServer(options.HTTPOptions...),
		InternalGRPC:  grpc.NewServer(options.InternalGRPCOptions...),
		InternalHTTP:  http.NewServer(options.InternalHTTPOptions...),
	}
	channelz.Register(s.InternalGRPC)
	s.RegisterTCPServer("gRPC", s.config.ListenGRPC, s.GRPC)
	s.RegisterTCPServer("internal gRPC", s.config.ListenInternalGRPC, s.InternalGRPC)
	s.RegisterTCPServer("HTTP", s.config.ListenHTTP, s.HTTP)
	s.RegisterTCPServer("internal HTTP", s.config.ListenInternalHTTP, s.InternalHTTP)
	return s
}

// RegisterTCPServer registers the named TCP server on address.
func (s *Server) RegisterTCPServer(name, address string, server interface {
	Serve(lis net.Listener) error
	GracefulStop() error
}) error {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return err
	}
	address = addr.String()
	for _, registered := range s.tcpServers {
		if address == registered.address {
			return fmt.Errorf("could not register %q server: %w",
				name, fmt.Errorf("%q already registered on %q",
					registered.name, address))
		}
	}
	s.tcpServers = append(s.tcpServers, tcpServer{name: name, address: address, server: server})
	return nil
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
	s.runGroup.Go(s.GRPC.ServeLoopback)
	s.runGroup.Go(s.InternalGRPC.ServeLoopback)
	if err = s.runTCPServers(ctx); err != nil {
		return err
	}
	<-s.runContext.Done()
	return s.runContext.Err()
}

func (s *Server) runTCPServers(ctx context.Context) error {
	for _, tcpServer := range s.tcpServers {
		if err := s.runTCPServer(ctx, tcpServer.name, tcpServer.address, tcpServer.server); err != nil {
			return err
		}
	}
	return nil
}

type tcpServer struct {
	name    string
	address string
	server  interface {
		Serve(lis net.Listener) error
		GracefulStop() error
	}
}

func (s *Server) runTCPServer(ctx context.Context, name, address string, server interface {
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
