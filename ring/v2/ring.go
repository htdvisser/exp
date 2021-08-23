// Package ring implements a ring data structure for efficiently storing the n
// most recent entries added to it.
// The ring data structure is comparable to a ring buffer, but with only one
// pointer for writing.
package ring

// Ring implements a ring data structure.
type Ring[T any] struct {
	entries []T
	pos     int
}

// New creates a new ring data structure of the given capacity.
func New[T any](capacity int) *Ring[T] {
	return &Ring[T]{
		entries: make([]T, 0, capacity),
	}
}

// Add adds an entry to the ring, appending it if the ring hasn't reached its
// capacity yet, or otherwise replacing the oldest entry.
func (r *Ring[T]) Add(entry T) {
	if len(r.entries) < cap(r.entries) {
		r.entries = append(r.entries, entry)
		r.pos = len(r.entries) - 1
	} else {
		r.pos = (r.pos + 1) % len(r.entries)
		r.entries[r.pos] = entry
	}
}

// Last returns the most recently added entry.
func (r *Ring[T]) Last() (T, bool) {
	var zero T
	if len(r.entries) == 0 {
		return zero, false
	}
	return r.entries[r.pos], true
}

// All returns a copy of the contents of the ring, starting with the oldest entry
// and ending with the most recently added entry.
func (r *Ring[T]) All() []T {
	entries := make([]T, 0, len(r.entries))
	if len(r.entries) < cap(r.entries) {
		return append(entries, r.entries...)
	}
	start := (r.pos + 1) % len(r.entries)
	entries = append(entries, r.entries[start:]...)
	if start > 0 {
		entries = append(entries, r.entries[:start]...)
	}
	return entries
}
