package main

import "text/template"

type Options struct {
	PackageName string
	Model       bool
	ModelSuffix string
	SetterTo    string
	SetterFrom  string
	Dialect     string
	Columns     string
	Pointers    string
	Values      string
	CRUD        bool
	Plural      string
	IDField     string
	Table       string
}

type Data struct {
	Options       Options
	Imports       []string
	EntityType    StructType
	FieldMaskType StructType
	IDField       Field
}

var fileTemplate = template.Must(template.New("").Parse(`
{{- $ := .Options -}}
// Code generated by sqlgen. DO NOT EDIT.

package {{ $.PackageName }}

import (
	{{- if $.CRUD }}
	"context"
	"fmt"
	"database/sql"
	hsql "htdvisser.dev/exp/sql"
	{{- end }}
	{{- range .Imports }}
	"{{ . }}"
	{{- end }}
)

{{- if $.Model }}

// {{ .EntityType.Name }}{{ $.ModelSuffix }} is the generated model for {{ .EntityType.FullName }}.
type {{ .EntityType.Name }}{{ $.ModelSuffix }} struct {
	{{- range .EntityType.Fields }}
	{{ .Name }}{{ with .Ref }}{{ .Name }}{{ end }} {{ if .NullType }}{{ .NullType.FullName }}{{ else if .Ref }}{{ .Ref.Type.FullName }}{{ else }}{{ if .Type.Array }}[]{{ end }}{{ .Type.FullName }}{{ end }}
	{{- end }}
}

{{- end }}

{{- if $.SetterTo }}

// {{ $.SetterTo }} sets the selected fields to the {{ .EntityType.FullName }}.
func (m *{{ .EntityType.Name }}{{ $.ModelSuffix }}) {{ $.SetterTo }}(e *{{ .EntityType.FullName }}, mask {{ .FieldMaskType.FullName }}) {
	{{- range .EntityType.Fields }}
	{{- if .Ref }}
	if mask.{{ .Name }} != nil && mask.{{ .Name }}.{{ .Ref.Name }} {
		e.{{ .Name }} = &{{ .Type.FullName }}{
			{{ .Ref.Name }}: m.{{ .Name }}{{ .Ref.Name }},
		}
	}
	{{- else }}
	if mask.{{ .Name }} {
		{{- if .NullType }}
		if m.{{ .Name }}.Valid {
			e.{{ .Name }} = &m.{{ .Name }}.{{ .NullType.Short }}
		} else {
			e.{{ .Name }} = nil
		}
		{{- else }}
		e.{{ .Name }} = m.{{ .Name }}
		{{- end }}
	}
	{{- end }}
	{{- end }}
}

{{- end }}

{{- if $.SetterFrom }}

// {{ $.SetterFrom }} sets the selected fields from the {{ .EntityType.FullName }}.
func (m *{{ .EntityType.Name }}{{ $.ModelSuffix }}) {{ $.SetterFrom }}(e *{{ .EntityType.FullName }}, mask {{ .FieldMaskType.FullName }}) {
	{{- range .EntityType.Fields }}
	{{- if .Ref }}
	if mask.{{ .Name }} != nil && mask.{{ .Name }}.{{ .Ref.Name }} {
		m.{{ .Name }}{{ .Ref.Name }} = e.{{ .Name }}.{{ .Ref.Name }}
	}
	{{- else }}
	if mask.{{ .Name }} {
		{{- if .NullType }}
		m.{{ .Name }} = {{ .NullType.FullName }}{}
		if e.{{ .Name }} != nil {
			m.{{ .Name }}.Valid = true
			m.{{ .Name }}.{{ .NullType.Short }} = *e.{{ .Name }}
		}
		{{- else }}
		m.{{ .Name }} = e.{{ .Name }}
		{{- end }}
	}
	{{- end }}
	{{- end }}
}

{{- end }}

{{- if $.Columns }}

// {{ $.Columns }} returns column names for the selected fields.
func (m *{{ .EntityType.Name }}{{ $.ModelSuffix }}) {{ $.Columns }}(mask {{ .FieldMaskType.FullName }}) []string {
	columns := make([]string, 0, mask.Len())
	{{- range .EntityType.Fields }}
	{{- if .Ref }}
	if mask.{{ .Name }} != nil && mask.{{ .Name }}.{{ .Ref.Name }} {
		columns = append(columns, "{{ .Tag }}_{{ .Ref.Tag }}")
	}
	{{- else }}
	if mask.{{ .Name }} {
		columns = append(columns, "{{ .Tag }}")
	}
	{{- end }}
	{{- end }}
	return columns
}

{{- end }}

{{- if $.Pointers }}

// {{ $.Pointers }} returns pointers to the selected fields.
func (m *{{ .EntityType.Name }}{{ $.ModelSuffix }}) {{ $.Pointers }}(mask {{ .FieldMaskType.FullName }}) []interface{} {
	pointers := make([]interface{}, 0, mask.Len())
	{{- range .EntityType.Fields }}
	{{- if .Ref }}
	if mask.{{ .Name }} != nil && mask.{{ .Name }}.{{ .Ref.Name }} {
		pointers = append(pointers, &m.{{ .Name }}{{ .Ref.Name }})
	}
	{{- else }}
	if mask.{{ .Name }} {
		pointers = append(pointers,
		{{- if .Type.Array -}}
		{{ .Type.Name }}Array(&m.{{ .Name }})
		{{- else -}}
		&m.{{ .Name }}
		{{- end -}}
		)
	}
	{{- end }}
	{{- end }}
	return pointers
}

{{- end }}

{{- if $.Values }}

// {{ $.Values }} returns the values of the selected fields.
func (m *{{ .EntityType.Name }}{{ $.ModelSuffix }}) {{ $.Values }}(mask {{ .FieldMaskType.FullName }}) []interface{} {
	values := make([]interface{}, 0, mask.Len())
	{{- range .EntityType.Fields }}
	{{- if .Ref }}
	if mask.{{ .Name }} != nil && mask.{{ .Name }}.{{ .Ref.Name }} {
		values = append(values, m.{{ .Name }}{{ .Ref.Name }})
	}
	{{- else }}
	if mask.{{ .Name }} {
		values = append(values,
		{{- if .Type.Array -}}
		{{ .Type.Name }}Array(&m.{{ .Name }})
		{{- else -}}
		m.{{ .Name }}
		{{- end -}}
		)
	}
	{{- end }}
	{{- end }}
	return values
}

{{- end }}

{{- if $.CRUD }}

// scan{{ .EntityType.Name }} scans a row with the masked fields into a {{ .EntityType.Name }}.
func scan{{ .EntityType.Name }}(row hsql.Row, mask {{ .FieldMaskType.FullName }}) (*{{ .EntityType.FullName }}, error) {
	var model {{ .EntityType.Name }}{{ $.ModelSuffix }}
	pointers := model.Pointers(mask)
	if err := row.Scan(pointers...); err != nil {
		return nil, err
	}
	var res {{ .EntityType.FullName }}
	model.{{ $.SetterTo }}(&res, mask)
	return &res, nil
}

// scan{{ $.Plural }} scans multiple rows with the masked fields into a slice of {{ $.Plural }}.
func scan{{ $.Plural }}(rows hsql.Rows, mask {{ .FieldMaskType.FullName }}) ([]*{{ .EntityType.FullName }}, error) {
	var model {{ .EntityType.Name }}{{ $.ModelSuffix }}
	pointers := model.Pointers(mask)
	var (
		res []*{{ .EntityType.FullName }}
		err error
	)
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(pointers...); err != nil {
			return nil, err
		}
		var dst {{ .EntityType.FullName }}
		model.{{ $.SetterTo }}(&dst, mask)
		res = append(res, &dst)
	}
	return res, nil
}

// create{{ .EntityType.Name }} creates the {{ .EntityType.Name }}.
func create{{ .EntityType.Name }}(ctx context.Context, db hsql.DB, e *{{ .EntityType.FullName }}, mask {{ .FieldMaskType.FullName }}) error {
	var model {{ .EntityType.Name }}{{ $.ModelSuffix }}
	model.{{ $.SetterFrom }}(e, mask)
	fields, values := ((*{{ .EntityType.Name }})(nil)).{{ $.Columns }}(mask), model.Values(mask)
	query := fmt.Sprintf(
		"INSERT INTO \"{{ $.Table }}\" %s",
		hsql.{{ $.Dialect }}.BuildInsert(fields...),
	)
	_, err := db.ExecContext(ctx, query, values...)
	return err
}

// get{{ .EntityType.Name }}By{{ .IDField.Name }} gets the {{ .EntityType.Name }} by {{ .IDField.Name }}.
func get{{ .EntityType.Name }}By{{ .IDField.Name }}(ctx context.Context, db hsql.DB, {{ .IDField.Tag }} {{ .IDField.Type.FullName }}, mask {{ .FieldMaskType.FullName }}) (*{{ .EntityType.FullName }}, error) {
	return get{{ .EntityType.Name }}Where(ctx, db, "{{ .IDField.Tag }}", {{ .IDField.Tag }}, mask)
}

// get{{ .EntityType.Name }}Where gets the {{ .EntityType.Name }} where column equals value.
// The column is expected to be safe. DO NOT pass it directly from user input.
func get{{ .EntityType.Name }}Where(ctx context.Context, db hsql.DB, column string, value interface{}, mask {{ .FieldMaskType.FullName }}) (*{{ .EntityType.FullName }}, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM \"{{ $.Table }}\" WHERE \"%s\" = $1 LIMIT 1",
		hsql.{{ $.Dialect }}.BuildSelect(((*{{ .EntityType.Name }})(nil)).{{ $.Columns }}(mask)...),
		column,
	)
	rows, err := db.QueryContext(ctx, query, value)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	return scan{{ .EntityType.Name }}{{ $.ModelSuffix }}(rows, mask)
}

// get{{ $.Plural }}By{{ .IDField.Name }} gets multiple {{ $.Plural }} by their {{ .IDField.Name }}s.
func get{{ $.Plural }}By{{ .IDField.Name }}(ctx context.Context, db hsql.DB, {{ .IDField.Tag }}s []{{ .IDField.Type.FullName }}, mask {{ .FieldMaskType.FullName }}) ([]*{{ .EntityType.FullName }}, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM \"{{ $.Table }}\" WHERE \"{{ .IDField.Tag }}\" IN (%s)",
		hsql.{{ $.Dialect }}.BuildSelect(((*{{ .EntityType.Name }})(nil)).{{ $.Columns }}(mask)...),
		hsql.{{ $.Dialect }}.BuildPlaceholders(1, len({{ .IDField.Tag }}s)),
	)
	args := make([]interface{}, len({{ .IDField.Tag }}s))
	for i, {{ .IDField.Tag }} := range {{ .IDField.Tag }}s {
		args[i] = {{ .IDField.Tag }}
	}
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return scan{{ $.Plural }}(rows, mask)
}

// count{{ $.Plural }}Where counts all {{ $.Plural }} in the database.
func count{{ $.Plural }}(ctx context.Context, db hsql.DB) (uint64, error) {
	query := "SELECT COUNT(*) FROM \"{{ $.Table }}\""
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	if !rows.Next() {
		return 0, sql.ErrNoRows
	}
	var count uint64
	if err = rows.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// count{{ $.Plural }}Where counts {{ $.Plural }} where column equals value.
// The column is expected to be safe. DO NOT pass it directly from user input.
func count{{ $.Plural }}Where(ctx context.Context, db hsql.DB, column string, value interface{}) (uint64, error) {
	query := fmt.Sprintf(
		"SELECT COUNT(*) FROM \"{{ $.Table }}\" WHERE \"%s\" = $1",
		column,
	)
	rows, err := db.QueryContext(ctx, query, value)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	if !rows.Next() {
		return 0, sql.ErrNoRows
	}
	var count uint64
	if err = rows.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// list{{ $.Plural }} lists all {{ $.Plural }} in the database.
// The orderBy is expected to be safe. DO NOT pass it directly from user input.
func list{{ $.Plural }}(ctx context.Context, db hsql.DB, mask {{ .FieldMaskType.FullName }}, orderBy string, limit, offset uint) ([]*{{ .EntityType.FullName }}, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM \"{{ $.Table }}\" ORDER BY \"%s\" LIMIT $1 OFFSET $2",
		hsql.{{ $.Dialect }}.BuildSelect(((*{{ .EntityType.Name }})(nil)).{{ $.Columns }}(mask)...),
		orderBy,
	)
	rows, err := db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	return scan{{ $.Plural }}(rows, mask)
}

// list{{ $.Plural }}Where lists {{ $.Plural }} where column equals value.
// The column and orderBy are expected to be safe. DO NOT pass those directly from user input.
func list{{ $.Plural }}Where(ctx context.Context, db hsql.DB, column string, value interface{}, mask {{ .FieldMaskType.FullName }}, orderBy string, limit, offset uint) ([]*{{ .EntityType.FullName }}, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM \"{{ $.Table }}\" WHERE \"%s\" = $1 ORDER BY \"%s\" LIMIT $2 OFFSET $3",
		hsql.{{ $.Dialect }}.BuildSelect(((*{{ .EntityType.Name }})(nil)).{{ $.Columns }}(mask)...),
		column,
		orderBy,
	)
	rows, err := db.QueryContext(ctx, query, value, limit, offset)
	if err != nil {
		return nil, err
	}
	return scan{{ $.Plural }}(rows, mask)
}

// update{{ .EntityType.Name }} updates the masked fields of the {{ .EntityType.Name }}.
// The {{ .EntityType.Name }} is looked up by {{ .IDField.Name }}.
func update{{ .EntityType.Name }}(ctx context.Context, db hsql.DB, e *{{ .EntityType.FullName }}, mask {{ .FieldMaskType.FullName }}) error {
	var model {{ .EntityType.Name }}{{ $.ModelSuffix }}
	model.{{ $.SetterFrom }}(e, mask)
	fields, values := ((*{{ .EntityType.Name }})(nil)).{{ $.Columns }}(mask), model.Values(mask)
	query := fmt.Sprintf(
		"UPDATE \"{{ $.Table }}\" SET %s WHERE \"{{ .IDField.Tag }}\" = $%d",
		hsql.{{ $.Dialect }}.BuildUpdate(fields...), len(fields)+1,
	)
	_, err := db.ExecContext(ctx, query, append(values, e.{{ .IDField.Name }})...)
	return err
}

// delete{{ .EntityType.Name }} deletes the {{ .EntityType.Name }} with the given {{ .IDField.Name }}.
func delete{{ .EntityType.Name }}{{ $.ModelSuffix }}(ctx context.Context, db hsql.DB, {{ .IDField.Tag }} {{ .IDField.Type.FullName }}) error {
	query := "DELETE FROM \"{{ $.Table }}\" WHERE \"{{ .IDField.Tag }}\" = $1"
	_, err := db.ExecContext(ctx, query, {{ .IDField.Tag }})
	return err
}

{{- end }}
`))
