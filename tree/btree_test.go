package tree

import (
	"fmt"
	"testing"
)

var (
	testDegree = 3
)

type depthTest struct {
	child, children, elems int
	values                 []interface{}
}

// TestNodeElementInsert tests simply adding to an array of elements
func TestNodeElementInsert(t *testing.T) {
	fmt.Println("-- TestNodeElementInsert")

	n := &btnode{
		elements: []*element{getElements(testNode{5})},
		tree:     getTestTree(3),
	}
	assertNodeValue(t, n, 0, 5)
	assertElementLength(t, n, 1)
	n.insertElement(testNode{8}, 1)
	assertElementLength(t, n, 2)
	assertNodeValue(t, n, 0, 5)
	assertNodeValue(t, n, 1, 8)
	n.insertElement(testNode{4}, 0)
	assertElementLength(t, n, 3)
	assertNodeValue(t, n, 0, 4)
	assertNodeValue(t, n, 1, 5)
	assertNodeValue(t, n, 2, 8)
}

// TestNodeFindKey validates the correct positioning of where to insert elements in an elements array
func TestNodeFindKey(t *testing.T) {
	fmt.Println("-- TestNodeFindKey")
	n := &btnode{
		elements: []*element{
			getElements(testNode{3}),
			getElements(testNode{4}),
			getElements(testNode{5}),
			getElements(testNode{7}),
			getElements(testNode{30}),
		},
		tree: getTestTree(3),
	}

	var tests = []struct {
		key, expected int
		found         bool
	}{
		{0, 0, false},
		{1, 0, false},
		{2, 0, false},
		{3, 0, true},
		{4, 1, true},
		{7, 3, true},
		{30, 4, true},
		{31, 5, false},
		{500, len(n.elements), false},
	}
	for i, test := range tests {
		idx, found := n.find(test.key)
		if idx != test.expected || found != test.found {
			t.Errorf("Test %d failed, expected index:%d returned_index:%d, expected found:%t returned_found:%t", i+1, test.expected, idx, test.found, found)
		}
	}
}

// TestNodeChildInsert tests inserting an internal node
func TestNodeChildInsert(t *testing.T) {
	fmt.Println("-- TestNodeChildInsert")
	//
	// n := &btnode{
	// 	children: make([]*btnode, 0),
	// }
	// n.insertChild(&btnode{elements: []Node{testNode{7}}}, 0)
	// n.insertChild(&btnode{elements: []Node{testNode{3}}}, 0)
}

// TestInsertRoot tests root insertion
func TestInsertRoot(t *testing.T) {
	fmt.Println("-- TestInsertRoot")

	btree := getTestTree(testDegree)
	btree.Insert(testNode{5})
	if btree.root.parent != nil {
		t.Error("Expected root parent to be nil")
	}
	assertElementLength(t, btree.root, 1)
	assertChildrenLength(t, btree.root, 0)
}

// TestRootBeforeSplit tests filling the root before it splits
func TestRootBeforeSplit(t *testing.T) {
	fmt.Println("-- TestRootBeforeSplit")

	btree := getTestTree(testDegree)
	btree.Insert(testNode{5})
	btree.Insert(testNode{8})
	btree.Insert(testNode{4})
	btree.Insert(testNode{7})
	btree.Insert(testNode{10})

	assertChildrenLength(t, btree.root, 0)
	assertElementLength(t, btree.root, 5)
}

// TestFirstRootSplit tests filling the root before it splits
func TestFirstRootSplit(t *testing.T) {
	fmt.Println("-- TestFirstRootSplit")

	btree := getTestTree(testDegree)
	btree.Insert(testNode{5})
	btree.Insert(testNode{8})
	btree.Insert(testNode{4})
	btree.Insert(testNode{7})
	btree.Insert(testNode{10})
	btree.Insert(testNode{6})

	assertChildrenLength(t, btree.root, 2)
	assertElementLength(t, btree.root, 1)
}

