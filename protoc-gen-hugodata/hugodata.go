package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/envoyproxy/protoc-gen-validate/validate"
	pgs "github.com/lyft/protoc-gen-star"
	"google.golang.org/genproto/googleapis/api/annotations"
	yaml "gopkg.in/yaml.v2"
)

type HugoDataModule struct {
	*pgs.ModuleBase

	packages map[string]pgs.Package
}

func HugoData() *HugoDataModule {
	return &HugoDataModule{
		ModuleBase: &pgs.ModuleBase{},
		packages:   make(map[string]pgs.Package),
	}
}

func (m *HugoDataModule) Name() string { return "hugodata" }

func (m *HugoDataModule) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {
	for _, pkg := range packages {
		m.generatePackage(pkg)
	}
	return m.Artifacts()
}

func EntityName(entity pgs.Entity) pgs.Name {
	return pgs.Name(strings.TrimPrefix(entity.FullyQualifiedName(), "."+entity.Package().ProtoName().String()+"."))
}

func (m *HugoDataModule) generatePackage(pkg pgs.Package) {
	var (
		enums    yaml.MapSlice
		messages yaml.MapSlice
		services yaml.MapSlice
	)

	for _, file := range pkg.Files() {
		if !file.BuildTarget() {
			continue
		}
		for _, enum := range file.AllEnums() {
			enums = append(enums, yaml.MapItem{
				Key:   EntityName(enum).String(),
				Value: BuildEnum(enum),
			})
		}
		for _, message := range file.AllMessages() {
			if message.IsMapEntry() {
				continue
			}
			messages = append(messages, yaml.MapItem{
				Key:   EntityName(message).String(),
				Value: BuildMessage(message),
			})
		}
		for _, service := range file.Services() {
			services = append(services, yaml.MapItem{
				Key:   service.Name().String(),
				Value: BuildService(service),
			})
		}
	}

	sort.Sort(mapSliceByKey(enums))
	sort.Sort(mapSliceByKey(messages))
	sort.Sort(mapSliceByKey(services))

	basePath := []string{"api", pkg.ProtoName().String()}

	var buf bytes.Buffer

	if len(enums) > 0 {
		if err := yaml.NewEncoder(&buf).Encode(enums); err != nil {
			m.AddError(err.Error())
		} else {
			m.OverwriteCustomFile(m.JoinPath(append(basePath, "enums.yml")...), buf.String(), 0644)
		}
		buf.Reset()
	}

	if len(messages) > 0 {
		if err := yaml.NewEncoder(&buf).Encode(messages); err != nil {
			m.AddError(err.Error())
		} else {
			m.OverwriteCustomFile(m.JoinPath(append(basePath, "messages.yml")...), buf.String(), 0644)
		}
		buf.Reset()
	}

	if len(services) > 0 {
		if err := yaml.NewEncoder(&buf).Encode(services); err != nil {
			m.AddError(err.Error())
		} else {
			m.OverwriteCustomFile(m.JoinPath(append(basePath, "services.yml")...), buf.String(), 0644)
		}
		buf.Reset()
	}
}

type mapSliceByKey yaml.MapSlice

func (a mapSliceByKey) Len() int           { return len(a) }
func (a mapSliceByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a mapSliceByKey) Less(i, j int) bool { return a[i].Key.(string) < a[j].Key.(string) }

type Entity struct {
	src     pgs.Entity
	Name    pgs.Name `yaml:"name"`
	Comment string   `yaml:"comment,omitempty"`
}

func cleanComments(comments string) string {
	commentLines := strings.Split(comments, "\n")
	hasLeadingSpace := true
	for _, line := range commentLines {
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, " ") {
			hasLeadingSpace = false
			break
		}
	}
	if !hasLeadingSpace {
		return strings.TrimRightFunc(comments, unicode.IsSpace)
	}
	for i, line := range commentLines {
		commentLines[i] = strings.TrimPrefix(line, " ")
	}
	return strings.TrimRightFunc(strings.Join(commentLines, "\n"), unicode.IsSpace)
}

func BuildEntity(src pgs.Entity) Entity {
	comments := src.SourceCodeInfo().LeadingComments()
	if comments == "" {
		comments = src.SourceCodeInfo().TrailingComments()
	}
	entity := Entity{
		src:     src,
		Name:    src.Name(),
		Comment: cleanComments(comments),
	}
	return entity
}

