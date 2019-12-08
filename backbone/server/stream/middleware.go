package stream

// Middleware for the HTTP Server.
type Middleware func(next Handler) Handler

// MiddlewareFunc returns middlewareFunc as Middleware.
func MiddlewareFunc(middlewareFunc func(next HandlerFunc) HandlerFunc) Middleware {
	return func(next Handler) Handler { return middlewareFunc(next.HandleStream) }
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
}

func chain(next Handler, m ...Middleware) Handler {
	if len(m) < 1 {
		return next
	}
	h := next
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h
}