// TestLeftBranchSplit tests the first split on the left branch
// Should end up with a tree in the form:
// [5, 7]
// [1,3,4] [5,6,6] [7,8,10]
func TestLeftBranchSplit(t *testing.T) {
	fmt.Println("-- TestLeftBranchSplit")

	btree := getTestTree(testDegree)
	btree.Insert(testNode{5})
	btree.Insert(testNode{8})
	btree.Insert(testNode{4})
	btree.Insert(testNode{7})
	btree.Insert(testNode{10})
	btree.Insert(testNode{6})
	btree.Insert(testNode{2})
	btree.Insert(testNode{3})
	btree.Insert(testNode{1})

	assertChildrenLength(t, btree.root, 3)
	assertElementLength(t, btree.root, 2)

	assertElementEquality(t, btree.root.children[0], []interface{}{1, 2, 3})
	assertElementEquality(t, btree.root.children[1], []interface{}{4, 5, 6})
	assertElementEquality(t, btree.root.children[2], []interface{}{7, 8, 10})
}

// TestOverflowInsertion1 tests adding same-valued keys over and over again into 1 node
func TestOverflowInsertion1(t *testing.T) {
	fmt.Println("-- TestOverflowInsertion")

	btree := getTestTree(testDegree)
	btree.Insert(testNode{1})
	assertElementLength(t, btree.root, 1)
	btree.Insert(testNode{1})
	btree.Insert(testNode{1})
	assertElementLength(t, btree.root, 1)

}

// TestLargeHeight1 creates a large 1 height tree and tests it
func TestLargeHeight1(t *testing.T) {
	fmt.Println("-- TestLargeHeight1")

	btree := getTestTree(testDegree)
	btree.Insert(testNode{5})
	btree.Insert(testNode{8})
	btree.Insert(testNode{4})
	btree.Insert(testNode{7})
	btree.Insert(testNode{10})
	btree.Insert(testNode{6})
	btree.Insert(testNode{6})
	btree.Insert(testNode{3})
	btree.Insert(testNode{1})
	btree.Insert(testNode{8})
	btree.Insert(testNode{2})
	btree.Insert(testNode{11})
	btree.Insert(testNode{12})
	btree.Insert(testNode{4})
	btree.Insert(testNode{2})
	btree.Insert(testNode{100})
	btree.Insert(testNode{8})
	btree.Insert(testNode{46})
	btree.Insert(testNode{26})
	btree.Insert(testNode{12})
	btree.Insert(testNode{9})
	btree.Insert(testNode{17})

	assertChildrenLength(t, btree.root, 5)
	assertElementLength(t, btree.root, 4)

	tests := []depthTest{
		{0, 0, 3, []interface{}{1, 2, 3}},
		{1, 0, 3, []interface{}{4, 5, 6}},
		{2, 0, 4, []interface{}{7, 8, 9, 10}},
		{3, 0, 3, []interface{}{11, 12, 17}},
		{4, 0, 3, []interface{}{26, 46, 100}},
	}
	assertDepthTests(t, btree.root, tests)
	assertElementOverflowLength(t, btree.root.children[1], 2, 2)
	assertElementOverflowLength(t, btree.root.children[2], 1, 3)
}

func TestHeight2Splits(t *testing.T) {
	btree := getTestTree(testDegree)
	for i := 1; i < 22; i++ {
		btree.Insert(testNode{i})
	}
	assertHeight(t, btree.root, 3)
}

// TestGetEmptyTree makes sure you can't Get an empty tree
func TestGetEmptyTree(t *testing.T) {
	fmt.Println("-- TestGetEmptyTree")
	tree := &btree{}
	if _, err := tree.Get(1); err != ErrKeyNotFound {
		t.Errorf("Expected %s for empty tree, but got %s", ErrKeyNotFound, err)
	}
}

// TestGetSingleItem tests retrieving an item that only matches once in a tree
func TestGetSingleItemHeight2(t *testing.T) {
	fmt.Println("-- TestGetSingleItemHeight2")
	tree := getTestGetTree(2)
	nodes, err := tree.Get(8)
	assertGet(t, 1, 8, nodes, err)
}

// TestGetSingleItemHeight3 retrieves a single item from a height=3 tree
func TestGetSingleItemHeight3(t *testing.T) {
	fmt.Println("-- TestGetSingleItemHeight3")
	tree := getTestGetTree(3)
	nodes, err := tree.Get(8)
	assertGet(t, 1, 8, nodes, err)
}

// TestGetMissingItemHeight2 makes sure a missing item isn't returned on a tree of height = 2
func TestGetMissingItemHeight2(t *testing.T) {
	fmt.Println("-- TestGetMissingItemHeight2")
	tree := getTestGetTree(2)
	if nodes, err := tree.Get(12312932138); err != ErrKeyNotFound {
		t.Errorf("Expected %s because this key shouldn't exist, instead got: %s %s", ErrKeyNotFound, nodes, err)
	}
}

