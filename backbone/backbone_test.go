package backbone_test

import (
	"context"
	"flag"
	"fmt"
	"os"

	"htdvisser.dev/exp/backbone/server"
	"htdvisser.dev/exp/backbone/server/recovery"
	"htdvisser.dev/exp/backbone/server/reflection"
	"htdvisser.dev/exp/clicontext"
	"htdvisser.dev/exp/flagenv"
)

var config struct {
	server server.Config
}

func init() {
	flag.StringVar(&config.server.ListenHTTP, "http.listen", ":8080", "Listen address for the HTTP server")
	flag.StringVar(&config.server.ListenGRPC, "grpc.listen", ":9090", "Listen address for the gRPC server")
	flag.StringVar(&config.server.ListenInternalHTTP, "internal.http.listen", "localhost:18080", "Listen address for the internal HTTP server")
	flag.StringVar(&config.server.ListenInternalGRPC, "internal.grpc.listen", "localhost:19090", "Listen address for the internal gRPC server")
}

func Example() {
	ctx, exit := clicontext.WithInterruptAndExit(context.Background())
	defer exit()

	if err := flagenv.NewParser(flagenv.Prefixes("backbone_")).ParseEnv(flag.CommandLine); err != nil {
		fmt.Fprintln(flag.CommandLine.Output(), err)
		flag.Usage()
		os.Exit(2)
	}

	flag.Parse()

	server := server.New(config.server)

	reflection.Register(server)
	recovery.Register(server)

	// TODO: Register services here.

	if err := server.Run(ctx); err != nil {
		fmt.Fprintln(flag.CommandLine.Output(), err)
		return
	}
}
