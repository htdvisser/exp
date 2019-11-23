package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"htdvisser.dev/exp/backbone/server"
	echo "htdvisser.dev/exp/echo/api/v1alpha1"
)

func NewEchoService() *EchoService {
	return &EchoService{}
}

type EchoService struct {
	echo.UnimplementedEchoServiceServer
}

func (es *EchoService) Register(ctx context.Context, bbs *server.Server) {
	echo.RegisterEchoServiceServer(bbs.GRPC.Server, es)
	echo.RegisterEchoServiceHandler(ctx, bbs.GRPC.Gateway, bbs.GRPC.LoopbackConn())
}

func (es *EchoService) Echo(ctx context.Context, req *echo.EchoRequest) (*echo.EchoResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %s", err)
	}
	return &echo.EchoResponse{
		Message: req.Message,
	}, nil
}