type Ref struct {
	src     pgs.Entity
	Package pgs.Name `yaml:"package,omitempty"`
	Name    pgs.Name `yaml:"name"`
}

func BuildRef(src pgs.Entity) Ref {
	ref := Ref{
		src:  src,
		Name: EntityName(src),
	}
	if !src.BuildTarget() {
		ref.Package = src.Package().ProtoName()
	}
	return ref
}

type EnumValue struct {
	src    pgs.EnumValue
	Entity `yaml:",inline"`
	Value  int32 `yaml:"value"`
}

func BuildEnumValue(src pgs.EnumValue) EnumValue {
	value := EnumValue{
		src:    src,
		Entity: BuildEntity(src),
		Value:  src.Value(),
	}
	return value
}

type Enum struct {
	src    pgs.Enum
	Entity `yaml:",inline"`
	Values []EnumValue
}

func BuildEnum(src pgs.Enum) Enum {
	enum := Enum{
		src:    src,
		Entity: BuildEntity(src),
	}
	enum.Name = EntityName(src)
	for _, value := range src.Values() {
		enum.Values = append(enum.Values, BuildEnumValue(value))
	}
	return enum
}

type FieldTypeElem struct {
	Type    string     `yaml:"type,omitempty"`
	Enum    Ref        `yaml:"enum,omitempty"`
	Message Ref        `yaml:"message,omitempty"`
	Rules   FieldRules `yaml:"rules,omitempty"`
}

func ProtoTypeString(t pgs.ProtoType) string {
	switch t {
	case pgs.DoubleT:
		return "double"
	case pgs.FloatT:
		return "float"
	case pgs.Int64T:
		return "int64"
	case pgs.UInt64T:
		return "uint64"
	case pgs.Int32T:
		return "int32"
	case pgs.Fixed64T:
		return "fixed64"
	case pgs.Fixed32T:
		return "fixed32"
	case pgs.BoolT:
		return "bool"
	case pgs.StringT:
		return "string"
	case pgs.BytesT:
		return "bytes"
	case pgs.UInt32T:
		return "uint32"
	case pgs.SFixed32:
		return "sfixed32"
	case pgs.SFixed64:
		return "sfixed64"
	case pgs.SInt32:
		return "sint32"
	case pgs.SInt64:
		return "sint64"
	default:
		panic(fmt.Errorf("unexpected ProtoType %q", t))
	}
}

type Bytes []byte

func (b Bytes) MarshalYAML() (interface{}, error) {
	return base64.StdEncoding.EncodeToString(b), nil
}

func ProtoTypeDefault(t pgs.ProtoType) interface{} {
	switch t {
	case pgs.DoubleT:
		return float64(0.0)
	case pgs.FloatT:
		return float32(0.0)
	case pgs.Int64T:
		return int64(0)
	case pgs.UInt64T:
		return uint64(0)
	case pgs.Int32T:
		return int32(0)
	case pgs.Fixed64T:
		return uint64(0)
	case pgs.Fixed32T:
		return uint32(0)
	case pgs.BoolT:
		return false
	case pgs.StringT:
		return ""
	case pgs.BytesT:
		return Bytes{}
	case pgs.UInt32T:
		return uint32(0)
	case pgs.SFixed32:
		return int32(0)
	case pgs.SFixed64:
		return int64(0)
	case pgs.SInt32:
		return int32(0)
	case pgs.SInt64:
		return int64(0)
	default:
		panic(fmt.Errorf("unexpected ProtoType %q", t))
	}
}

type PGSFieldType interface {
	ProtoType() pgs.ProtoType
	IsEmbed() bool
	IsEnum() bool
	Enum() pgs.Enum
	Embed() pgs.Message
}

func BuildFieldTypeElem(src PGSFieldType) FieldTypeElem {
	fieldTypeElem := FieldTypeElem{}
	switch {
	case src.IsEnum():
		fieldTypeElem.Enum = BuildRef(src.Enum())
	case src.IsEmbed():
		fieldTypeElem.Message = BuildRef(src.Embed())
	default:
		fieldTypeElem.Type = ProtoTypeString(src.ProtoType())
	}
	return fieldTypeElem
}

type FieldType struct {
	FieldTypeElem `yaml:",inline"`
	Repeated      *FieldTypeElem `yaml:"repeated,omitempty"`
	MapKey        *FieldTypeElem `yaml:"map_key,omitempty"`
	MapValue      *FieldTypeElem `yaml:"map_value,omitempty"`
}

