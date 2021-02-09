package commands

import (
	"github.com/urfave/cli/v2"
)

func All() []*cli.Command {
	return []*cli.Command{
		goCommand(),
	}
}
