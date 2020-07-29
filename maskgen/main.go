// Command maskgen generates code that helps masking the fields of structs.
package main

import (
	"bytes"
	"context"
	"fmt"
	"go/format"
	"io"
	"os"
	"sort"

	"github.com/spf13/pflag"
	"golang.org/x/tools/go/packages"
	"htdvisser.dev/exp/clicontext"
	"htdvisser.dev/exp/stringslice"
)

const usage = `maskgen [options] [package] [types...]`

var (
	flags         = pflag.NewFlagSet("maskgen", pflag.ContinueOnError)
	pkg           = flags.String("pkg", "", "Package name")
	tagName       = flags.String("tag-name", "field", "Name of the struct tag to extract the field name from")
	setter        = flags.String("setter", "Set", "Name of the method to set fields from another struct")
	jsonMarshaler = flags.String("json-marshaler", "", "Name of the method to marshal JSON")
	out           = flags.StringP("out", "o", "", "Output file (default is STDOUT)")
)

func main() {
	ctx, exit := clicontext.WithInterruptAndExit(context.Background())
	defer exit()

	flags.Usage = func() {
		fmt.Fprintln(os.Stderr, usage)
		flags.PrintDefaults()
	}

	if err := flags.Parse(os.Args[1:]); err != nil {
		switch err {
		case pflag.ErrHelp:
		default:
			fmt.Fprintln(os.Stderr, err)
		}
		return
	}

	args := flags.Args()
	if len(args) < 1 {
		flags.Usage()
		clicontext.SetExitCode(ctx, 2)
		return
	}

	if err := Main(ctx, args...); err != nil {
		fmt.Fprintln(os.Stderr, err)
		clicontext.SetExitCode(ctx, 1)
		return
	}
}

func Main(ctx context.Context, args ...string) (err error) {
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedImports |
			packages.NeedTypes,
		Context: ctx,
	}

	lpkgs, err := packages.Load(cfg, args[0])
	if err != nil {
		return err
	}
	if len(lpkgs) != 1 {
		return fmt.Errorf("found more than one package")
	}

	data := Data{
		Options: Options{
			PackageName:   *pkg,
			Setter:        *setter,
			JSONMarshaler: *jsonMarshaler,
		},
	}

	if data.Options.PackageName == "" {
		data.Options.PackageName = lpkgs[0].Name
	}

	for _, typeName := range args[1:] {
		structType, err := BuildStructType(lpkgs[0], typeName)
		if err != nil {
			return err
		}
		data.Types = append(data.Types, structType)
	}

	uniq := stringslice.Unique(len(data.Imports))
	uniq(data.Options.PackageName)
	data.Imports = stringslice.Filter(data.Imports, uniq)
	sort.Strings(data.Imports)

	var buf bytes.Buffer
	if err = fileTemplate.Execute(&buf, data); err != nil {
		return err
	}

	source := buf.Bytes()

	source, err = format.Source(source)
	if err != nil {
		return err
	}

	var w io.Writer = os.Stdout
	if *out != "" {
		f, err := os.OpenFile(*out, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("could not open %q for writing: %w", *out, err)
		}
		defer func() {
			if closeErr := f.Close(); err == nil {
				err = closeErr
			}
		}()
		w = f
	}

	_, err = w.Write(source)
	return err
}
