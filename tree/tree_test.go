package tree

import (
	"reflect"
	"testing"
)

func TestTree(t *testing.T) {
	var tree Tree

	expectLoad := func(p Path, expectVal interface{}, expectExists bool) {
		t.Helper()
		val, exists := tree.Load(p)
		if exists != expectExists {
			t.Errorf("tree.Load(%q) exists = %v, want %v", p, exists, expectExists)
		}
		if val != expectVal {
			t.Errorf("tree.Load(%q) val = %v, want %v", p, val, expectVal)
		}
	}
	expectStore := func(p Path, val interface{}, expectOld interface{}, expectReplaced bool) {
		t.Helper()
		old, replaced := tree.Store(p, val)
		if replaced != expectReplaced {
			t.Errorf("tree.Store(%q) replaced = %v, want %v", p, replaced, expectReplaced)
		}
		if old != expectOld {
			t.Errorf("tree.Store(%q) old = %v, want %v", p, old, expectOld)
		}
	}
	expectLoadOrStore := func(p Path, val interface{}, expectOld interface{}, expectLoaded bool) {
		t.Helper()
		old, loaded := tree.LoadOrStore(p, val)
		if loaded != expectLoaded {
			t.Errorf("tree.LoadOrStore(%q) loaded = %v, want %v", p, loaded, expectLoaded)
		}
		if old != expectOld {
			t.Errorf("tree.LoadOrStore(%q) old = %v, want %v", p, old, expectOld)
		}
	}
	expectAll := func(p Path, expect []PathValue) {
		t.Helper()
		all := tree.All(p)
		if !reflect.DeepEqual(all, expect) {
			t.Errorf("tree.All(%q) = %v, want %v", p, all, expect)
		}
	}
	expectDelete := func(p Path, expectVal interface{}, expectExists bool) {
		t.Helper()
		val, exists := tree.Delete(p)
		if exists != expectExists {
			t.Errorf("tree.Delete(%q) exists = %v, want %v", p, exists, expectExists)
		}
		if val != expectVal {
			t.Errorf("tree.Delete(%q) val = %v, want %v", p, val, expectVal)
		}
	}

	expectLoad(nil, nil, false)
	expectDelete(nil, nil, false)
	expectStore(nil, 1, nil, false)
	expectLoad(nil, 1, true)

	expectLoad(Path{"a", "b", "c"}, nil, false)
	expectDelete(Path{"a", "b", "c"}, nil, false)
	expectStore(Path{"a", "b", "c"}, 1, nil, false)
	expectLoad(Path{"a", "b", "c"}, 1, true)

	expectLoad(Path{"a", "b"}, nil, false)
	expectDelete(Path{"a", "b"}, nil, false)

	expectStore(Path{"a", "b"}, 1, nil, false)
	expectLoad(Path{"a", "b"}, 1, true)

	expectLoadOrStore(Path{"a", "b", "c"}, 2, 1, true)
	expectLoad(Path{"a", "b", "c"}, 1, true)

	expectStore(Path{"a", "b", "c"}, 2, 1, true)
	expectLoad(Path{"a", "b", "c"}, 2, true)

	expectLoadOrStore(Path{"b", "c", "d"}, 2, 2, false)
	expectLoad(Path{"b", "c", "d"}, 2, true)

	expectAll(Path{"a"}, []PathValue{
		{Path{"a", "b"}, 1},
		{Path{"a", "b", "c"}, 2},
	})

	clone := tree.Clone()
	if !reflect.DeepEqual(clone.All(nil), tree.All(nil)) {
		t.Error("tree.Clone().All() != tree.All()")
	}

	expectDelete(Path{"a", "b"}, 1, true)
	expectDelete(Path{"a", "b", "c"}, 2, true)
}
