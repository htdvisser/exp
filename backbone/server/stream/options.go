package stream

import "github.com/pires/go-proxyproto"

type options struct {
	middleware    []Middleware
	proxyProtocol bool
	proxyPolicy   proxyproto.PolicyFunc
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt.apply(o)
	}
}

// Option is an option for the stream server.
type Option interface {
	apply(*options)
}

type option func(*options)

func (f option) apply(opts *options) {
	f(opts)
}
