package stringslice

import (
	"reflect"
	"strings"
	"testing"
)

func TestFilter(t *testing.T) {
	got := Filter(
		[]string{"match", "no match", "match", "no match"},
		Equal("match"),
	)
	want := []string{"match", "match"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Filter(Equal) result is %v, want %v", got, want)
	}
}

func TestUnique(t *testing.T) {
	got := Filter(
		[]string{"match", "no match", "match", "no match"},
		Unique(4),
	)
	want := []string{"match", "no match"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Filter(Unique) result is %v, want %v", got, want)
	}
}

func TestMatchAny(t *testing.T) {
	any := MatchAny(
		[]string{"match", "no match", "match", "no match"},
		Equal("match"),
	)
	if !any {
		t.Errorf("MatchAny result is %v, want true", any)
	}

	none := MatchAny(
		[]string{"no match", "no match"},
		Equal("match"),
	)
	if none {
		t.Errorf("MatchAny result is %v, want false", none)
	}
}

func TestMatchAll(t *testing.T) {
	all := MatchAll(
		[]string{"match", "match"},
		Equal("match"),
	)
	if !all {
		t.Errorf("MatchAll result is %v, want true", all)
	}

	some := MatchAll(
		[]string{"match", "no match", "match", "no match"},
		Equal("match"),
	)
	if some {
		t.Errorf("MatchAll result is %v, want false", some)
	}
}

func TestMap(t *testing.T) {
	got := Map(
		[]string{"foo", "bar", "baz"},
		strings.ToUpper,
	)
	want := []string{"FOO", "BAR", "BAZ"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Map result is %v, want %v", got, want)
	}
}

func TestMapFuncs(t *testing.T) {
	if got := AddPrefix("foo")("bar"); got != "foobar" {
		t.Errorf(`AddPrefix("foo")("bar") = %q, want "foobar"`, got)
	}
	if got := AddSuffix("bar")("foo"); got != "foobar" {
		t.Errorf(`AddSuffix("bar")("foo") = %q, want "foobar"`, got)
	}
}
