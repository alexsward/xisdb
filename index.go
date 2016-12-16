package xisdb

import (
	"strings"

	"github.com/alexsward/xisdb/indexes"
	"github.com/alexsward/xisdb/tree"
)

// IndexType -- an index can be added to either the Key or Value of an Item
type IndexType int

const (
	// KeyIndex will index on an item's Key
	KeyIndex IndexType = iota
	// ValueIndex will index on an item's Value
	ValueIndex
)

// indexMatcher is a fucntion that determines if an Item matches an index
type indexMatcher func(*Item) bool

func newIndexMatcher(it IndexType, matcher indexes.Matcher) indexMatcher {
	return func(item *Item) bool {
		str := item.Key
		if it == ValueIndex {
			str = item.Value
		}
		return matcher(str)
	}
}

type index struct {
	name  string
	match indexMatcher
	tree  tree.BTree
}

func (i *index) String() string {
	return i.name
}

var (
	// NaturalOrderKeyComparison -- string.Compare two Items by Key
	NaturalOrderKeyComparison = func(k1, k2 tree.Key) int {
		// TODO: these comparators need to be way better
		return strings.Compare(k1.(string), k1.(string))
	}
)

type indexNode struct {
	item *Item
}

func (in indexNode) Key() tree.Key {
	return in.item.Key
}

func (in indexNode) Value() interface{} {
	return in.item
}

func newIndex(name string, it IndexType, m indexes.Matcher, comp tree.Comparator) (*index, error) {
	tree, err := tree.NewTree(3, NaturalOrderKeyComparison)
	if err != nil {
		return nil, err
	}

	idx := &index{
		name:  name,
		match: newIndexMatcher(it, m),
		tree:  tree,
	}
	return idx, err
}

func (i *index) add(item *Item) {
	i.tree.Insert(&indexNode{item})
}

func (i *index) remove(item *Item) {
	i.tree.Remove(&indexNode{item})
}

func (i *index) iterate() <-chan Item {
	ch := make(chan Item)
	go func(c chan Item) {
		defer close(c)
		if i.tree.Size() == 0 {
			return
		}
		for item := range i.tree.IterateAll() {
			ch <- *(item.(*indexNode).item)
		}
	}(ch)
	return ch
}
