package server

import (
	"bufio"
	"context"
	"net"
	"time"

	"github.com/pires/go-proxyproto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"htdvisser.dev/exp/backbone/server"
	"htdvisser.dev/exp/backbone/server/packet"
	"htdvisser.dev/exp/backbone/server/stream"
	echo "htdvisser.dev/exp/echo/api/v1alpha1"
)

type Config struct {
	ListenTCP  string
	TCPTimeout time.Duration
	TCPProxy           bool
	TCPProxyAllowedIPs []string
	tcpServerOptions   []stream.Option
	ListenUDP  string
	Prefix     string
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
