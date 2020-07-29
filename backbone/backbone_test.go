package backbone_test

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"htdvisser.dev/exp/backbone/server"
	"htdvisser.dev/exp/backbone/server/recovery"
	"htdvisser.dev/exp/backbone/server/reflection"
	"htdvisser.dev/exp/clicontext"
	"htdvisser.dev/exp/pflagenv"
)

var config struct {
	server server.Config
}

func init() {
	pflag.CommandLine.AddFlagSet(config.server.Flags(""))
}

func Example() {
	ctx, exit := clicontext.WithInterruptAndExit(context.Background())
	defer exit()

	if err := pflagenv.NewParser(pflagenv.Prefixes("backbone_")).ParseEnv(pflag.CommandLine); err != nil {
		fmt.Fprintln(os.Stderr, err)
		pflag.Usage()
		os.Exit(2)
	}

	pflag.Parse()

	server := server.New(config.server)

	reflection.Register(server)
	recovery.Register(server)

	// TODO: Register services here.

	if err := server.Run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}
