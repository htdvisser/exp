package sql

import "strings"

// BuildSelect builds part of a select statement for the given columns.
// The table and columns are expected to be safe. DO NOT pass those directly from user input.
func BuildSelect(table string, columns ...string) string {
	if len(columns) == 0 {
		return ""
	}
	n := 2 * (len(columns) - 1) // Separators
	for i := 0; i < len(columns); i++ {
		n += 2 + len(columns[i]) // "column"
		if table != "" {
			n += len(table) + 1
		}
	}
	var b strings.Builder
	b.Grow(n)
	b.WriteString(`"`)
	if table != "" {
		b.WriteString(table)
		b.WriteByte('.')
	}
	b.WriteString(columns[0])
	for _, s := range columns[1:] {
		b.WriteString(`", "`)
		if table != "" {
			b.WriteString(table)
			b.WriteByte('.')
		}
		b.WriteString(s)
	}
	b.WriteString(`"`)
	return b.String()
}

// BuildInsert builds part of an insert statement for the given columns.
// The columns are expected to be safe. DO NOT pass those directly from user input.
func BuildInsert(columns ...string) string {
	if len(columns) == 0 {
		return ""
	}
	n := 12                         // () VALUES ()
	n += 2 * 2 * (len(columns) - 1) // Separators on both sides
	for i := 0; i < len(columns); i++ {
		n += 2 + len(columns[i]) + 1 // "column" on one side, ? on the other
	}
	var b strings.Builder
	b.Grow(n)
	b.WriteString(`("`)
	b.WriteString(columns[0])
	for _, s := range columns[1:] {
		b.WriteString(`", "`)
		b.WriteString(s)
	}
	b.WriteString(`") VALUES (?`)
	for i := 0; i < len(columns)-1; i++ {
		b.WriteString(", ?")
	}
	b.WriteByte(')')
	return b.String()
}

// BuildUpdate builds part of an update statement for the given columns.
// The columns are expected to be safe. DO NOT pass those directly from user input.
func BuildUpdate(columns ...string) string {
	if len(columns) == 0 {
		return ""
	}
	n := 2 * (len(columns) - 1) // Separators
	for i := 0; i < len(columns); i++ {
		n += 2 + len(columns[i]) + 4 // "column" = ?
	}
	var b strings.Builder
	b.Grow(n)
	b.WriteString(`"`)
	b.WriteString(columns[0])
	b.WriteString(`" = ?`)
	for _, s := range columns[1:] {
		b.WriteString(`, "`)
		b.WriteString(s)
		b.WriteString(`" = ?`)
	}
	return b.String()
}
