package http

import (
	"context"
	"net/http"
)

// WithContextExtender adds functions that extend request contexts.
func WithContextExtender(contextExtender ...func(context.Context) context.Context) Option {
	return option(func(opts *options) {
		opts.contextExtenders = append(opts.contextExtenders, contextExtender...)
	})
}

// AddContextExtender adds functions that extend request contexts.
func (s *Server) AddContextExtender(contextExtender ...func(context.Context) context.Context) {
	s.contextExtenders = append(s.contextExtenders, contextExtender...)
}

func (s *Server) extendContext(r *http.Request) *http.Request {
	ctx := r.Context()
	for _, extendContext := range s.contextExtenders {
		ctx = extendContext(ctx)
	}
	return r.WithContext(ctx)
}
