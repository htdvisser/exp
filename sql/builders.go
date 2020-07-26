package sql

import (
	"strings"
)

type Dialect interface {
	BuildSelect(columns ...string) string
	BuildInsert(columns ...string) string
	BuildUpdate(columns ...string) string
	BuildPlaceholders(start, end int) string
}

func buildSelect(columns ...string) string {
	if len(columns) == 0 {
		return ""
	}
	n := 2 * (len(columns) - 1) // Separators
	for i := range columns {
		n += len(columns[i]) + 2 // "column"
	}
	var b strings.Builder
	b.Grow(n)
	b.WriteByte('"')
	b.WriteString(columns[0])
	for _, s := range columns[1:] {
		b.WriteString(`", "`)
		b.WriteString(s)
	}
	b.WriteByte('"')
	return b.String()
}
