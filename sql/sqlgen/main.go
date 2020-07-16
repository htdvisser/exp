// Command sqlgen generates code that makes it easier to scan SQL rows into structs.
package main

import (
	"bytes"
	"context"
	"fmt"
	"go/format"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/spf13/pflag"
	"golang.org/x/tools/go/packages"
	"htdvisser.dev/exp/clicontext"
	"htdvisser.dev/exp/stringslice"
)

const usage = `sqlgen [options] [package] [type]`

var (
	flags           = pflag.NewFlagSet("sqlgen", pflag.ContinueOnError)
	tagName         = flags.String("tag-name", "json", "Name of the struct tag to extract the field name from")
	plural          = flags.String("plural", "", "Plural form of the type (default: type+s)")
	idFieldName     = flags.String("id-field", "ID", "Name of the ID field")
	fieldMaskSuffix = flags.String("fieldmask-suffix", "FieldMask", "Suffix of the struct that is the field mask")
	pkg             = flags.String("pkg", "", "Package name")
	model           = flags.Bool("model", true, "Generate model")
	modelSuffix     = flags.String("model-suffix", "", "Suffix of the model to generate")
	setterTo        = flags.String("setter-to", "SetTo", "Name of the method that sets fields to the source struct")
	setterFrom      = flags.String("setter-from", "SetFrom", "Name of the method that sets fields from the source struct")
	pointers        = flags.String("pointers", "Pointers", "Name of the method that returns pointers")
	values          = flags.String("values", "Values", "Name of the method that returns values")
	crud            = flags.Bool("crud", true, "Generate CRUD functions")
	table           = flags.String("table", "", "Table name (default: lowercase type+s)")
	out             = flags.StringP("out", "o", "", "Output file (default is STDOUT)")
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
	if len(args) != 2 {
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
			Model:       *model,
			ModelSuffix: *modelSuffix,
			SetterTo:    *setterTo,
			SetterFrom:  *setterFrom,
			Pointers:    *pointers,
			Values:      *values,
			CRUD:        *crud,
			Plural:      *plural,
			IDField:     *idFieldName,
			Table:       *table,
		},
		Imports: []string{
			lpkgs[0].PkgPath,
		},
	}

	if data.Options.PackageName == "" {
		data.Options.PackageName = lpkgs[0].Name
	}
	if data.Options.Plural == "" {
		data.Options.Plural = args[1] + "s"
	}
	if data.Options.Table == "" {
		data.Options.Table = strings.ToLower(data.Options.Plural)
	}

	data.EntityType, err = BuildStructType(lpkgs[0], args[1])
	if err != nil {
		return err
	}
	data.Imports = append(data.Imports, data.EntityType.Imports()...)

	data.FieldMaskType, err = BuildStructType(lpkgs[0], args[1]+*fieldMaskSuffix)
	if err != nil {
		return err
	}
	data.Imports = append(data.Imports, data.FieldMaskType.Imports()...)

	data.Imports = stringslice.Filter(data.Imports, stringslice.Unique(len(data.Imports)))
	sort.Strings(data.Imports)

	for _, field := range data.EntityType.Fields {
		if field.Name == data.Options.IDField {
			data.IDField = field
			break
		}
	}

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
