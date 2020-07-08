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
			expect: `("id") VALUES (?)`,
		},
		{
			name:   "two columns",
			fields: []string{"id", "name"},
			expect: `("id", "name") VALUES (?, ?)`,
		},
		{
			name:   "three columns",
			fields: []string{"id", "name", "date"},
			expect: `("id", "name", "date") VALUES (?, ?, ?)`,
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
			expect: `"id" = ?`,
		},
		{
			name:   "two columns",
			fields: []string{"id", "name"},
			expect: `"id" = ?, "name" = ?`,
		},
		{
			name:   "three columns",
			fields: []string{"id", "name", "date"},
			expect: `"id" = ?, "name" = ?, "date" = ?`,
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
