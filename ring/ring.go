// Package ring implements a ring data structure for efficiently storing the n
// most recent entries added to it.
// The ring data structure is comparable to a ring buffer, but with only one
// pointer for writing.
package ring // import "htdvisser.dev/exp/ring"

// Entry is the Entry type stored in the data structure.
// You may want to change the type from interface{} to something else if you're
// actually using this.
type Entry interface{}

var zero Entry

// Ring implements a ring data structure.
type Ring struct {
	entries []Entry
	pos     int
}

// New creates a new ring data structure of the given capacity.
func New(capacity int) *Ring {
	return &Ring{
		entries: make([]Entry, 0, capacity),
	}
}

// Add adds an entry to the ring, appending it if the ring hasn't reached its
// capacity yet, or otherwise replacing the oldest entry.
func (r *Ring) Add(entry Entry) {
	if len(r.entries) < cap(r.entries) {
		r.entries = append(r.entries, entry)
		r.pos = len(r.entries) - 1
	} else {
		r.pos = (r.pos + 1) % len(r.entries)
		r.entries[r.pos] = entry
	}
}

// Last returns the most recently added entry.
func (r *Ring) Last() (Entry, bool) {
	if len(r.entries) == 0 {
		return zero, false
	}
	return r.entries[r.pos], true
}

// All returns a copy of the contents of the ring, starting with the oldest entry
// and ending with the most recently added entry.
func (r *Ring) All() []Entry {
	entries := make([]Entry, 0, len(r.entries))
	if len(r.entries) < cap(r.entries) {
		return append(entries, r.entries...)
	}
	start := (r.pos + 1) % len(r.entries)
	entries = append(entries, r.entries[start:]...)
	if start > 0 {
		entries = append(entries, r.entries[0:start]...)
	}
	return entries
}
