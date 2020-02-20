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
	"github.com/spf13/pflag"
	bbserver "htdvisser.dev/exp/backbone/server"
	"htdvisser.dev/exp/backbone/server/grpc"
	"htdvisser.dev/exp/backbone/server/jaeger"
	"htdvisser.dev/exp/backbone/server/opentelemetry"
	"htdvisser.dev/exp/backbone/server/prometheus"
	"htdvisser.dev/exp/backbone/server/recovery"
	"htdvisser.dev/exp/backbone/server/reflection"
	"htdvisser.dev/exp/clicontext"
	"htdvisser.dev/exp/echo/internal/server"
	"htdvisser.dev/exp/pflagenv"
)

var config struct {
	server bbserver.Config
	echo   server.Config
}

func init() {
	pflag.StringVar(&config.server.ListenHTTP, "http.listen", ":8080", "Listen address for the HTTP server")
	pflag.StringVar(&config.server.ListenGRPC, "grpc.listen", ":9090", "Listen address for the gRPC server")
	pflag.StringVar(&config.server.ListenInternalHTTP, "internal.http.listen", "localhost:18080", "Listen address for the internal HTTP server")
	pflag.StringVar(&config.server.ListenInternalGRPC, "internal.grpc.listen", "localhost:19090", "Listen address for the internal gRPC server")
	pflag.StringVar(&config.echo.Prefix, "prefix", "<echo>: ", "Prefix for the echo")
	pflag.StringVar(&config.echo.ListenTCP, "tcp.listen", ":7070", "Listen address for the TCP server")
	pflag.BoolVar(&config.echo.TCPProxy, "tcp.proxy", false, "Support PROXY protocol on the TCP server")
	pflag.StringSliceVar(&config.echo.TCPProxyAllowedIPs, "tcp.proxy.allowed-ips", nil, "Optional list of IPs/CIDRs from which PROXY headers are read")
	pflag.DurationVar(&config.echo.TCPTimeout, "tcp.timeout", time.Minute, "Connection timeout for the TCP server")
	pflag.StringVar(&config.echo.ListenUDP, "udp.listen", ":6060", "Listen address for the UDP server")
}

func main() {
	ctx, exit := clicontext.WithInterruptAndExit(context.Background())
	defer exit()

	if err := pflagenv.NewParser(pflagenv.Prefixes("echo_")).ParseEnv(pflag.CommandLine); err != nil {
		fmt.Fprintln(os.Stderr, err)
		flag.Usage()
		os.Exit(2)
	}

	pflag.Parse()

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

	echoService, err := server.NewEchoService(config.echo)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	echoService.Register(ctx, backbone)

	if err := backbone.Run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}
