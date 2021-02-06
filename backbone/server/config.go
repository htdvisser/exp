package server

import "github.com/spf13/pflag"

// Config is the required configuration for the server.
type Config struct {
	ListenHTTP         string
	ListenGRPC         string
	ListenInternalHTTP string
	ListenInternalGRPC string
}

// DefaultConfig returns the default config for the server.
func DefaultConfig() *Config {
	return &Config{
		ListenHTTP:         ":8080",
		ListenGRPC:         ":9090",
		ListenInternalHTTP: "localhost:18080",
		ListenInternalGRPC: "localhost:19090",
	}
}

// Flags returns a flagset that can be added to the command line.
func (c *Config) Flags(prefix string, defaults *Config) *pflag.FlagSet {
	var flags pflag.FlagSet
	if defaults == nil {
		defaults = DefaultConfig()
	}
	flags.StringVar(&c.ListenHTTP, prefix+"http.listen", defaults.ListenHTTP, "Listen address for the HTTP server")
	flags.StringVar(&c.ListenGRPC, prefix+"grpc.listen", defaults.ListenGRPC, "Listen address for the gRPC server")
	flags.StringVar(&c.ListenInternalHTTP, prefix+"internal.http.listen", defaults.ListenInternalHTTP, "Listen address for the internal HTTP server")
	flags.StringVar(&c.ListenInternalGRPC, prefix+"internal.grpc.listen", defaults.ListenInternalGRPC, "Listen address for the internal gRPC server")
	return &flags
}
