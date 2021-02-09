package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"htdvisser.dev/exp/clicontext"
	"htdvisser.dev/exp/tool/task/commands"
)

func main() {
	ctx, exit := clicontext.WithInterruptAndExit(context.Background())
	defer exit()

	app := &cli.App{
		Name:     "task",
		Usage:    "task tool",
		Commands: commands.All(),
	}
	err := app.RunContext(ctx, os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		clicontext.SetExitCode(ctx, 1)
		return
	}
}
