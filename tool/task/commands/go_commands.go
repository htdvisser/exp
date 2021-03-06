package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/urfave/cli/v2"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/semver"
)

func moduleDirs() ([]string, error) {
	var modDirs []string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if path == "go.mod" {
			return nil
		}
		if filepath.Base(path) == "go.mod" {
			modDirs = append(modDirs, filepath.Dir(path))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return modDirs, nil
}

func execInModuleDirs(ctx *cli.Context, before, after string, cmd string, args ...string) error {
	moduleDirs, err := moduleDirs()
	if err != nil {
		return err
	}
	for _, moduleDir := range moduleDirs {
		moduleFileName := filepath.Join(moduleDir, "go.mod")
		moduleFileData, err := ioutil.ReadFile(moduleFileName)
		if err != nil {
			return err
		}
		moduleFile, err := modfile.Parse(moduleFileName, moduleFileData, nil)
		if err != nil {
			return err
		}
		if moduleFile.Go == nil {
			return fmt.Errorf("--- %s: SKIP (no Go version in module file)", moduleDir)
		}
		moduleGoVersion := semver.MajorMinor("v" + moduleFile.Go.Version)
		runtimeGoVersion := semver.MajorMinor("v" + strings.TrimPrefix(runtime.Version(), "go"))
		if semver.Compare(runtimeGoVersion, moduleGoVersion) < 0 {
			log.Printf("--- %s: SKIP (module needs Go %s, this is Go %s)", moduleDir, moduleGoVersion, runtimeGoVersion)
			continue
		}
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
