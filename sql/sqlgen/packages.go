package main

import (
	"fmt"
	"go/types"

	"github.com/fatih/structtag"
	"golang.org/x/tools/go/packages"
	"htdvisser.dev/exp/stringslice"
)

type Type struct {
	PackagePath string
	Package     string
	Pointer     bool
	Name        string
}

func (t Type) FullName() string {
	if t.Package == "" {
		return t.Name
	}
	return t.Package + "." + t.Name
}

func (t *Type) SetTo(obj types.Object) {
	t.PackagePath = obj.Pkg().Path()
	t.Package = obj.Pkg().Name()
	t.Name = obj.Name()
}

type Field struct {
	Name     string
	Tag      string
	Type     Type
	Ref      *Field
	NullType *NullType
}

func (f Field) Imports() []string {
	imports := make([]string, 0, 2)
	if f.Type.PackagePath != "" {
		imports = append(imports, f.Type.PackagePath)
	}
	if f.NullType != nil && f.NullType.PackagePath != "" {
		imports = append(imports, f.NullType.PackagePath)
	}
	return stringslice.Filter(imports, stringslice.Unique(len(imports)))
}

type NullType struct {
	Type
	Short string
}

func BuildNullType(name string) *NullType {
	return &NullType{
		Type: Type{
			PackagePath: "database/sql",
			Package:     "sql",
			Name:        "Null" + name,
		},
		Short: name,
	}
}

type StructType struct {
	Type
	Fields []Field
}

func (t StructType) Imports() []string {
	imports := make([]string, 0, 1+len(t.Fields)*2)
	if t.Type.PackagePath != "" {
		imports = append(imports, t.Type.PackagePath)
	}
	for _, field := range t.Fields {
		imports = append(imports, field.Imports()...)
	}
	return stringslice.Filter(imports, stringslice.Unique(len(imports)))
}

func BuildStructType(pkg *packages.Package, typeName string) (StructType, error) {
	obj := pkg.Types.Scope().Lookup(typeName)
	if obj == nil {
		return StructType{}, fmt.Errorf(
			"could not find type %q in package %q",
			typeName, pkg.Name,
		)
	}
	var structType StructType
	structType.Type.SetTo(obj)
	namedObj, ok := obj.Type().(*types.Named)
	if !ok {
		return StructType{}, fmt.Errorf(
			"type %q is not a named object",
			typeName,
		)
	}
	structObj, ok := namedObj.Underlying().(*types.Struct)
	if !ok {
		return StructType{}, fmt.Errorf(
			"type %q is not a struct object",
			typeName,
		)
	}
	for i := 0; i < structObj.NumFields(); i++ {
		fieldObj := structObj.Field(i)
		if !fieldObj.Exported() {
			continue
		}
		field := Field{
			Name: fieldObj.Name(),
		}
		tags, err := structtag.Parse(structObj.Tag(i))
		if err != nil {
			return StructType{}, fmt.Errorf(
				"invalid struct tag on field %q of type %q: %w",
				fieldObj.Name(), typeName, err,
			)
		}
		if tag, err := tags.Get(*tagName); err == nil {
			if tag.Name == "-" {
				continue
			}
			field.Tag = tag.Name
		}
		fieldType := fieldObj.Type()
		if ptr, ok := fieldType.(*types.Pointer); ok {
			field.Type.Pointer = true
			fieldType = ptr.Elem()
		}
		switch fieldType := fieldType.(type) {
		case *types.Basic:
			field.Type.Name = fieldType.String()
		case *types.Named:
			field.Type.SetTo(fieldType.Obj())
			if tag, err := tags.Get("ref"); err == nil && tag.Name != "-" {
				structObj := fieldType.Underlying().(*types.Struct)
				for i := 0; i < structObj.NumFields(); i++ {
					fieldObj := structObj.Field(i)
					if !fieldObj.Exported() || fieldObj.Name() != tag.Name {
						continue
					}
					field.Ref = &Field{
						Name: fieldObj.Name(),
					}
					tags, err := structtag.Parse(structObj.Tag(i))
					if err != nil {
						return StructType{}, fmt.Errorf(
							"invalid struct tag on field %q of type %q: %w",
							fieldObj.Name(), typeName, err,
						)
					}
					if tag, err := tags.Get(*tagName); err == nil {
						if tag.Name == "-" {
							continue
						}
						field.Ref.Tag = tag.Name
					}
					fieldType := fieldObj.Type()
					if ptr, ok := fieldType.(*types.Pointer); ok {
						field.Ref.Type.Pointer = true
						fieldType = ptr.Elem()
					}
					switch fieldType := fieldType.(type) {
					case *types.Basic:
						field.Ref.Type.Name = fieldType.String()
					case *types.Named:
						field.Ref.Type.SetTo(fieldType.Obj())
					}
				}
			}
		default:
			return StructType{}, fmt.Errorf("field of unsupported type %q", fieldType)
		}
		if field.Type.Pointer {
			switch field.Type.Name {
			case "bool":
				field.NullType = BuildNullType("Bool")
			case "float64":
				field.NullType = BuildNullType("Float64")
			case "int32":
				field.NullType = BuildNullType("Int32")
			case "int64":
				field.NullType = BuildNullType("Int64")
			case "string":
				field.NullType = BuildNullType("String")
			case "Time":
				field.NullType = BuildNullType("Time")
			}
		}
		structType.Fields = append(structType.Fields, field)
	}
	return structType, nil
}
