// Package reflection can be used to support gRPC server reflection on the server.
package reflection

import (
	"google.golang.org/grpc/reflection"
	"htdvisser.dev/exp/backbone/server"
)

// Register registers the gRPC server reflection service on the server.
func Register(s *server.Server) {
	reflection.Register(s.GRPC.Server)
	reflection.Register(s.InternalGRPC.Server)
}
