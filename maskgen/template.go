package main

import (
	"text/template"
)

type Options struct {
	PackageName   string
	Setter        string
	JSONMarshaler string
}

type Data struct {
	Options Options
	Imports []string
	Types   []StructType
}

var fileTemplate = template.Must(template.New("").Parse(`
{{- $ := .Options -}}
// Code generated by maskgen. DO NOT EDIT.

package {{ $.PackageName }}

import (
	"fmt"
	"htdvisser.dev/exp/fieldpath"
	{{- range .Imports }}
	"{{ . }}"
	{{- end }}
)

{{- range .Types }}

// {{ .Name }}FieldMask masks the fields of {{ .Name }}.
type {{ .Name }}FieldMask struct {
	{{- range .Fields }}
	{{ .Name }} {{ with .MaskType.Name }}*{{ . }}{{ else }}bool{{ end }} {{ if .JSONTag.Key }}` + "`" + `{{ .JSONTag.String }}` + "`" + `{{ end }}
	{{- end}}
}

func (m *{{ .Name }}FieldMask) set(selected bool, fields ...fieldpath.Path) error {
	for _, field := range fields {
		if len(field) == 0 {
			continue
		}
		switch field[0] {
		{{- range .Fields }}
		case "{{ .Tag }}":
			{{- if .MaskType.Name }}
			if m.{{ .Name }} == nil {
				m.{{ .Name }} = &{{ .MaskType.Name }}{}
			}
			if len(field) == 1 {
				m.{{ .Name }}.setAll(selected)
			} else {
				m.{{ .Name }}.set(selected, field[1:])
			}
			{{- else }}
			m.{{ .Name }} = selected
			{{- end }}
		{{- end }}
		default:
			return fmt.Errorf("no field %q in {{ .Name }}", field)
		}
	}
	return nil
}

func (m *{{ .Name }}FieldMask) setAll(selected bool) {
	{{- range .Fields }}
	{{- if .MaskType.Name }}
	if selected {
		m.{{ .Name }} = &{{ .MaskType.Name }}{}
		m.{{ .Name }}.setAll(selected)
	} else {
		m.{{ .Name }} = nil
	}
	{{- else }}
	m.{{ .Name }} = selected
	{{- end }}
	{{- end }}
}

// Select selects the given fields in the field mask.
func (m *{{ .Name }}FieldMask) Select(fields ...fieldpath.Path) error {
	return m.set(true, fields...)
}

// SelectAll selects all fields in the field mask.
func (m *{{ .Name }}FieldMask) SelectAll() {
	m.setAll(true)
}

// Unselect unselects the given fields in the field mask.
func (m *{{ .Name }}FieldMask) Unselect(fields ...fieldpath.Path) error {
	return m.set(false, fields...)
}

// UnselectAll unselects all fields in the field mask.
func (m *{{ .Name }}FieldMask) UnselectAll() {
	m.setAll(false)
}

// IsAll returns whether the entire mask is selected.
func (m {{ .Name }}FieldMask) IsAll() bool {
	{{- range .Fields }}
	{{- if .MaskType.Name }}
	if m.{{ .Name }} == nil || !m.{{ .Name }}.IsAll() {
		return false
	}
	{{- else }}
	if !m.{{ .Name }} {
		return false
	}
	{{- end }}
	{{- end }}
	return true
}

// Len returns the number of selected fields. If all fields of a subfield are
// selected, that subfield is counted as 1.
func (m {{ .Name }}FieldMask) Len() int {
	var count int
	{{- range .Fields }}
	{{- if .MaskType.Name }}
	if m.{{ .Name }} != nil {
		if m.{{ .Name }}.IsAll() {
			count++
		} else {
			count += m.{{ .Name }}.Len()
		}
	}
	{{- else }}
	if m.{{ .Name }} {
		count++
	}
	{{- end }}
	{{- end }}
	return count
}

// Fields returns the selected fields.
func (m {{ .Name }}FieldMask) Fields() fieldpath.List {
	fields := make(fieldpath.List, 0, m.Len())
	{{- range .Fields }}
	{{- if .MaskType.Name }}
	if m.{{ .Name }} != nil {
		if m.{{ .Name }}.IsAll() {
			fields = append(fields, fieldpath.Path{"{{ .Tag }}"})
		} else {
			subFields := m.{{ .Name }}.Fields()
			fields = append(fields, subFields.AddPrefix(fieldpath.Path{"{{ .Tag }}"})...)
		}
	}
	{{- else }}
	if m.{{ .Name }} {
		fields = append(fields, fieldpath.Path{"{{ .Tag }}"})
	}
	{{- end }}
	{{- end }}
	return fields
}

{{- if $.Setter }}

// {{ $.Setter }} sets the selected fields from src to e.
func (e *{{ .Name }}) {{ $.Setter }}(src *{{ .Name }}, mask {{ .Name }}FieldMask) {
	if src == nil {
		src = &{{ .Name }}{}
	}
	{{- range .Fields }}
	{{- if .MaskType.Name }}
	if mask.{{ .Name }} != nil {
		if mask.IsAll() {
			e.{{ .Name }} = src.{{ .Name }}
		} else {
			if e.{{ .Name }} == nil {
				e.{{ .Name }} = &{{ .Type.Name }}{}
			}
			e.{{ .Name }}.Set(src.{{ .Name }}, *mask.{{ .Name }})
		}
	}
	{{- else }}
	if mask.{{ .Name }} {
		e.{{ .Name }} = src.{{ .Name }}
	}
	{{- end }}
	{{- end }}
}
{{- end }}

{{- if $.JSONMarshaler }}

// {{ $.JSONMarshaler }} marshals the selected fields of e to JSON.
// Any "omitempty" options in JSON struct tags of {{ .Name }} are ignored.
func (e *{{ .Name }}) {{ $.JSONMarshaler }}(mask {{ .Name }}FieldMask) ([]byte, error) {
	if e == nil {
		return []byte{'n', 'u', 'l', 'l'}, nil
	}
	s := streamPool.BorrowStream(nil)
	defer streamPool.ReturnStream(s)
	s.WriteObjectStart()
	var isNotFirst bool
	{{- range .Fields }}
	{{- if eq .JSONTag.Name "-" }}
	// Omit {{ .Name }} from JSON.
	{{- else }}
	{{- if .MaskType.Name }}
	if mask.{{ .Name }} != nil {
		if isNotFirst {
			s.WriteMore()
		}
		s.WriteObjectField("{{ with .JSONTag.Name }}{{ . }}{{ else }}{{ .Tag }}{{ end }}")
		sub, err := e.{{ .Name }}.{{ $.JSONMarshaler }}(*mask.{{ .Name }})
		if err != nil {
			return nil, err
		}
		s.Write(sub)
		isNotFirst = true
	}
	{{- else }}
	if mask.{{ .Name }} {
		if isNotFirst {
			s.WriteMore()
		}
		s.WriteObjectField("{{ with .JSONTag.Name }}{{ . }}{{ else }}{{ .Tag }}{{ end }}")
		s.WriteVal(e.{{ .Name }})
		isNotFirst = true
	}
	{{- end }}
	{{- end }}
	{{- end }}
	s.WriteObjectEnd()
	return s.Buffer(), nil
}
{{- end }}

{{- end }}
`))