// TestGetMultiItemHeight2 ensures multi-item returns work correctly on a tree of height=2
func TestGetMultiItemHeight2(t *testing.T) {
	fmt.Println("-- TestGetMultiItemHeight2")
	tree := getTestGetTree(2)
	nodes, err := tree.Get(4)
	assertGet(t, 2, 4, nodes, err)
}

// TestGetMultiItemHeight3 ensures multi-item returns work correctly on tree of height=3
func TestGetMultiItemHeight3(t *testing.T) {
	fmt.Println("-- TestGetMultiItemHeight3")
	tree := getTestGetTree(3)
	nodes, err := tree.Get(4)
	assertGet(t, 3, 4, nodes, err)
}

// TestGetLeftItem makes sure Left() returns the left-most node of a tree
func TestGetLeftRightItems(t *testing.T) {
	fmt.Println("-- TestGetLeftRightItems")
	tree2 := getTestGetTree(2)
	left2 := tree2.left().(testNode)
	right2 := tree2.right().(testNode)
	if left2.v != 1 {
		t.Errorf("Expected left node of height 2 to be 1, got %d", left2.v)
	}
	if right2.v != 11 {
		t.Errorf("Expected right node of height 3 to be 11, got %d", right2.v)
	}
	tree3 := getTestGetTree(3)
	left3 := tree3.left().(testNode)
	right3 := tree3.right().(testNode)
	if left3.v != 1 {
		t.Errorf("Expected left node of height 3 to be 1, got %d", left3.v)
	}
	if right3.v != 21 {
		t.Errorf("Expected right node of height 3 to be 21, got %d", right3.v)
	}
}

// TestTreeHeight validates tree height after inserts
func TestTreeHeight(t *testing.T) {
	fmt.Println("-- TestTreeHeight")

	btree := getTestTree(testDegree)
	if btree.Height() != 0 {
		t.Errorf("Expected tree height of 0, got %d instead", btree.Height())
	}
	btree.Insert(testNode{5})
	if btree.Height() != 1 {
		t.Errorf("Expected tree height of 0, got %d instead", btree.Height())
	}
}

// TestInvalidDegree makes sure you can't create too small a BTree
func TestInvalidDegree(t *testing.T) {
	fmt.Println("-- TestInvalidDegree")
	if _, err := NewTree(2, nil); err != ErrInvalidDegree {
		t.Errorf("Expected a %s because of invalid degree, got %s", ErrInvalidDegree, err)
	}
}

// TestTreeSize verifies counting elements in the tree is correct
func TestTreeSize(t *testing.T) {
	fmt.Println("-- TestTreeSize")
	btree := getTestTree(testDegree)
	for i := 1; i < 22; i++ {
		btree.Insert(testNode{i})
	}
	if btree.Size() != 21 {
		t.Errorf("Expected tree size %d, instead it was %d", 21, btree.Size())
	}
}

// TestElementsFind tests the elements struct
func TestElementsFind(t *testing.T) {
	fmt.Println("-- TestElementsFind")
	less := func(k1, k2 Key) bool {
		return k1.(int) < k2.(int)
	}
	tester := func(name string, es elements) {
		var cases = []struct {
			search, index int
			found         bool
		}{
			{2, 0, true}, {4, 1, true}, {6, 2, true}, {8, 3, true}, {10, 4, true},
			{12, 5, true}, {14, 6, true}, {16, 7, true}, {18, 8, true}, {20, 9, true},
			{22, 10, false}, {7, 3, false}, {1, 0, false},
		}
		str := "%s Case: %d -- Expected: element{%d}, index:%d, found:%t -- instead got index:%d, found:%t"
		for i, test := range cases {
			idx, found := es.indexOf((&element{test.search, nil}).Key(), less)
			if idx != test.index || found != test.found {
				t.Errorf(str, name, i+1, test.search, test.index, test.found, idx, found)
			}
		}
	}
	appended := make(elements, 0)
	for i := 1; i < 11; i++ {
		appended = append(appended, &element{i * 2, nil})
	}
	tester("BuiltWithAppend", appended)
	added := make(elements, 0)
	for i := 1; i < 11; i++ {
		added.add(&element{i * 2, nil}, less)
	}
	tester("BuiltWithAdding", added)
}

