package sql

import "testing"

func TestBuildSelect(t *testing.T) {
	tt := []struct {
		name   string
		fields []string
		expect string
	}{
		{
			name:   "empty",
			fields: nil,
			expect: "",
		},
		{
			name:   "one column",
			fields: []string{"id"},
			expect: `"id"`,
		},
		{
			name:   "two columns",
			fields: []string{"id", "name"},
			expect: `"id", "name"`,
		},
		{
			name:   "three columns",
			fields: []string{"id", "name", "date"},
			expect: `"id", "name", "date"`,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if got := BuildSelect("", tc.fields...); got != tc.expect {
				t.Errorf("Expected %q, got %q", tc.expect, got)
			}
		})
	}
}

func TestBuildSelectTable(t *testing.T) {
	tt := []struct {
		name   string
		fields []string
		expect string
	}{
		{
			name:   "empty",
			fields: nil,
			expect: "",
		},
		{
			name:   "one column",
			fields: []string{"id"},
			expect: `"table.id"`,
		},
		{
			name:   "two columns",
			fields: []string{"id", "name"},
			expect: `"table.id", "table.name"`,
		},
		{
			name:   "three columns",
			fields: []string{"id", "name", "date"},
			expect: `"table.id", "table.name", "table.date"`,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if got := BuildSelect("table", tc.fields...); got != tc.expect {
				t.Errorf("Expected %q, got %q", tc.expect, got)
			}
		})
	}
}

func TestBuildInsert(t *testing.T) {
	tt := []struct {
		name   string
		fields []string
		expect string
	}{
		{
			name:   "empty",
			fields: nil,
			expect: "",
		},
		{
			name:   "one column",
			fields: []string{"id"},
			expect: `("id") VALUES ($1)`,
		},
		{
			name:   "two columns",
			fields: []string{"id", "name"},
			expect: `("id", "name") VALUES ($1, $2)`,
		},
		{
			name:   "three columns",
			fields: []string{"id", "name", "date"},
			expect: `("id", "name", "date") VALUES ($1, $2, $3)`,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if got := BuildInsert(tc.fields...); got != tc.expect {
				t.Errorf("Expected %q, got %q", tc.expect, got)
			}
		})
	}
}

func TestBuildUpdate(t *testing.T) {
	tt := []struct {
		name   string
		fields []string
		expect string
	}{
		{
			name:   "empty",
			fields: nil,
			expect: "",
		},
		{
			name:   "one column",
			fields: []string{"id"},
			expect: `"id" = $1`,
		},
		{
			name:   "two columns",
			fields: []string{"id", "name"},
			expect: `"id" = $1, "name" = $2`,
		},
		{
			name:   "three columns",
			fields: []string{"id", "name", "date"},
			expect: `"id" = $1, "name" = $2, "date" = $3`,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if got := BuildUpdate(tc.fields...); got != tc.expect {
				t.Errorf("Expected %q, got %q", tc.expect, got)
			}
		})
	}
}

func TestBuildPlaceholders(t *testing.T) {
	tt := []struct {
		name   string
		start  int
		end    int
		expect string
	}{
		{
			name:   "zero",
			expect: "$0",
		},
		{
			name:   "1 2 3",
			start:  1,
			end:    3,
			expect: `$1, $2, $3`,
		},
		{
			name:   "99 100 101",
			start:  99,
			end:    101,
			expect: `$99, $100, $101`,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if got := BuildPlaceholders(tc.start, tc.end); got != tc.expect {
				t.Errorf("Expected %q, got %q", tc.expect, got)
			}
		})
	}
}
