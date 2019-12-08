package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gogo/gateway"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	bbserver "htdvisser.dev/exp/backbone/server"
	"htdvisser.dev/exp/backbone/server/grpc"
	"htdvisser.dev/exp/backbone/server/jaeger"
	"htdvisser.dev/exp/backbone/server/opentelemetry"
	"htdvisser.dev/exp/backbone/server/prometheus"
	"htdvisser.dev/exp/backbone/server/recovery"
	"htdvisser.dev/exp/backbone/server/reflection"
	"htdvisser.dev/exp/clicontext"
	"htdvisser.dev/exp/echo/internal/server"
	"htdvisser.dev/exp/flagenv"
)

var config struct {
	server bbserver.Config
	echo   server.Config
}

func init() {
	flag.StringVar(&config.server.ListenHTTP, "http.listen", ":8080", "Listen address for the HTTP server")
	flag.StringVar(&config.server.ListenGRPC, "grpc.listen", ":9090", "Listen address for the gRPC server")
	flag.StringVar(&config.server.ListenInternalHTTP, "internal.http.listen", "localhost:18080", "Listen address for the internal HTTP server")
	flag.StringVar(&config.server.ListenInternalGRPC, "internal.grpc.listen", "localhost:19090", "Listen address for the internal gRPC server")
	flag.StringVar(&config.echo.Prefix, "prefix", "<echo>: ", "Prefix for the echo")
	flag.StringVar(&config.echo.ListenTCP, "tcp.listen", ":7070", "Listen address for the TCP server")
	flag.DurationVar(&config.echo.TCPTimeout, "tcp.timeout", time.Minute, "Connection timeout for the TCP server")
	flag.StringVar(&config.echo.ListenUDP, "udp.listen", ":7070", "Listen address for the UDP server")
}

func main() {
	ctx, exit := clicontext.WithInterruptAndExit(context.Background())
	defer exit()

	if err := flagenv.NewParser(flagenv.Prefixes("echo_")).ParseEnv(flag.CommandLine); err != nil {
		fmt.Fprintln(flag.CommandLine.Output(), err)
		flag.Usage()
		os.Exit(2)
	}

	flag.Parse()

	jsonpb := &gateway.JSONPb{
		EmitDefaults: true,
		Indent:       "  ",
		OrigName:     true,
	}

	backbone := bbserver.New(
		config.server,
		bbserver.WithGRPCOptions(
			grpc.WithRuntimeServeMuxOption(
				runtime.WithMarshalerOption(runtime.MIMEWildcard, jsonpb),
			),
		),
	)

	backbone.HTTP.ServeMux.Handle("/api/", http.StripPrefix("/api", backbone.GRPC.Gateway))

	je, err := jaeger.NewExporter()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	opentelemetry.Register(backbone, opentelemetry.WithSyncer(je))
	prometheus.Register(backbone)
	reflection.Register(backbone)
	recovery.Register(backbone)

	echoService := server.NewEchoService(config.echo)
	echoService.Register(ctx, backbone)

	if err := backbone.Run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}
