// Package http provides the backbone for an HTTP server with some opinionated defaults.
package http

import (
	"context"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// Server wraps the HTTP server.
type Server struct {
	ServeMux         *http.ServeMux
	Router           *mux.Router
	server           *http.Server
	http2server      *http2.Server
	contextExtenders []func(context.Context) context.Context
	middleware       []Middleware
	chain            http.Handler
}

// NewServer instantiates a new HTTP server with the given options.
func NewServer(opts ...Option) *Server {
	options := &options{
		serveMux: http.NewServeMux(),
		router:   mux.NewRouter(),
	}
	options.apply(opts...)
	s := &Server{
		ServeMux:         options.serveMux,
		Router:           options.router,
		contextExtenders: options.contextExtenders,
		middleware:       options.middleware,
	}
	s.chain = chain(s.ServeMux, s.middleware...)
	s.ServeMux.Handle("/", s.Router)
	var handler http.Handler = s
	if options.h2c {
		s.http2server = &http2.Server{}
		if s.http2server != nil {
			handler = h2c.NewHandler(s, s.http2server)
		}
	}
	s.server = &http.Server{Handler: handler}
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.chain.ServeHTTP(w, s.extendContext(r))
}

// Serve serves the HTTP server on lis.
func (s *Server) Serve(lis net.Listener) error {
	return s.server.Serve(lis)
}

// GracefulStop stops the HTTP server gracefully.
func (s *Server) GracefulStop() error {
	s.server.Shutdown(context.Background())
	return nil
}
