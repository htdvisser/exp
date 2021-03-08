package http

import "net/http"

// Middleware for the HTTP Server.
type Middleware func(next http.Handler) http.Handler

// MiddlewareFunc returns middlewareFunc as Middleware.
func MiddlewareFunc(middlewareFunc func(next http.HandlerFunc) http.HandlerFunc) Middleware {
	return func(next http.Handler) http.Handler { return middlewareFunc(next.ServeHTTP) }
}

// WithMiddleware adds middleware.
func WithMiddleware(middleware ...Middleware) Option {
	return option(func(opts *options) {
		opts.middleware = append(opts.middleware, middleware...)
	})
}

// AddMiddleware adds middleware.
func (s *Server) AddMiddleware(middleware ...Middleware) {
	s.middleware = append(s.middleware, middleware...)
	s.chain = chain(s.ServeMux, s.middleware...)
}

func chain(next http.Handler, m ...Middleware) http.Handler {
	if len(m) < 1 {
		return next
	}
	h := next
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h
}
