package fieldpath

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	errNilMap                *ErrNilMap
	errEmptyFieldPath        *ErrEmptyFieldPath
	errNoFieldAtPath         *ErrNoFieldAtPath
	errUnexpectedValueAtPath *ErrUnexpectedValueAtPath
)

func TestFields(t *testing.T) {
	if Map(nil).Fields() != nil {
		t.Fatal("Expected nil fields from nil map")
	}

	m := Map{
		"hello": "world",
		"foo": map[string]interface{}{
			"bar": "baz",
			"baz": "qux",
		},
		"ans": 42,
	}

	fps := m.Fields()
	if diff := cmp.Diff(fps, List{
		Path{"ans"},
		Path{"foo", "bar"},
		Path{"foo", "baz"},
		Path{"hello"},
	}); diff != "" {
		t.Errorf("Result not as expected: %v", diff)
	}
}

func TestGet(t *testing.T) {
	m := Map{
		"hello": "world",
		"foo": map[string]interface{}{
			"bar": "baz",
			"baz": "qux",
		},
		"ans": 42,
	}

	_, err := Map(nil).Get(Path{"hello"})
	if !errors.Is(err, errNilMap) {
		t.Fatal("Expected error when getting from nil map")
	}

	_, err = m.Get(nil)
	if !errors.Is(err, errEmptyFieldPath) {
		t.Fatal("Expected error when getting from empty field path")
	}

	v, err := m.Get(Path{"hello"})
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(v, "world"); diff != "" {
		t.Errorf("Result not as expected: %v", diff)
	}

	v, err = m.Get(Path{"ans"})
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(v, 42); diff != "" {
		t.Errorf("Result not as expected: %v", diff)
	}

	v, err = m.Get(Path{"foo", "bar"})
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(v, "baz"); diff != "" {
		t.Errorf("Result not as expected: %v", diff)
	}

	_, err = m.Get(Path{"bar"})
	if !errors.Is(err, errNoFieldAtPath) {
		t.Fatal("Expected error when getting unset field")
	}

	_, err = m.Get(Path{"hello", "world"})
	if !errors.Is(err, errUnexpectedValueAtPath) {
		t.Fatalf("Expected error when getting field of wrong type, got: %T", err)
	}
}

func TestSet(t *testing.T) {
	m := Map{
		"hello": "world",
	}

	err := Map(nil).Set(Path{"hello"}, "world")
	if !errors.Is(err, errNilMap) {
		t.Fatal("Expected error when setting to nil map")
	}

	err = m.Set(nil, "world")
	if !errors.Is(err, errEmptyFieldPath) {
		t.Fatal("Expected error when setting at empty field path")
	}

	err = m.Set(Path{"hello"}, "you")
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(m, Map{
		"hello": "you",
	}); diff != "" {
		t.Errorf("Result not as expected: %v", diff)
	}

	err = m.Set(Path{"foo", "bar"}, "baz")
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(m, Map{
		"hello": "you",
		"foo": map[string]interface{}{
			"bar": "baz",
		},
	}); diff != "" {
		t.Errorf("Result not as expected: %v", diff)
	}

	err = m.Set(Path{"foo", "baz"}, "qux")
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(m, Map{
		"hello": "you",
		"foo": map[string]interface{}{
			"bar": "baz",
			"baz": "qux",
		},
	}); diff != "" {
		t.Errorf("Result not as expected: %v", diff)
	}

	err = m.Set(Path{"hello", "world"}, "foo")
	if !errors.Is(err, errUnexpectedValueAtPath) {
		t.Fatal("Expected error when setting field of wrong type")
	}
}

func TestUnset(t *testing.T) {
	m := Map{
		"hello": "world",
		"foo": map[string]interface{}{
			"bar": "baz",
			"baz": "qux",
		},
		"ans": 42,
	}

	err := Map(nil).Unset(Path{"hello"})
	if !errors.Is(err, errNilMap) {
		t.Fatal("Expected error when unsetting from nil map")
	}

	err = m.Unset(nil)
	if !errors.Is(err, errEmptyFieldPath) {
		t.Fatal("Expected error when unsetting empty field path")
	}

	err = m.Unset(Path{"hello", "world"})
	if !errors.Is(err, errUnexpectedValueAtPath) {
		t.Fatal("Expected error when unsetting field of wrong type")
	}

	err = m.Unset(Path{"hello"})
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(m, Map{
		"foo": map[string]interface{}{
			"bar": "baz",
			"baz": "qux",
		},
		"ans": 42,
	}); diff != "" {
		t.Errorf("Result not as expected: %v", diff)
	}

	err = m.Unset(Path{"hello", "world"})
	if err != nil {
		t.Fatal(err)
	}

	err = m.Unset(Path{"hello"})
	if err != nil {
		t.Fatal(err)
	}

	err = m.Unset(Path{"foo", "bar", "baz"})
	if !errors.Is(err, errUnexpectedValueAtPath) {
		t.Fatal("Expected error when unsetting field of wrong type")
	}

	err = m.Unset(Path{"foo", "bar"})
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(m, Map{
		"foo": map[string]interface{}{
			"baz": "qux",
		},
		"ans": 42,
	}); diff != "" {
		t.Errorf("Result not as expected: %v", diff)
	}

	err = m.Unset(Path{"foo"})
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(m, Map{
		"ans": 42,
	}); diff != "" {
		t.Errorf("Result not as expected: %v", diff)
	}
}

func TestSetFrom(t *testing.T) {
	m := Map{
		"hello": "world",
		"foo": map[string]interface{}{
			"bar": "baz",
			"baz": "qux",
		},
		"ans": 42,
	}
	out := make(Map)

	err := out.SetFrom(m, m.Fields()...)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(out, Map{
		"hello": "world",
		"foo": map[string]interface{}{
			"bar": "baz",
			"baz": "qux",
		},
		"ans": 42,
	}); diff != "" {
		t.Errorf("Result not as expected: %v", diff)
	}

	out["foo"] = "bar"

	err = out.SetFrom(m, m.Fields()...)
	if !errors.Is(err, errUnexpectedValueAtPath) {
		t.Fatal("Expected error when setting field of wrong type")
	}

	err = out.SetFrom(m, Path{"qux"})
	if !errors.Is(err, errNoFieldAtPath) {
		t.Fatal("Expected error when getting unset field")
	}
}
