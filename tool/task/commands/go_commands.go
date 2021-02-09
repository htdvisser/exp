package commands

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
	"htdvisser.dev/exp/stringslice"
)

func moduleDirs() ([]string, error) {
	modFiles, err := filepath.Glob("./**/*.mod")
	if err != nil {
		return nil, err
	}
	return stringslice.Map(modFiles, filepath.Dir), nil
}

func execInModuleDirs(ctx *cli.Context, before, after string, cmd string, args ...string) error {
	moduleDirs, err := moduleDirs()
	if err != nil {
		return err
	}
	for _, moduleDir := range moduleDirs {
		log.Printf("--- %s: %s", moduleDir, before)
		cmd := exec.CommandContext(ctx.Context, cmd, append(args, ctx.Args().Slice()...)...)
		cmd.Dir = moduleDir
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		if err = cmd.Run(); err != nil {
			return err
		}
		log.Printf("--- %s: %s", moduleDir, after)
	}
	return nil
}

func goCommand() *cli.Command {
	return &cli.Command{
		Name:  "go",
		Usage: "Run the same Go command in each Go module",
		Action: func(ctx *cli.Context) error {
			return execInModuleDirs(
				ctx,
				fmt.Sprintf("Running go %s", strings.Join(ctx.Args().Slice(), " ")),
				"Done",
				"go",
			)
		},
	}
}
