package ring

import (
	"reflect"
	"testing"
)

func TestRing(t *testing.T) {
	r := New[int](4)

	for _, tt := range []struct {
		Name       string
		mut        func(r *Ring[int])
		wantLast   int
		wantLastOK bool
		wantAll    []int
	}{
		{
			Name:       "Initial",
			wantLast:   0,
			wantLastOK: false,
			wantAll:    []int{},
		},
		{
			Name: "Add 1 2",
			mut: func(r *Ring[int]) {
				r.Add(1)
				r.Add(2)
			},
			wantLast:   2,
			wantLastOK: true,
			wantAll:    []int{1, 2},
		},
		{
			Name: "Add 3 4",
			mut: func(r *Ring[int]) {
				r.Add(3)
				r.Add(4)
			},
			wantLast:   4,
			wantLastOK: true,
			wantAll:    []int{1, 2, 3, 4},
		},
		{
			Name: "Add 5 6",
			mut: func(r *Ring[int]) {
				r.Add(5)
				r.Add(6)
			},
			wantLast:   6,
			wantLastOK: true,
			wantAll:    []int{3, 4, 5, 6},
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
