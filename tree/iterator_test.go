package tree

import (
	"fmt"
	"testing"
)

// TestIterateEmpty
func TestIterateEmpty(t *testing.T) {
	fmt.Println("-- TestIterateEmpty")
	tree := getTestTree(testDegree)
	assertFullIteration(t, tree, 0, 10)
}

// TestIterateSingleLevelOrderedInsert
func TestIterateSingleLevelOrderedInsert(t *testing.T) {
	fmt.Println("-- TestIterateSingleLevelOrderedInsert")
	tree := getTestTreeForIteration(5, 1)
	assertFullIteration(t, tree, 4, 10)
}

// TestIterateTwoLevelOrderedInsert
func TestIterateTwoLevelOrderedInsert(t *testing.T) {
	fmt.Println("-- TestIterateTwoLevelOrderedInsert")
	tree := getTestTreeForIteration(23, 1)
	assertFullIteration(t, tree, 22, 30)
}

// TestIterateTwoLevelOrderedInsertOverflow -- overflows the tree from TestIterateTwoLevelOrderedInsert
func TestIterateTwoLevelOrderedInsertOverflow(t *testing.T) {
	fmt.Println("-- TestIterateTwoLevelOrderedInsertOverflow")
	tree := getTestTreeForIteration(23, 3)
	assertFullIteration(t, tree, 66, 90)
}

// TestIterateRangeSingleLevel
func TestIterateRangeSingleLevel(t *testing.T) {
	fmt.Println("-- TestIterateRangeSingleLevel")
	tree := getTestTreeForIteration(5, 1)
	assertRangeIteration(t, tree, 2, 10, 2, 3)
}

// TestIterateRangeTwoLevel
func TestIterateRangeTwoLevel(t *testing.T) {
	fmt.Println("-- TestIterateRangeTwoLevel")
	tree := getTestTreeForIteration(23, 1)
	assertRangeIteration(t, tree, 8, 10, 8, 15)
}

func getTestTreeForIteration(upper, each int) *btree {
	tree := getTestTree(testDegree)
	for i := 1; i < upper; i++ {
		for j := 0; j < each; j++ {
			tree.Insert(testNode{i})
		}
	}
	return tree
}

func assertFullIteration(t *testing.T, tree *btree, size, limit int) []Node {
	var received []Node
	i := 0
	ch := tree.IterateAll()
	for n := range ch {
		i++
		received = append(received, n)
		if i >= limit {
			break
		}
	}
	assertIteratedItems(t, received, size)
	return received
}

func assertRangeIteration(t *testing.T, tree *btree, size, limit int, start, end Key) []Node {
	var received []Node
	i := 0
	ch := tree.Iterate(start, end)
	for n := range ch {
		i++
		received = append(received, n)
		if i >= limit {
			break
		}
	}
	assertIteratedItems(t, received, size)
	return received
}

func assertIteratedItems(t *testing.T, received []Node, size int) {
	if len(received) != size {
		t.Errorf("Expected %d items, instead got %d", size, len(received))
	}
}