func BuildFieldType(src pgs.FieldType) FieldType {
	fieldType := FieldType{}
	switch {
	case src.IsRepeated():
		elem := BuildFieldTypeElem(src.Element())
		fieldType.Repeated = &elem
	case src.IsMap():
		key := BuildFieldTypeElem(src.Key())
		fieldType.MapKey = &key
		elem := BuildFieldTypeElem(src.Element())
		fieldType.MapValue = &elem
	default:
		fieldType.FieldTypeElem = BuildFieldTypeElem(src)
	}
	return fieldType
}

func BuildFieldDefault(src pgs.FieldType) interface{} {
	switch {
	case src.IsRepeated():
		return []interface{}{}
	case src.IsMap():
		return map[string]interface{}{}
	case src.IsEnum():
		return src.Enum().Values()[0].Name().String()
	case src.IsEmbed():
		if src.Embed().IsWellKnown() && strings.HasSuffix(src.Embed().WellKnownType().Name().String(), "Value") {
			return nil
		}
		switch src.Embed().WellKnownType() {
		case pgs.AnyWKT:
			return nil
		case pgs.DurationWKT:
			return "0s"
		case pgs.TimestampWKT:
			return "0001-01-01T00:00:00Z"
		}
		return map[string]interface{}{}
	default:
		return ProtoTypeDefault(src.ProtoType())
	}
}

type Field struct {
	src       pgs.Field
	Entity    `yaml:",inline"`
	FieldType `yaml:",inline"`
	Default   interface{} `yaml:"default"`
}

func BuildField(src pgs.Field) Field {
	field := Field{
		src:       src,
		Entity:    BuildEntity(src),
		FieldType: BuildFieldType(src.Type()),
		Default:   BuildFieldDefault(src.Type()),
	}
	var fieldRules validate.FieldRules
	if ok, _ := src.Extension(validate.E_Rules, &fieldRules); ok {
		field.AddFieldRules(&fieldRules)
	}
	return field
}

type OneOf struct {
	src        pgs.OneOf
	Entity     `yaml:",inline"`
	FieldNames []pgs.Name `yaml:"field_names,omitempty"`
}

func BuildOneOf(src pgs.OneOf) OneOf {
	oneof := OneOf{
		src:    src,
		Entity: BuildEntity(src),
	}
	for _, field := range src.Fields() {
		oneof.FieldNames = append(oneof.FieldNames, field.Name())
	}
	return oneof
}

type Message struct {
	src    pgs.Message
	Entity `yaml:",inline"`
	Fields []Field `yaml:"fields,omitempty"`
	OneOfs []OneOf `yaml:"oneofs,omitempty"`
}

func BuildMessage(src pgs.Message) Message {
	message := Message{
		src:    src,
		Entity: BuildEntity(src),
	}
	message.Name = EntityName(src)
	for _, field := range src.Fields() {
		message.Fields = append(message.Fields, BuildField(field))
	}
	for _, oneof := range src.OneOfs() {
		message.OneOfs = append(message.OneOfs, BuildOneOf(oneof))
	}
	return message
}

type Stream struct {
	Ref    `yaml:",inline"`
	Stream bool `yaml:"stream,omitempty"`
}

type Method struct {
	src    pgs.Method
	Entity `yaml:",inline"`
	Input  Stream     `yaml:"input"`
	Output Stream     `yaml:"output"`
	HTTP   []HTTPRule `yaml:"http,omitempty"`
}

func BuildMethod(src pgs.Method) Method {
	method := Method{
		src:    src,
		Entity: BuildEntity(src),
		Input: Stream{
			Ref:    BuildRef(src.Input()),
			Stream: src.ClientStreaming(),
		},
		Output: Stream{
			Ref:    BuildRef(src.Output()),
			Stream: src.ServerStreaming(),
		},
	}
	var httpRules annotations.HttpRule
	if ok, _ := src.Extension(annotations.E_Http, &httpRules); ok {
		method.AddHTTPRules(&httpRules)
	}
	return method
}

type Service struct {
	src     pgs.Service
	Entity  `yaml:",inline"`
	Methods yaml.MapSlice `yaml:"methods,omitempty"`
}

func BuildService(src pgs.Service) Service {
	service := Service{
		src:    src,
		Entity: BuildEntity(src),
	}
	for _, method := range src.Methods() {
		service.Methods = append(service.Methods, yaml.MapItem{
			Key:   method.Name().String(),
			Value: BuildMethod(method),
		})
	}
	return service
}
