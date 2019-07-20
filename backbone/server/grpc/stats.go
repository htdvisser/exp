package grpc

import (
	"context"

	"google.golang.org/grpc/stats"
)

// WithStatsHandler adds stats handlers.
func WithStatsHandler(handler ...stats.Handler) Option {
	return option(func(opts *options) {
		opts.gRPCStatsHandlers = append(opts.gRPCStatsHandlers, handler...)
	})
}

// AddStatsHandler adds stats handlers.
func (s *Server) AddStatsHandler(handler ...stats.Handler) {
	s.statsHandlers = append(s.statsHandlers, handler...)
}

type statsHandler struct {
	*Server
}

func (s *statsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	for _, h := range s.statsHandlers {
		ctx = h.TagRPC(ctx, info)
	}
	return ctx
}

func (s *statsHandler) HandleRPC(ctx context.Context, stats stats.RPCStats) {
	for _, h := range s.statsHandlers {
		h.HandleRPC(ctx, stats)
	}
}

func (s *statsHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	for _, h := range s.statsHandlers {
		ctx = h.TagConn(ctx, info)
	}
	return ctx
}

func (s *statsHandler) HandleConn(ctx context.Context, stats stats.ConnStats) {
	for _, h := range s.statsHandlers {
		h.HandleConn(ctx, stats)
	}
}
