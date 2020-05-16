package clicontext_test

import (
	"context"
	"fmt"
	"os"

	"htdvisser.dev/exp/clicontext"
)

func Example() {
	var app interface {
		Run(ctx context.Context, args ...string) error
	}

	// func main()
	{
		ctx, exit := clicontext.WithInterruptAndExit(context.Background())
		defer exit()

		if err := app.Run(ctx, os.Args[1:]...); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	}
}
