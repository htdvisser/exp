// Package tree implements a tree structure.
package tree

import "sort"

// Path is the path of an element in the tree.
type Path []string

type leaf struct {
	val interface{}
}

type node struct {
	path  Path
	edges []*edge
	leaf  *leaf
}

func (n *node) getEdge(element string) *edge {
	i := sort.Search(len(n.edges), func(i int) bool {
		return n.edges[i].element >= element
	})
	if i < len(n.edges) && n.edges[i].element == element {
		return n.edges[i]
	}
	return nil
}

func (n *node) getNode(p Path) *node {
	if len(p) == 0 {
		return n
	}
	edge := n.getEdge(p[0])
	if edge == nil {
		return nil
	}
	return edge.node.getNode(p[1:])
}

func (n *node) getOrCreateEdge(element string) (e *edge, exists bool) {
	i := sort.Search(len(n.edges), func(i int) bool {
		return n.edges[i].element >= element
	})
	if i < len(n.edges) && n.edges[i].element == element {
		return n.edges[i], true
	}
	e = &edge{element: element, node: &node{}}
	n.edges = append(n.edges, e)
	copy(n.edges[i+1:], n.edges[i:])
	n.edges[i] = e
	return e, false
}

func (n *node) getOrCreateNode(p Path) *node {
	if len(p) == 0 {
		return n
	}
	edge, _ := n.getOrCreateEdge(p[0])
	return edge.node.getOrCreateNode(p[1:])
}

func (n *node) delEdge(element string) (e *edge, exists bool) {
	i := sort.Search(len(n.edges), func(i int) bool {
		return n.edges[i].element >= element
	})
	if i < len(n.edges) && n.edges[i].element == element {
		e = n.edges[i]
		n.edges = append(n.edges[:i], n.edges[i+1:]...)
		return e, true
	}
	return nil, false
}

func (n *node) del(p Path) (val interface{}, deleted bool) {
	if len(p) == 0 {
		n.path = nil
		if n.leaf != nil {
			val, n.leaf = n.leaf.val, nil
			return val, true
		}
		return nil, false
	}
	edge := n.getEdge(p[0])
	if edge == nil {
		return nil, false
	}
	val, deleted = edge.node.del(p[1:])
	if len(edge.node.edges) == 0 && edge.node.leaf == nil {
		n.delEdge(p[0])
	}
	return
}

func (n *node) walk(fn WalkFunc) error {
	if n.leaf != nil {
		if err := fn(n.path, n.leaf.val); err != nil {
			return err
		}
	}
	for _, e := range n.edges {
		if err := e.node.walk(fn); err != nil {
			return err
		}
	}
	return nil
}

func (n *node) clone() *node {
	clone := &node{path: n.path}
	if len(n.edges) > 0 {
		clone.edges = make([]*edge, len(n.edges))
		for i, e := range n.edges {
			clone.edges[i] = &edge{
				element: e.element,
				node:    e.node.clone(),
			}
		}
	}
	if n.leaf != nil {
		clone.leaf = &leaf{val: n.leaf.val}
	}
	return clone
}

type edge struct {
	element string
	node    *node
}

// Tree structure storing values at paths.
// It is not safe for concurrent use by multiple goroutines without locking.
type Tree struct {
	root *node
}

func (t *Tree) init() {
	if t.root == nil {
		t.root = &node{}
	}
}

// Load returns the value stored in the tree under the given path, if it
// exists.
func (t *Tree) Load(p Path) (val interface{}, ok bool) {
	t.init()
	node := t.root.getNode(p)
	if node == nil || node.leaf == nil {
		return nil, false
	}
	return node.leaf.val, true
}

// Store stores the given value in the tree under the given path.
// It returns the old value if an existing value was replaced.
func (t *Tree) Store(p Path, val interface{}) (old interface{}, replaced bool) {
	t.init()
	node := t.root.getOrCreateNode(p)
	if node.leaf != nil {
		old, replaced = node.leaf.val, true
	}
	node.path, node.leaf = p, &leaf{val}
	return old, replaced
}

// LoadOrStore returns the value stored in the tree under the given path,
// if it exists. Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (t *Tree) LoadOrStore(p Path, val interface{}) (actual interface{}, loaded bool) {
	t.init()
	node := t.root.getOrCreateNode(p)
	if node.leaf != nil {
		loaded = true
	} else {
		node.leaf = &leaf{val}
	}
	return node.leaf.val, loaded
}

// Delete deletes the value stored in the tree under the given path, if it
// exists. It returns the deleted value if it was deleted.
func (t *Tree) Delete(p Path) (val interface{}, deleted bool) {
	t.init()
	return t.root.del(p)
}

// Walk walks the tree under the given prefix. If f returns a non-nil error,
// walk stops and returns that error.
func (t *Tree) Walk(prefix Path, f WalkFunc) error {
	node := t.root.getNode(prefix)
	if node == nil {
		return nil
	}
	return node.walk(f)
}

// WalkFunc is called by Walk.
type WalkFunc func(p Path, val interface{}) error

// PathValue represents a value stored at a certain path in the tree.
type PathValue struct {
	Path  Path
	Value interface{}
}

// All returns all items stored in the tree under the given prefix.
func (t *Tree) All(prefix Path) []PathValue {
	var res []PathValue
	t.Walk(prefix, func(p Path, val interface{}) error {
		res = append(res, PathValue{p, val})
		return nil
	})
	return res
}

// Clone returns a clone of the tree.
func (t *Tree) Clone() *Tree {
	return &Tree{root: t.root.clone()}
}
