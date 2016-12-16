package tree

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrInvalidDegree when an invalid degree is used to initialize a tree
	ErrInvalidDegree = errors.New("Degree must be >= 3")
	// ErrKeyNotFound when a desired key doesn't exist in the tree
	ErrKeyNotFound = errors.New("Key doesn't exist in BTree")
)

// BTree is the tree structure itself
type BTree interface {
	// Get finds all values with the provided key, or ErrKeyNotFound if none
	Get(key Key) ([]Node, error)
	// Insert inserts the given nodes into the tree, will return an error on first failure
	Insert(...Node) error
	// Remove removes the given nodes from the tree, will return an error on first failure
	Remove(...Node) error
	// IterateAll will return a channel over every Node in the BTree
	IterateAll() <-chan Node
	Iterate(start, end Key) <-chan Node

	Height() int
	Size() uint
}

// Comparator will compare two objects (keys) for equality
// Returns 0 if equal, -1 if less, 1 if greater
type Comparator func(Key, Key) int

type btree struct {
	degree int
	root   *btnode
	comp   Comparator
	less   func(Key, Key) bool
	size   uint

	list elements
}

// NewTree creates a tree of degree d using the supplied comparator
func NewTree(d int, c Comparator) (BTree, error) {
	if d < 3 {
		return nil, ErrInvalidDegree
	}

	t := &btree{
		degree: d,
		comp:   c,
		less: func(k1, k2 Key) bool {
			return c(k1, k2) < 0
		},
		root: nil,
		list: make(elements, 0),
	}
	return t, nil
}

func (bt *btree) Get(key Key) ([]Node, error) {
	for n := bt.root; n != nil; {
		idx, found := n.find(key)
		if n.isLeaf() {
			if !found {
				return nil, ErrKeyNotFound
			}
			return n.elements[idx].overflow, nil
		}

		if found {
			idx++
		}
		n = n.children[idx]
	}
	return nil, ErrKeyNotFound
}

func (bt *btree) Insert(nodes ...Node) error {
	for _, n := range nodes {
		if err := bt.insertNode(n); err != nil {
			return err
		}
		bt.size = bt.size + 1
	}
	return nil
}

// insert performs the actual heavy lifting of an insert, including splitting of nodes
func (bt *btree) insertNode(n Node) error {
	var inserted *element
	defer func() {
		bt.list.add(inserted, bt.less)
	}()

	if bt.root == nil {
		bt.root = newEmptyNode(bt, nil)
		inserted = bt.root.insertElement(n, 0)
		return nil
	}

	// insert at the end of the path in the proper position
	path, idx := bt.findPath(n.Key())
	end := path[len(path)-1]
	inserted = end.insertElement(n, idx)

	// then split the tree up if necessary
	return bt.split(end)
}

func (bt *btree) split(node *btnode) error {
	if node == nil || !node.shouldSplit() {
		return nil
	}

	// 1. A single median is chosen from among the leaf's elements and the new element.
	m := node.median()
	middle := newNode(bt, node.parent, nil, node.getElements(m, m+1))
	// 2. Values less than the median are put in the new left node
	left := newNode(bt, node.parent, nil, node.leftElements(m))
	// 3. Values greater than the median are put in the new right node, with the median acting as a separation value.
	right := newNode(bt, node.parent, nil, node.rightElements(m))

	if node == bt.root {
		// assign left and right to the middle node children from the split
		middle.children = []*btnode{left, right}
		assignNonLeafChildren(node, left, right, m)

		// reassign root to the middle with two children
		left.parent = middle
		right.parent = middle
		bt.root = middle
	} else {
		idx, _ := node.parent.find(middle.elements[0].Key())
		node.parent.insertElement(middle.elements[0].overflow[0], idx)
		assignNonLeafChildren(node, left, right, m)

		// remove the old child that's been split from the parent and add the new children
		node.parent.deleteChild(idx)
		node.parent.insertChild(left, idx)
		node.parent.insertChild(right, idx+1)
	}

	// recurse through the entire tree to split
	return bt.split(node.parent)
}

// assignNonLeafChildren properly sets children on non-leaf nodes during a split
func assignNonLeafChildren(node, left, right *btnode, m int) {
	if node.isLeaf() {
		return
	}
	// copy the split node children down to left and right
	left.children = append([]*btnode{}, node.children[:m+1]...)
	right.children = append([]*btnode{}, node.children[m+1:]...)

	// properly assign parent pointers
	left.assignParent()
	right.assignParent()
}

// find returns a path to where a Node should be inserted into the tree
// the returned path is in root -> leaf order for where the Node n belongs
// the returned integer is the position in the last item of the path to insertt
func (bt *btree) findPath(key Key) ([]*btnode, int) {
	var path []*btnode
	idx := 0
	for node := bt.root; node != nil; node = node.children[idx] {
		path = append(path, node)
		i, found := node.find(key)
		idx = transformPathIndex(node, found, i)
		if len(node.children) == 0 {
			break
		}
	}
	return path, idx
}

func (bt *btree) Remove(nodes ...Node) error {
	return nil
}

func (bt *btree) left() Node {
	child := func(n []*btnode) int {
		return 0
	}
	element := func(e elements) int {
		return 0
	}
	return getEndNode(bt, child, element)
}

func (bt *btree) right() Node {
	child := func(n []*btnode) int {
		return len(n) - 1
	}
	element := func(e elements) int {
		return len(e) - 1
	}
	return getEndNode(bt, child, element)
}

func (bt *btree) Height() int {
	if bt.root == nil {
		return 0
	}
	return bt.root.height()
}

func (bt *btree) Size() uint {
	return bt.size
}

// TODO: DELETE ALL THIS BELOW

func printTree(t *btree) {
	printNodes(0, t.root)
}

func printNodes(level int, nodes ...*btnode) {
	if len(nodes) == 0 {
		return
	}

	var children []*btnode
	fmt.Printf("Level %d: ", level)
	for _, n := range nodes {
		children = append(children, n.children...)
		fmt.Print("[ ")
		fmt.Printf("%s ", getElementsString(n))
		fmt.Print("]")
	}
	fmt.Print("\n")
	printNodes(level+1, children...)
}

func printNode(n *btnode) {
	fmt.Printf("elements:[%s]\n", getElementsString(n))
}

func printElements(n *btnode) {
	fmt.Printf("elements: %s\n", getElementsString(n))
}

func getElementsString(n *btnode) string {
	var s []string
	for _, e := range n.elements {
		s = append(s, fmt.Sprintf("%d", e.Key()))
	}
	return strings.Join(s, " ")
}
