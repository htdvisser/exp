package sql

import (
	"strconv"
	"strings"
)

type postgreSQL struct{}

// PostgreSQL is the PostgreSQL dialect.
var PostgreSQL Dialect = postgreSQL{}

// BuildSelect builds part of a select statement for the given columns.
// The columns are expected to be safe. DO NOT pass those directly from user input.
func (postgreSQL) BuildSelect(columns ...string) string {
	return buildSelect(columns...)
}

func arglen(i int) (n int) {
	n = 1
	for e := i; e >= 10; e /= 10 {
		n++
	}
	return n
}

// BuildInsert builds part of an insert statement for the given columns.
// The columns are expected to be safe. DO NOT pass those directly from user input.
func (postgreSQL) BuildInsert(columns ...string) string {
	if len(columns) == 0 {
		return ""
	}
	n := 12                         // () VALUES ()
	n += 2 * 2 * (len(columns) - 1) // Separators on both sides
	for i := range columns {
		n += len(columns[i]) + 3 + arglen(i+1) // "column" on one side, $i on the other
	}
	var b strings.Builder
	b.Grow(n)
	b.WriteString(`("`)
	b.WriteString(columns[0])
	for _, s := range columns[1:] {
		b.WriteString(`", "`)
		b.WriteString(s)
	}
	b.WriteString(`") VALUES ($1`)
	for i := range columns[1:] {
		b.WriteString(", $")
		b.WriteString(strconv.Itoa(i + 2))
	}
	b.WriteByte(')')
	return b.String()
}

// BuildUpdate builds part of an update statement for the given columns.
// The columns are expected to be safe. DO NOT pass those directly from user input.
func (postgreSQL) BuildUpdate(columns ...string) string {
	if len(columns) == 0 {
		return ""
	}
	n := 2 * (len(columns) - 1) // Separators
	for i := 0; i < len(columns); i++ {
		n += len(columns[i]) + 6 + arglen(i+1) // "column" = $i
	}
	var b strings.Builder
	b.Grow(n)
	b.WriteString(`"`)
	b.WriteString(columns[0])
	b.WriteString(`" = $1`)
	for i, s := range columns[1:] {
		b.WriteString(`, "`)
		b.WriteString(s)
		b.WriteString(`" = $`)
		b.WriteString(strconv.Itoa(i + 2))
	}
	return b.String()
}

// BuildPlaceholders builds placeholders for SQL queries.
func (postgreSQL) BuildPlaceholders(start, end int) string {
	var n int
	for i := start; i <= end; i++ {
		if i > start {
			n += 2
		}
		n += 1 + arglen(i)
	}
	var b strings.Builder
	b.Grow(n)
	for i := start; i < end+1; i++ {
		if i > start {
			b.WriteString(", ")
		}
		b.WriteString(`$`)
		b.WriteString(strconv.Itoa(i))
	}
	return b.String()
}
