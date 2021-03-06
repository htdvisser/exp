// Package fieldpath implements utilities for field paths as used in protobuf field masks.
package fieldpath

import (
	"sort"
	"strings"
)

// List is a list of field paths.
type List []Path

func (fps List) Len() int { return len(fps) }

func (fps List) Swap(i, j int) { fps[i], fps[j] = fps[j], fps[i] }

func (fps List) Less(i, j int) bool {
	li, lj := len(fps[i]), len(fps[j])
	for k := 0; k < li && k < lj; k++ {
		if fps[i][k] < fps[j][k] {
			return true
		}
		if fps[i][k] == fps[j][k] {
			continue
		}
		return false
	}
	return li < lj
}

// ParseList parses field paths. The result is sorted.
func ParseList(s ...string) (List, error) {
	var fps []string
	for _, s := range s {
		fps = append(fps, strings.Split(s, ",")...)
	}
	out := make(List, len(fps))
	var err error
	for i, fp := range fps {
		out[i], err = Parse(fp)
		if err != nil {
			return nil, err
		}
	}
	sort.Sort(out)
	return out, nil
}

// Strings returns a list of dot-separated field paths.
func (fps List) Strings() []string {
	ps := make([]string, len(fps))
	for i, fp := range fps {
		ps[i] = fp.String()
	}
	return ps
}

// String returns a comma-separated list of dot-separated field paths.
func (fps List) String() string {
	return strings.Join(fps.Strings(), ",")
}

// Sort returns a sorted copy of the List.
func (fps List) Sort() List {
	out := make(List, len(fps))
	copy(out, fps)
	sort.Sort(out)
	return out
}

// Filter filters the fieldpaths by predicate p.
func (fps List) Filter(p func(Path) bool) List {
	out := make(List, 0, len(fps))
	for _, fp := range fps {
		if p(fp) {
			out = append(out, fp)
		}
	}
	return out
}

// Unique returns a List containing the unique paths in the List. If the List
// contains a field and a prefix of that field, only the prefix will be in the
// result, unless exact is true. The result is sorted.
func (fps List) Unique(exact bool) List {
	out := make(List, 0, len(fps))
	for _, fp := range fps.Sort() {
		if len(out) > 0 {
			last := out[len(out)-1]
			if fp.Equal(last) || (!exact && fp.HasPrefix(last)) {
				continue
			}
		}
		out = append(out, fp)
	}
	return out
}

// MatchAny returns true if any element of the list matches predicate p.
func (fps List) MatchAny(p func(Path) bool) bool {
	for _, fp := range fps {
		if p(fp) {
			return true
		}
	}
	return false
}

// Contains returns true if fps contains search or a prefix of search (if exact is false).
func (fps List) Contains(search Path, exact bool) bool {
	return fps.MatchAny(func(fp Path) bool {
		return fp.Equal(search) || (!exact && search.HasPrefix(fp))
	})
}

// MatchAll returns true if all elements of the list match predicate p.
func (fps List) MatchAll(p func(Path) bool) bool {
	for _, fp := range fps {
		if !p(fp) {
			return false
		}
	}
	return true
}

// ContainsOnly returns true if the list contains only field paths present in search.
func (fps List) ContainsOnly(search List) bool {
	return fps.MatchAll(func(fp Path) bool {
		return search.Contains(fp, true)
	})
}

// Map returns a List containing the results of calling m on every element of fps.
func (fps List) Map(m func(Path) Path) List {
	out := make(List, len(fps))
	for i, fp := range fps {
		out[i] = m(fp)
	}
	return out
}

// AddPrefix returns a List with all elements of fps with the given prefix prepended.
func (fps List) AddPrefix(prefix Path) List {
	return fps.Map(func(fp Path) Path { return prefix.Join(fp...) })
}

// RemovePrefix returns a List with all elements of fps that have the given prefix,
// but without that prefix.
func (fps List) RemovePrefix(prefix Path) List {
	return fps.Filter(func(fp Path) bool {
		return fp.HasPrefix(prefix)
	}).Map(func(fp Path) Path {
		wp := make(Path, len(fp)-len(prefix))
		copy(wp, fp[len(prefix):])
		return wp
	})
}

// Path is the path to a field in a struct.
type Path []string

// Parse parses a field path.
func Parse(s string) (Path, error) {
	fp := strings.Split(s, ".")
	for i, e := range fp {
		if interned, ok := internedPathElements[e]; ok {
			fp[i] = interned
		}
	}
	return fp, nil
}

// String returns the field path as a dot-separated string.
func (fp Path) String() string {
	return strings.Join(fp, ".")
}

// Join returns a Path that joins f together with the extra elements.
func (fp Path) Join(elements ...string) Path {
	newPath := make(Path, len(fp)+len(elements))
	if len(fp) > 0 {
		copy(newPath[:len(fp)], fp)
	}
	for i, e := range elements {
		if interned, ok := internedPathElements[e]; ok {
			newPath[len(fp)+i] = interned
		} else {
			newPath[len(fp)+i] = e
		}
	}
	return newPath
}

// Equal returns whether f is equal to other.
func (fp Path) Equal(other Path) bool {
	if len(other) != len(fp) {
		return false
	}
	for i, e := range other {
		if fp[i] != e {
			return false
		}
	}
	return true
}

// HasPrefix returns whether f has other as a prefix.
func (fp Path) HasPrefix(other Path) bool {
	if len(other) >= len(fp) {
		return false
	}
	for i, e := range other {
		if fp[i] != e {
			return false
		}
	}
	return true
}

var internedPathElements = make(map[string]string)

// InternPathStrings interns the field path strings of t.
// This func is typically called from init().
// It is not safe for concurrent use.
func InternPathStrings(s ...string) error {
	fps, err := ParseList(s...)
	if err != nil {
		return err
	}
	for _, fp := range fps {
		for _, e := range fp {
			internedPathElements[e] = e
		}
	}
	return nil
}
