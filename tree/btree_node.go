package tree

import "sort"

// Node is anything put into the tree
type Node interface {
	Key() Key
	Value() interface{}
}

// Key is a key for a node
type Key interface{}

// element contains a Key and all the actual Node items for that key (enables duplicates easily)
type element struct {
	key      Key
	overflow []Node
}

func (e *element) Key() Key {
	return e.key
}

type elements []*element

// add places a new pointer into this elements slice
func (es *elements) add(e *element, less func(k1, k2 Key) bool) bool {
	idx, exists := es.indexOf(e.Key(), less)
	if exists {
		return false
	}

	*es = append(*es, nil)
	if idx < len(*es) {
		copy((*es)[idx+1:], (*es)[idx:])
	}
	(*es)[idx] = e
	return true
}

// indexOf returns the index where to put the element, and if it exists
func (es *elements) indexOf(key Key, less func(k1, k2 Key) bool) (int, bool) {
	idx := sort.Search(len(*es), func(i int) bool {
		return less(key, (*es)[i].Key())
	})

	if 0 < idx && !less((*es)[idx-1].Key(), key) {
		return idx - 1, true
	}
	return idx, false
}

type btnode struct {
	tree     *btree
	parent   *btnode
	children []*btnode
	elements elements
}

func newEmptyNode(t *btree, p *btnode) *btnode {
	return newNode(t, p, make([]*btnode, 0), make(elements, 0))
}

func newNode(t *btree, p *btnode, c []*btnode, e elements) *btnode {
	if c == nil {
		c = make([]*btnode, 0)
	}
	if e == nil {
		e = make(elements, 0)
	}

	return &btnode{
		tree:     t,
		parent:   p,
		children: c,
		elements: e,
	}
}

func (bn *btnode) shouldSplit() bool {
	return len(bn.elements) > bn.maximumSize()
}

func (bn *btnode) isLeaf() bool {
	return len(bn.children) == 0
}

func (bn *btnode) maximumSize() int {
	return (2 * bn.tree.degree) - 1
}

func (bn *btnode) median() int {
	return bn.tree.degree
}

func (bn *btnode) getElements(s, e int) elements {
	return append(elements{}, bn.elements[s:e]...)
}

func (bn *btnode) leftElements(m int) elements {
	return bn.getElements(0, m)
}

func (bn *btnode) rightElements(m int) elements {
	return bn.getElements(m, len(bn.elements))
}

func (bn *btnode) height() int {
	height := 0
	for n := bn; n != nil; n = n.children[0] {
		height++
		if len(n.children) == 0 {
			break
		}
	}
	return height
}

func (bn *btnode) assignParent() {
	for _, child := range bn.children {
		child.parent = bn
	}
}

func (bn *btnode) insertElement(node Node, i int) *element {
	idx, contains := bn.find(node.Key())
	if contains {
		bn.elements[idx].overflow = append(bn.elements[idx].overflow, node)
		return bn.elements[idx]
	}

	bn.elements = append(bn.elements, nil)
	if i < len(bn.elements) {
		// need to actually move elements around
		copy(bn.elements[i+1:], bn.elements[i:])
	}

	bn.elements[i] = &element{
		key:      node.Key(),
		overflow: []Node{node},
	}
	return bn.elements[i]
}

func (bn *btnode) insertChild(n *btnode, i int) {
	bn.children = append(bn.children, nil)
	if i < len(bn.children) {
		copy(bn.children[i+1:], bn.children[i:])
	}
	bn.children[i] = n
}

func (bn *btnode) deleteChild(i int) {
	copy(bn.children[i:], bn.children[i+1:])
	bn.children[len(bn.children)-1] = nil
	bn.children = bn.children[:len(bn.children)-1]
}

// find determines what position this element would be placed at, or if it's found
func (bn *btnode) find(key Key) (int, bool) {
	idx := sort.Search(len(bn.elements), func(i int) bool {
		return bn.tree.less(key, bn.elements[i].Key())
	})

	if 0 < idx && !bn.tree.less(bn.elements[idx-1].Key(), key) {
		// it's less than an item, but not less than the previous, means it's already here
		return idx - 1, true
	}
	return idx, false
}
