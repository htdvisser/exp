package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	bbgrpc "htdvisser.dev/exp/backbone/client/grpc"
	"htdvisser.dev/exp/clicontext"
	echo "htdvisser.dev/exp/echo/api/v1alpha1"
	"htdvisser.dev/exp/flagenv"
)

var config struct {
	server struct {
		GRPCAddress string
		GRPCTLS     bool
	}
}

func init() {
	flag.StringVar(&config.server.GRPCAddress, "grpc.address", "localhost:9090", "Address of the gRPC server")
	flag.BoolVar(&config.server.GRPCTLS, "grpc.tls", false, "Use TLS to connect to the gRPC server")
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

	if err := Main(ctx, flag.Args()...); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Main(ctx context.Context, args ...string) error {
	if config.server.GRPCTLS {
		ctx = bbgrpc.NewContextWithDialOptions(ctx, grpc.WithTransportCredentials(credentials.NewTLS(nil)))
	} else {
		ctx = bbgrpc.NewContextWithDialOptions(ctx, grpc.WithInsecure())
	}

	cc, err := bbgrpc.DialContext(ctx, config.server.GRPCAddress)
	if err != nil {
		return err
	}
	defer cc.Close()

	res, err := echo.NewEchoServiceClient(cc).Echo(ctx, &echo.EchoRequest{
		Message: strings.Join(args, " "),
	})
	if err != nil {
		return err
	}
	fmt.Println(res.Message)

	return nil
}
