// Package channelz can be used to expose the channelz service on the gRPC server.
package channelz

import (
	"google.golang.org/grpc/channelz/service"
	"htdvisser.dev/exp/backbone/server/grpc"
)

// Register registers the gRPC channelz service on the gRPC server.
func Register(s *grpc.Server) {
	service.RegisterChannelzServiceToServer(s.Server)
}
