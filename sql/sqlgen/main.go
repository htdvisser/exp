// Command sqlgen generates code that makes it easier to scan SQL rows into structs.
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
	"htdvisser.dev/exp/stringslice"
)

const usage = `sqlgen [options] [package] [types...]`

var (
	flags           = pflag.NewFlagSet("sqlgen", pflag.ContinueOnError)
	tagName         = flags.String("tag-name", "json", "Name of the struct tag to extract the field name from")
	fieldMaskSuffix = flags.String("fieldmask-suffix", "FieldMask", "Suffix of the struct that is the field mask")
	pkg             = flags.String("pkg", "", "Package name")
	models          = flags.Bool("models", true, "Generate models")
	settersTo       = flags.String("setters-to", "SetTo", "Name of the method that sets fields to the source struct")
	settersFrom     = flags.String("setters-from", "SetFrom", "Name of the method that sets fields from the source struct")
	pointers        = flags.String("pointers", "Pointers", "Name of the method that returns pointers")
	values          = flags.String("values", "Values", "Name of the method that returns values")
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
			Models:      *models,
			SettersTo:   *settersTo,
			SettersFrom: *settersFrom,
			Pointers:    *pointers,
			Values:      *values,
		},
		Package: lpkgs[0],
		Imports: []string{
			lpkgs[0].PkgPath,
		},
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

		fieldMaskObj := scope.Lookup(typeName + *fieldMaskSuffix)
		if fieldMaskObj == nil {
			return fmt.Errorf(
				"could not find type %q in package %q",
				typeName, data.Package.Name,
			)
		}

		typeData := Type{
			Name:              obj.Name(),
			FullName:          obj.Pkg().Name() + "." + obj.Name(),
			FieldMaskName:     fieldMaskObj.Name(),
			FieldMaskFullName: fieldMaskObj.Pkg().Name() + "." + fieldMaskObj.Name(),
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

			fieldData := Field{Name: field.Name()}

			fieldType := field.Type()
			if ptr, ok := fieldType.(*types.Pointer); ok {
				fieldData.Type = "*"
				fieldType = ptr.Elem()
			}
			switch fieldType := fieldType.(type) {
			case *types.Basic:
				fieldData.Type += fieldType.String()
			case *types.Named:
				data.Imports = append(data.Imports, fieldType.Obj().Pkg().Path())
				fieldData.Type += fieldType.Obj().Pkg().Name() + "." + fieldType.Obj().Name()
			default:
				fieldData.Type = fmt.Sprintf("%#v", fieldType)
			}
			switch fieldData.Type {
			case "*bool":
				data.Imports = append(data.Imports, "database/sql")
				fieldData.Type = "sql.NullBool"
				fieldData.NullType = "Bool"
			case "*float64":
				data.Imports = append(data.Imports, "database/sql")
				fieldData.Type = "sql.NullFloat64"
				fieldData.NullType = "Float64"
			case "*int32":
				data.Imports = append(data.Imports, "database/sql")
				fieldData.Type = "sql.NullInt32"
				fieldData.NullType = "Int32"
			case "*int64":
				data.Imports = append(data.Imports, "database/sql")
				fieldData.Type = "sql.NullInt64"
				fieldData.NullType = "Int64"
			case "*string":
				data.Imports = append(data.Imports, "database/sql")
				fieldData.Type = "sql.NullString"
				fieldData.NullType = "String"
			case "*time.Time":
				data.Imports = append(data.Imports, "database/sql")
				fieldData.Type = "sql.NullTime"
				fieldData.NullType = "Time"
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

	data.Imports = stringslice.Filter(data.Imports, stringslice.Unique(len(data.Imports)))

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
