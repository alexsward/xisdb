package tree

// TODO: maybe this shouldn't be its own file?

func getEndNode(bt *btree, childIndex func([]*btnode) int, elemIndex func(elements) int) Node {
	if bt.root == nil {
		return nil
	}

	end := bt.root
	for len(end.children) > 0 {
		end = end.children[childIndex(end.children)]
	}
	return end.elements[elemIndex(end.elements)].overflow[0]
}

func transformPathIndex(n *btnode, found bool, idx int) int {
	if found && len(n.children) > 0 {
		return idx + 1
	}

	if idx != 0 && len(n.children) == idx {
		return idx - 1
	}

	return idx
}
