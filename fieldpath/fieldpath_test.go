package fieldpath

import (
	"reflect"
	"testing"
)

func TestParseList(t *testing.T) {
	InternFieldPathStrings("a", "b")
	for _, tt := range []struct {
		strs []string
		fps  List
		str  string
	}{
		{
			strs: []string{"a,a.b", "a.b.c"},
			fps:  List{{"a"}, {"a", "b"}, {"a", "b", "c"}},
			str:  "a,a.b,a.b.c",
		},
	} {
		fps, err := ParseList(tt.strs...)
		if err != nil {
			t.Errorf("ParseList(%q) err = %v, want nil", tt.str, err)
		}
		if !reflect.DeepEqual(fps, tt.fps) {
			t.Errorf("ParseList(%q) = %v, want %v", tt.str, fps, tt.fps)
		}
		if got := fps.String(); got != tt.str {
			t.Errorf("fps.String() = %q, want %q", got, tt.str)
		}
	}
}

func TestSort(t *testing.T) {
	fps := List{
		{"a", "b"},
		{"b"},
		{"a", "b", "c"},
		{"c"},
		{"a"},
	}
	got := fps.Sort()
	want := List{
		{"a"},
		{"a", "b"},
		{"a", "b", "c"},
		{"b"},
		{"c"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("fps.Sort() = %v, want %v", got, want)
	}
}

func TestUnique(t *testing.T) {
	fps := List{
		{"a"},
		{"a", "b"},
		{"b"},
		{"a", "b", "c"},
		{"c"},
		{"a"},
	}

	{
		got := fps.Unique(false)
		want := List{
			{"a"},
			{"b"},
			{"c"},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("fps.Unique(false) = %v, want %v", got, want)
		}
	}

	{
		got := fps.Unique(true)
		want := List{
			{"a"},
			{"a", "b"},
			{"a", "b", "c"},
			{"b"},
			{"c"},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("fps.Unique(true) = %v, want %v", got, want)
		}
	}
}

func TestContains(t *testing.T) {
	fps := List{
		{"a"},
		{"a", "b"},
		{"a", "b", "c"},
		{"b"},
		{"c"},
	}
	for _, tt := range []struct {
		search FieldPath
		exact  bool
		want   bool
	}{
		{FieldPath{"b"}, true, true},
		{FieldPath{"a", "b", "c", "d"}, true, false},
		{FieldPath{"a", "b", "c", "d"}, false, true},
		{FieldPath{"d", "e", "f"}, false, false},
	} {
		if got := fps.Contains(tt.search, tt.exact); got != tt.want {
			t.Errorf("fps.Contains(%s, %v) = %v, want %v", tt.search, tt.exact, got, tt.want)
		}
	}
}

func TestContainsOnly(t *testing.T) {
	fps := List{
		{"a"},
		{"a", "b"},
		{"a", "b", "c"},
	}
	for _, tt := range []struct {
		only List
		want bool
	}{
		{List{{"a"}, {"a", "b"}, {"a", "b", "c"}}, true},
		{List{{"a", "b"}}, false},
	} {
		if got := fps.ContainsOnly(tt.only); got != tt.want {
			t.Errorf("fps.ContainsOnly(%s) = %v, want %v", tt.only, got, tt.want)
		}
	}
}

func TestAddPrefix(t *testing.T) {
	fps := List{
		{"b"},
		{"b", "c"},
	}
	want := List{
		{"a", "b"},
		{"a", "b", "c"},
	}
	if got := fps.AddPrefix(FieldPath{"a"}); !reflect.DeepEqual(got, want) {
		t.Errorf(`fps.AddPrefix("a") = %v, want %v`, got, want)
	}
}

func TestRemovePrefix(t *testing.T) {
	fps := List{
		{"a"},
		{"a", "b"},
		{"a", "b", "c"},
		{"b"},
		{"c"},
	}
	want := List{
		{"b"},
		{"b", "c"},
	}
	if got := fps.RemovePrefix(FieldPath{"a"}); !reflect.DeepEqual(got, want) {
		t.Errorf(`fps.RemovePrefix("a") = %v, want %v`, got, want)
	}
}
