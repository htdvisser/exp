package sql

import (
	"strconv"
	"strings"
)

// BuildSelect builds part of a select statement for the given columns.
// The table and columns are expected to be safe. DO NOT pass those directly from user input.
func BuildSelect(table string, columns ...string) string {
	if len(columns) == 0 {
		return ""
	}
	n := 2 * (len(columns) - 1) // Separators
	for i := range columns {
		n += len(columns[i]) + 2 // "column"
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
	for i := range columns {
		n += len(columns[i]) + 4 // "column" on one side, $i on the other
		if i+1 >= 10 {
			n++
		}
		if i+1 >= 100 {
			n++
		}
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
func BuildUpdate(columns ...string) string {
	if len(columns) == 0 {
		return ""
	}
	n := 2 * (len(columns) - 1) // Separators
	for i := 0; i < len(columns); i++ {
		n += len(columns[i]) + 7 // "column" = $i
		if i+1 >= 10 {
			n++
		}
		if i+1 >= 100 {
			n++
		}
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

// BuildPlaceholders builds placeholders
func BuildPlaceholders(start, end int) string {
	var n int
	for i := start; i <= end; i++ {
		if i > start {
			n += 2
		}
		n += 2
		if i >= 10 {
			n++
		}
		if i >= 100 {
			n++
		}
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
