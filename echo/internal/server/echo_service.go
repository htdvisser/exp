package server

import (
	"bufio"
	"context"
	"net"
	"time"

	"github.com/pires/go-proxyproto"
	"github.com/spf13/pflag"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"htdvisser.dev/exp/backbone/server"
	"htdvisser.dev/exp/backbone/server/packet"
	"htdvisser.dev/exp/backbone/server/stream"
	echo "htdvisser.dev/exp/echo/api/v1alpha1"
)

// Config is the configuration for the Echo service.
type Config struct {
	ListenTCP          string
	TCPTimeout         time.Duration
	TCPProxy           bool
	TCPProxyAllowedIPs []string
	tcpServerOptions   []stream.Option
	ListenUDP          string
	Prefix             string
}

// DefaultConfig returns the default configuration for the Echo service.
func DefaultConfig() *Config {
	return &Config{
		ListenTCP:  ":7070",
		TCPTimeout: time.Minute,
		ListenUDP:  ":6060",
		Prefix:     "<echo>: ",
	}
}

// Flags returns a flagset that can be added to the command line.
func (c *Config) Flags(prefix string, defaults *Config) *pflag.FlagSet {
	var flags pflag.FlagSet
	if defaults == nil {
		defaults = DefaultConfig()
	}
	flags.StringVar(&c.ListenTCP, prefix+"tcp.listen", defaults.ListenTCP, "Listen address for the TCP server")
	flags.DurationVar(&c.TCPTimeout, prefix+"tcp.timeout", defaults.TCPTimeout, "Connection timeout for the TCP server")
	flags.BoolVar(&c.TCPProxy, prefix+"tcp.proxy", defaults.TCPProxy, "Support PROXY protocol on the TCP server")
	flags.StringSliceVar(&c.TCPProxyAllowedIPs, prefix+"tcp.proxy.allowed-ips", defaults.TCPProxyAllowedIPs, "Optional list of IPs/CIDRs from which PROXY headers are read")
	flags.StringVar(&c.ListenUDP, prefix+"udp.listen", defaults.ListenUDP, "Listen address for the UDP server")
	flags.StringVar(&c.Prefix, prefix+"prefix", defaults.Prefix, "Prefix for the echo")
	return &flags
}

func NewEchoService(config Config) (*EchoService, error) {
	if config.TCPProxy {
		var (
			policy proxyproto.PolicyFunc
			err    error
		)
		if len(config.TCPProxyAllowedIPs) > 0 {
			policy, err = proxyproto.LaxWhiteListPolicy(config.TCPProxyAllowedIPs)
			if err != nil {
				return nil, err
			}
		}
		config.tcpServerOptions = append(config.tcpServerOptions, stream.WithProxyProtocol(policy))
	}
	return &EchoService{config: config}, nil
}

type EchoService struct {
	config Config
	echo.UnimplementedEchoServiceServer
}

func (es *EchoService) echoBytes(in []byte) []byte {
	if es.config.Prefix == "" {
		return in
	}
	prefix := []byte(es.config.Prefix)
	out := make([]byte, 0, len(prefix)+len(in))
	out = append(out, prefix...)
	out = append(out, in...)
	return out
}

func (es *EchoService) Register(ctx context.Context, bbs *server.Server) {
	echo.RegisterEchoServiceServer(bbs.GRPC.Server, es)
	echo.RegisterEchoServiceHandler(ctx, bbs.GRPC.Gateway, bbs.GRPC.LoopbackConn())
	bbs.RegisterTCPServer("Echo-TCP", es.config.ListenTCP, stream.NewServer(es, es.config.tcpServerOptions...))
	bbs.RegisterUDPServer("Echo-UDP", es.config.ListenUDP, packet.NewServer(es))
}

func (es *EchoService) Echo(ctx context.Context, req *echo.EchoRequest) (*echo.EchoResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %s", err)
	}
	return &echo.EchoResponse{
		Message: es.config.Prefix + req.Message,
	}, nil
}

func (es *EchoService) HandleStream(ctx context.Context, conn net.Conn) error {
	r := bufio.NewReader(conn)
	for {
		conn.SetReadDeadline(time.Now().Add(es.config.TCPTimeout))
		msg, err := r.ReadBytes(byte('\n'))
		if err != nil {
			return err
		}
		conn.SetWriteDeadline(time.Now().Add(es.config.TCPTimeout))
		if _, err = conn.Write(es.echoBytes(msg)); err != nil {
			return err
		}
	}
}

func (es *EchoService) HandlePacket(ctx context.Context, msg []byte, addr net.Addr, reply func([]byte) error) error {
	if es.config.Prefix != "" {
		return reply(append(append([]byte{}, []byte(es.config.Prefix)...), msg...))
	}
	return reply(es.echoBytes(msg))
}
