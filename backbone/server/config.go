package server

import "github.com/spf13/pflag"

// Config is the required configuration for the server.
type Config struct {
	ListenHTTP         string
	ListenGRPC         string
	ListenInternalHTTP string
	ListenInternalGRPC string
}

// Flags returns a flagset that can be added to the command line.
func (c *Config) Flags(prefix string) *pflag.FlagSet {
	var flags pflag.FlagSet
	flags.StringVar(&c.ListenHTTP, prefix+"http.listen", ":8080", "Listen address for the HTTP server")
	flags.StringVar(&c.ListenGRPC, prefix+"grpc.listen", ":9090", "Listen address for the gRPC server")
	flags.StringVar(&c.ListenInternalHTTP, prefix+"internal.http.listen", "localhost:18080", "Listen address for the internal HTTP server")
	flags.StringVar(&c.ListenInternalGRPC, prefix+"internal.grpc.listen", "localhost:19090", "Listen address for the internal gRPC server")
	return &flags
}
