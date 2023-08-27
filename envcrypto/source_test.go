package envcrypto

import (
	"os"
	"testing"

	"golang.org/x/exp/slices"
)

func testLookup(t *testing.T, s Source, key string, want string, wantOk bool) {
	t.Helper()

	got, gotOk := s.Lookup(key)

	if got != want {
		t.Errorf("Lookup(%q) returned %q, want %q", key, got, want)
	}
	if gotOk != wantOk {
		t.Errorf("Lookup(%q) returned %t, want %t", key, gotOk, wantOk)
	}
}

func sortedClone(s []string) []string {
	if slices.IsSorted(s) {
		return s
	}
	clone := slices.Clone(s)
	slices.Sort(clone)
	return clone
}

func testKeys(t *testing.T, s Source, want []string) {
	t.Helper()

	got := sortedClone(s.Keys())
	want = sortedClone(want)

	if !slices.Equal(got, want) {
		t.Errorf("Keys() returned %v, want %v", got, want)
	}
}

func TestMapSource(t *testing.T) {
	s := MapSource{
		"FOO": "foo",
		"BAR": "bar",
	}

	t.Run("Lookup", func(t *testing.T) {
		testLookup(t, s, "FOO", "foo", true)
		testLookup(t, s, "BAR", "bar", true)
		testLookup(t, s, "BAZ", "", false)
	})

	t.Run("Keys", func(t *testing.T) {
		testKeys(t, s, []string{"FOO", "BAR"})
	})
}

func TestEnvFileSource(t *testing.T) {
	s, err := NewEnvFileSource(os.DirFS("testdata"), "source/source1.env")
	if err != nil {
		t.Fatalf("NewEnvFileSource() returned error: %v", err)
	}

	t.Run("Lookup", func(t *testing.T) {
		testLookup(t, s, "FOO", "foo", true)
		testLookup(t, s, "BAR", "bar", true)
		testLookup(t, s, "BAZ", "", false)
	})

	t.Run("Keys", func(t *testing.T) {
		testKeys(t, s, []string{"FOO", "BAR"})
	})
}

func TestEnvFilesSource(t *testing.T) {
	s, err := NewEnvFilesSource(os.DirFS("testdata"), "source/source1.env", "source/source2.env")
	if err != nil {
		t.Fatalf("NewEnvFilesSource() returned error: %v", err)
	}

	t.Run("Lookup", func(t *testing.T) {
		testLookup(t, s, "FOO", "foo", true)
		testLookup(t, s, "BAR", "bar", true)
		testLookup(t, s, "BAZ", "baz", true)
		testLookup(t, s, "QUX", "", false)
	})

	t.Run("Keys", func(t *testing.T) {
		testKeys(t, s, []string{"FOO", "BAR", "BAZ"})
	})
}
