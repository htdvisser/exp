// Package stringslice provides some utilities on top of []string for lazy developers.
package stringslice

// MatchFunc is a function that matches elements in a string slice.
type MatchFunc func(string) bool

// Filter returns a slice containing the elements of slice for which match returns true.
func Filter(slice []string, match MatchFunc) []string {
	out := make([]string, 0, len(slice))
	for _, s := range slice {
		if match(s) {
			out = append(out, s)
		}
	}
	return out
}

// Unique returns a filter function that returns true for an element if it is
// the first time the function has been called with that element.
func Unique(size int) MatchFunc {
	seen := make(map[string]struct{}, size)
	return func(s string) bool {
		if _, seen := seen[s]; seen {
			return false
		}
		seen[s] = struct{}{}
		return true
	}
}

// MatchAny returns true if match returns true for any element of slice.
func MatchAny(slice []string, match MatchFunc) bool {
	for _, s := range slice {
		if match(s) {
			return true
		}
	}
	return false
}

// MatchAll returns true if match returns true for all elements of slice.
func MatchAll(slice []string, match MatchFunc) bool {
	for _, s := range slice {
		if !match(s) {
			return false
		}
	}
	return true
}

// Equal returns a filter function that returns true if the element equals the given argument.
func Equal(t string) MatchFunc {
	return func(s string) bool { return s == t }
}

// MapFunc is a function that maps a string to a string.
type MapFunc func(string) string

// Map returns a slice containing the result of mapping for each element in the slice.
func Map(slice []string, mapping MapFunc) []string {
	out := make([]string, len(slice))
	for i, s := range slice {
		out[i] = mapping(s)
	}
	return out
}

// AddPrefix returns a map function that adds the given prefix to the element.
func AddPrefix(prefix string) MapFunc {
	return func(s string) string { return prefix + s }
}

// AddSuffix returns a map function that adds the given suffix to the element.
func AddSuffix(suffix string) MapFunc {
	return func(s string) string { return s + suffix }
}
