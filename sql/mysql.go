package sql

import (
	"strings"
)

type mySQL struct{}

// MySQL is the MySQL dialect.
var MySQL Dialect = mySQL{}

// BuildSelect builds part of a select statement for the given columns.
// The columns are expected to be safe. DO NOT pass those directly from user input.
func (mySQL) BuildSelect(columns ...string) string {
	return buildSelect(columns...)
}

// BuildInsert builds part of an insert statement for the given columns.
// The columns are expected to be safe. DO NOT pass those directly from user input.
func (mySQL) BuildInsert(columns ...string) string {
	if len(columns) == 0 {
		return ""
	}
	n := 12                         // () VALUES ()
	n += 2 * 2 * (len(columns) - 1) // Separators on both sides
	for i := range columns {
		n += len(columns[i]) + 3 // "column" on one side, ? on the other
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
	for range columns[1:] {
		b.WriteString(", ?")
	}
	b.WriteByte(')')
	return b.String()
}

// BuildUpdate builds part of an update statement for the given columns.
// The columns are expected to be safe. DO NOT pass those directly from user input.
func (mySQL) BuildUpdate(columns ...string) string {
	if len(columns) == 0 {
		return ""
	}
	n := 2 * (len(columns) - 1) // Separators
	for i := 0; i < len(columns); i++ {
		n += len(columns[i]) + 6 // "column" = ?
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

// BuildPlaceholders builds placeholders for SQL queries.
func (mySQL) BuildPlaceholders(start, end int) string {
	var n int
	for i := start; i <= end; i++ {
		if i > start {
			n += 2
		}
		n += 1
	}
	var b strings.Builder
	b.Grow(n)
	for i := start; i < end+1; i++ {
		if i > start {
			b.WriteString(", ")
		}
		b.WriteString(`?`)
	}
	return b.String()
}
