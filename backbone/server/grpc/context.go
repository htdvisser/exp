package grpc

import "context"

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

func (s *Server) extendContext(ctx context.Context) context.Context {
	for _, extendContext := range s.contextExtenders {
		ctx = extendContext(ctx)
	}
	return ctx
}
