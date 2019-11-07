package ring

import (
	"reflect"
	"testing"
)

func TestRing(t *testing.T) {
	r := New(4)

	for _, tt := range []struct {
		Name       string
		mut        func(r *Ring)
		wantLast   Entry
		wantLastOK bool
		wantAll    []Entry
	}{
		{
			Name:       "Initial",
			wantLast:   zero,
			wantLastOK: false,
			wantAll:    []Entry{},
		},
		{
			Name: "Add 1 2",
			mut: func(r *Ring) {
				r.Add(1)
				r.Add(2)
			},
			wantLast:   Entry(2),
			wantLastOK: true,
			wantAll:    []Entry{1, 2},
		},
		{
			Name: "Add 3 4",
			mut: func(r *Ring) {
				r.Add(3)
				r.Add(4)
			},
			wantLast:   Entry(4),
			wantLastOK: true,
			wantAll:    []Entry{1, 2, 3, 4},
		},
		{
			Name: "Add 5 6",
			mut: func(r *Ring) {
				r.Add(5)
				r.Add(6)
			},
			wantLast:   Entry(6),
			wantLastOK: true,
			wantAll:    []Entry{3, 4, 5, 6},
		},
	} {
		if tt.mut != nil {
			tt.mut(r)
		}

		last, ok := r.Last()
		if last != tt.wantLast || ok != tt.wantLastOK {
			t.Errorf("t.Last() = (%v, %v), want (%v, %v)", last, ok, tt.wantLast, tt.wantLastOK)
		}

		all := r.All()
		if !reflect.DeepEqual(all, tt.wantAll) {
			t.Errorf("t.All() = %v, want %v", all, tt.wantAll)
		}
	}
}
