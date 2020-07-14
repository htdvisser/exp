// Command maskgen generates code that helps masking the fields of structs.
package main

import (
	"bytes"
	"context"
	"fmt"
	"go/format"
	"go/types"
	"io"
	"os"

	"github.com/fatih/structtag"
	"github.com/spf13/pflag"
	"golang.org/x/tools/go/packages"
	"htdvisser.dev/exp/clicontext"
)

const usage = `maskgen [options] [package] [types...]`

var (
	flags   = pflag.NewFlagSet("maskgen", pflag.ContinueOnError)
	pkg     = flags.String("pkg", "", "Package name")
	tagName = flags.String("tag-name", "field", "Name of the struct tag to extract the field name from")
	setter  = flags.String("setter", "Set", "Name of the method to set fields from another struct")
	out     = flags.StringP("out", "o", "", "Output file (default is STDOUT)")
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
			PackageName: *pkg,
			Setter:      *setter,
		},
		Package: lpkgs[0],
	}

	scope := lpkgs[0].Types.Scope()

	for _, typeName := range args[1:] {
		obj := scope.Lookup(typeName)
		if obj == nil {
			return fmt.Errorf(
				"could not find type %q in package %q",
				typeName, data.Package.Name,
			)
		}
		structObj, ok := obj.Type().Underlying().(*types.Struct)
		if !ok {
			return fmt.Errorf(
				"type %q is not a struct",
				typeName,
			)
		}

		typeData := Type{
			FullName: obj.Pkg().Name() + "." + obj.Name(),
			Name:     typeName,
		}

		for i := 0; i < structObj.NumFields(); i++ {
			field := structObj.Field(i)
			if !field.Exported() {
				continue
			}
			tags, err := structtag.Parse(structObj.Tag(i))
			if err != nil {
				return fmt.Errorf(
					"invalid struct tag on field %q of type %q: %w",
					field.Name(), typeName, err,
				)
			}
			fieldData := Field{
				Name:  field.Name(),
				Field: field.Name(),
			}
			if tag, err := tags.Get(*tagName); err == nil {
				if tag.Name == "-" {
					continue
				}
				fieldData.Field = tag.Name
			}

			typeData.Fields = append(typeData.Fields, fieldData)
		}

		data.Types = append(data.Types, typeData)
	}

	var buf bytes.Buffer
	if err = fileTemplate.Execute(&buf, data); err != nil {
		return err
	}

	source, err := format.Source(buf.Bytes())
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
