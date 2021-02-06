package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/encoding/protojson"
	bbserver "htdvisser.dev/exp/backbone/server"
	"htdvisser.dev/exp/backbone/server/grpc"
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
	pflag.CommandLine.AddFlagSet(config.server.Flags("", nil))
	pflag.CommandLine.AddFlagSet(config.echo.Flags("", nil))
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

	jsonpb := &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			Multiline:     true,
			Indent:        "  ",
			UseProtoNames: true,
		},
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
