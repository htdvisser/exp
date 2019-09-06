// Package http provides the backbone for an HTTP server with some opinionated defaults.
package http

import (
	"context"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

// Server wraps the HTTP server.
type Server struct {
	ServeMux         *http.ServeMux
	Router           *mux.Router
	server           *http.Server
	contextExtenders []func(context.Context) context.Context
	middleware       []Middleware
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
	s.ServeMux.Handle("/", s.Router)
	s.server = &http.Server{Handler: s}
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	middleware := chain(s.ServeMux, s.middleware...)
	middleware.ServeHTTP(w, s.extendContext(r))
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
