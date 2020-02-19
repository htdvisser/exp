package stream

import (
	"net"

	"github.com/pires/go-proxyproto"
)

func (s *Server) withProxy(conn net.Conn) (net.Conn, error) {
	var (
		policy = proxyproto.USE
		err    error
	)
	if s.proxyPolicy != nil {
		policy, err = s.proxyPolicy(conn.RemoteAddr())
		if err != nil {
			return nil, err
		}
	}
	return proxyproto.NewConn(
		conn,
		proxyproto.WithPolicy(policy),
	), nil
}

// WithProxyProtocol returns an option that enables the PROXY protocol and applies the optional policy.
func WithProxyProtocol(policy proxyproto.PolicyFunc) Option {
	return option(func(opts *options) {
		opts.proxyProtocol, opts.proxyPolicy = true, policy
	})
}
