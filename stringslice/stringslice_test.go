package stringslice

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	got := Filter(
		[]string{"match", "no match", "match", "no match"},
		Equal("match"),
	)
	assert.Equal(t, []string{"match", "match"}, got)
}

func TestUnique(t *testing.T) {
	got := Filter(
		[]string{"match", "no match", "match", "no match"},
		Unique(4),
	)
	assert.Equal(t, []string{"match", "no match"}, got)
}

func TestMatchAny(t *testing.T) {
	any := MatchAny(
		[]string{"match", "no match", "match", "no match"},
		Equal("match"),
	)
	assert.True(t, any)

	none := MatchAny(
		[]string{"no match", "no match"},
		Equal("match"),
	)
	assert.False(t, none)
}

func TestMatchAll(t *testing.T) {
	all := MatchAll(
		[]string{"match", "match"},
		Equal("match"),
	)
	assert.True(t, all)

	some := MatchAll(
		[]string{"match", "no match", "match", "no match"},
		Equal("match"),
	)
	assert.False(t, some)
}

func TestMap(t *testing.T) {
	res := Map(
		[]string{"foo", "bar", "baz"},
		strings.ToUpper,
	)
	assert.Equal(t, []string{"FOO", "BAR", "BAZ"}, res)
}

func TestMapFuncs(t *testing.T) {
	assert.Equal(t, "foobar", AddPrefix("foo")("bar"))
	assert.Equal(t, "foobar", AddSuffix("bar")("foo"))
}
