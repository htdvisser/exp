package fieldpath

import "fmt"

// ErrNilMap indicates that an operation can not be executed with a nil map.
type ErrNilMap struct{}

func (*ErrNilMap) Error() string {
	return "nil map"
}

// Is returns whether the target error is a *ErrNilMap.
func (*ErrNilMap) Is(target error) bool {
	if _, ok := target.(*ErrNilMap); ok {
		return true
	}
	return false
}

// ErrEmptyFieldPath indicates that an operation can not be executed with an empty field path.
type ErrEmptyFieldPath struct{}

func (*ErrEmptyFieldPath) Error() string {
	return "empty field path"
}

// Is returns whether the target error is a *ErrEmptyFieldPath.
func (*ErrEmptyFieldPath) Is(target error) bool {
	if _, ok := target.(*ErrEmptyFieldPath); ok {
		return true
	}
	return false
}

// ErrNoFieldAtPath indicates that no field is present at the path.
type ErrNoFieldAtPath struct {
	Path Path
}

func (e *ErrNoFieldAtPath) Error() string {
	return fmt.Sprintf("no field at path %q", e.Path)
}

// Is returns whether the target error is a *ErrNoFieldAtPath.
// If the paths of both errors are non-nil, it also checks equality of those paths.
func (e *ErrNoFieldAtPath) Is(target error) bool {
	if t, ok := target.(*ErrNoFieldAtPath); ok {
		return e == nil || t == nil || e.Path.Equal(t.Path)
	}
	return false
}

// ErrUnexpectedValueAtPath indicates that an unexpected value was present at the path.
type ErrUnexpectedValueAtPath struct {
	Path  Path
	Value interface{}
}

func (e *ErrUnexpectedValueAtPath) Error() string {
	return fmt.Sprintf("unexpected value at path %q: %T instead of map[string]interface{}", e.Path, e.Value)
}

// Is returns whether the target error is a *ErrUnexpectedValueAtPath.
// If the paths of both errors are non-nil, it also checks equality of those paths.
func (e *ErrUnexpectedValueAtPath) Is(target error) bool {
	if t, ok := target.(*ErrUnexpectedValueAtPath); ok {
		return e == nil || t == nil || e.Path.Equal(t.Path)
	}
	return false
}

// Map is a map[string]interface{} that understands field paths.
type Map map[string]interface{}

// Fields returns the paths of the fields present in m.
// The result is sorted.
func (m Map) Fields() List {
	return fields(m, nil)
}

func fields(m map[string]interface{}, pos Path) List {
	if m == nil {
		return nil
	}
	list := make(List, 0)
	for k, v := range m {
		fp := pos.Join(k)
		switch v := v.(type) {
		case Map:
			list = append(list, fields(v, fp)...)
		case map[string]interface{}:
			list = append(list, fields(v, fp)...)
		default:
			list = append(list, fp)
		}
	}
	return list.Sort()
}

// Get gets the value of the field at fp.
func (m Map) Get(fp Path) (interface{}, error) {
	return get(m, nil, fp)
}

func get(m map[string]interface{}, pos, fp Path) (interface{}, error) {
	if m == nil {
		return nil, &ErrNilMap{}
	}
	if len(fp) < 1 {
		return nil, &ErrEmptyFieldPath{}
	}
	nextpos := pos.Join(fp[0])
	mf, ok := m[fp[0]]
	if !ok {
		return nil, &ErrNoFieldAtPath{Path: nextpos}
	}
	if len(fp) == 1 {
		return mf, nil
	}
	switch sub := mf.(type) {
	case Map:
		return get(sub, nextpos, fp[1:])
	case map[string]interface{}:
		return get(sub, nextpos, fp[1:])
	default:
		return nil, &ErrUnexpectedValueAtPath{Path: nextpos, Value: sub}
	}
}

// Set sets the value of the field at fp to v.
func (m Map) Set(fp Path, v interface{}) error {
	return set(m, nil, fp, v)
}

func set(m map[string]interface{}, pos, fp Path, v interface{}) error {
	if m == nil {
		return &ErrNilMap{}
	}
	if len(fp) < 1 {
		return &ErrEmptyFieldPath{}
	}
	if len(fp) == 1 {
		m[fp[0]] = v
		return nil
	}
	nextpos := pos.Join(fp[0])
	mf, ok := m[fp[0]]
	if !ok {
		sub := make(map[string]interface{})
		m[fp[0]] = sub
		return set(sub, nextpos, fp[1:], v)
	}
	switch sub := mf.(type) {
	case Map:
		return set(sub, nextpos, fp[1:], v)
	case map[string]interface{}:
		return set(sub, nextpos, fp[1:], v)
	default:
		return &ErrUnexpectedValueAtPath{Path: nextpos, Value: sub}
	}
}

// Unset unsets the value of the field at fp.
func (m Map) Unset(fp Path) error {
	return unset(m, nil, fp)
}

func unset(m map[string]interface{}, pos, fp Path) error {
	if m == nil {
		return &ErrNilMap{}
	}
	if len(fp) < 1 {
		return &ErrEmptyFieldPath{}
	}
	if len(fp) == 1 {
		delete(m, fp[0])
		return nil
	}
	nextpos := pos.Join(fp[0])
	mf, ok := m[fp[0]]
	if !ok {
		return nil
	}
	switch sub := mf.(type) {
	case Map:
		return unset(sub, nextpos, fp[1:])
	case map[string]interface{}:
		return unset(sub, nextpos, fp[1:])
	default:
		return &ErrUnexpectedValueAtPath{Path: nextpos, Value: sub}
	}
}

// SetFrom sets values at fps in src to m.
func (m Map) SetFrom(src Map, fps ...Path) error {
	for _, fp := range fps {
		v, err := src.Get(fp)
		if err != nil {
			return err
		}
		if err = m.Set(fp, v); err != nil {
			return err
		}
	}
	return nil
}
