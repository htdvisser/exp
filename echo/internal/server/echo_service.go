package server

import (
	"bufio"
	"context"
	"net"
	"time"

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
	ListenUDP  string
}

func NewEchoService(config Config) *EchoService {
	return &EchoService{config: config}
}

type EchoService struct {
	config Config
	echo.UnimplementedEchoServiceServer
}

func (es *EchoService) Register(ctx context.Context, bbs *server.Server) {
	echo.RegisterEchoServiceServer(bbs.GRPC.Server, es)
	echo.RegisterEchoServiceHandler(ctx, bbs.GRPC.Gateway, bbs.GRPC.LoopbackConn())
	bbs.RegisterTCPServer("Echo-TCP", es.config.ListenTCP, stream.NewServer(es))
	bbs.RegisterUDPServer("Echo-UDP", es.config.ListenUDP, packet.NewServer(es))
}

func (es *EchoService) Echo(ctx context.Context, req *echo.EchoRequest) (*echo.EchoResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %s", err)
	}
	return &echo.EchoResponse{
		Message: req.Message,
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
		if _, err = conn.Write(msg); err != nil {
			return err
		}
	}
}

func (es *EchoService) HandlePacket(ctx context.Context, msg []byte, addr net.Addr, reply func([]byte) error) error {
	return reply(msg)
}