func intComparator(i1, i2 Key) int {
	item1 := i1.(int)
	item2 := i2.(int)

	if item1 < item2 {
		return -1
	} else if item1 == item2 {
		return 0
	} else {
		return 1
	}
}

type testNode struct {
	v int
}

func (tn testNode) Key() Key {
	return tn.v
}

func (tn testNode) Value() interface{} {
	return tn.v
}

func (tn testNode) String() string {
	return fmt.Sprintf("%d", tn.v)
}

func getElements(n ...Node) *element {
	return &element{
		key:      n[0].Key(),
		overflow: n,
	}
}

func getTestTree(degree int) *btree {
	tree, err := NewTree(degree, intComparator)
	if err != nil {
		fmt.Printf("Error creating tree? %s\n", err)
		return nil
	}

	return tree.(*btree)
}

func getTestGetTree(height int) *btree {
	tree := getTestTree(3)
	if height >= 2 {
		tree.Insert(testNode{4}, testNode{4}, testNode{8}, testNode{2}, testNode{7}, testNode{11},
			testNode{1})
	}
	if height >= 3 {
		tree.Insert(testNode{4}, testNode{3}, testNode{5}, testNode{6}, testNode{9}, testNode{10}, testNode{12}, testNode{13})
		tree.Insert(testNode{14}, testNode{15}, testNode{16}, testNode{17}, testNode{18}, testNode{19}, testNode{20}, testNode{21})
	}
	return tree
}

func assertNodeValue(t *testing.T, n *btnode, position int, expected interface{}) {
	if intComparator(n.elements[position].Key(), expected) != 0 {
		t.Errorf("Expected element at position %d to be equal to %s", position, expected)
	}
}

func assertChildrenLength(t *testing.T, n *btnode, length int) {
	if len(n.children) != length {
		t.Errorf("Expected node to have %d children, got %d instead", length, len(n.children))
	}
}

func assertElementLength(t *testing.T, n *btnode, length int) {
	if len(n.elements) != length {
		t.Errorf("Expected node to have %d elements, got %d instead", length, len(n.elements))
	}
}

func assertElementOverflowLength(t *testing.T, n *btnode, idx, total int) {
	nodes := len(n.elements[idx].overflow)
	if nodes != total {
		t.Errorf("Expected element at index %d to have %d Nodes, got %d", idx, total, nodes)
	}
}

func assertElementEquality(t *testing.T, n *btnode, expected []interface{}) {
	if len(n.elements) != len(expected) {
		t.Errorf("Expected the elements len to be %d, it was %d", len(expected), len(n.elements))
	}
	for i := range expected {
		if intComparator(n.elements[i].Key(), expected[i]) != 0 {
			t.Errorf("Item at index %d expected:%s got:%s", i, expected[i], n.elements[i].Key())
		}
	}
}

func assertHeight(t *testing.T, n *btnode, expected int) {
	if n.height() != expected {
		t.Errorf("Expected height at node to be %d, but it was %d", expected, n.height())
	}
}

func assertChildrenElementsSizing(t *testing.T, n *btnode) {
	internal := len(n.children) == len(n.elements)+1    // internal node
	leaf := len(n.children) == 0 && len(n.elements) > 0 //leaf node
	if !internal && !leaf {
		t.Errorf("incorrect sizing len(children)=%d len(elements)=%d", len(n.children), len(n.elements))
	}
}

func assertGet(t *testing.T, length, key int, nodes []Node, err error) {
	if err != nil {
		t.Errorf("Shouldn't have gotten an error, got: %s", err)
		return
	}
	if len(nodes) != length {
		t.Errorf("Expected number of items found to be %d, got: %d", length, len(nodes))
		return
	}
	for i, node := range nodes {
		if node.Key() != key {
			t.Errorf("Item returned at index %d was not %d, instead found %s", i, key, node.Key())
			return
		}
	}
}

func assertDepthTests(t *testing.T, n *btnode, cases []depthTest) {
	for _, test := range cases {
		assertChildrenLength(t, n.children[test.child], test.children)
		assertElementLength(t, n.children[test.child], test.elems)
		assertElementEquality(t, n.children[test.child], test.values)
		assertChildrenElementsSizing(t, n.children[test.child])
	}
}
