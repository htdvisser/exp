package sql

import "testing"

func TestMySQL(t *testing.T) {
	t.Run("BuildSelect", func(t *testing.T) {
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
				if got := MySQL.BuildSelect(tc.fields...); got != tc.expect {
					t.Errorf("Expected %q, got %q", tc.expect, got)
				}
			})
		}
	})

	t.Run("TestBuildInsert", func(t *testing.T) {
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
				if got := MySQL.BuildInsert(tc.fields...); got != tc.expect {
					t.Errorf("Expected %q, got %q", tc.expect, got)
				}
			})
		}
	})

	t.Run("BuildUpdate", func(t *testing.T) {
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
				if got := MySQL.BuildUpdate(tc.fields...); got != tc.expect {
					t.Errorf("Expected %q, got %q", tc.expect, got)
				}
			})
		}
	})

	t.Run("BuildPlaceholders", func(t *testing.T) {
		tt := []struct {
			name   string
			start  int
			end    int
			expect string
		}{
			{
				name:   "zero",
				expect: "?",
			},
			{
				name:   "1 2 3",
				start:  1,
				end:    3,
				expect: `?, ?, ?`,
			},
			{
				name:   "99 100 101",
				start:  99,
				end:    101,
				expect: `?, ?, ?`,
			},
		}
		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				if got := MySQL.BuildPlaceholders(tc.start, tc.end); got != tc.expect {
					t.Errorf("Expected %q, got %q", tc.expect, got)
				}
			})
		}
	})
}
