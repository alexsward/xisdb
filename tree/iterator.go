package tree

func (bt *btree) IterateAll() <-chan Node {
	ch := make(chan Node)
	go func(c chan Node) {
		defer close(c)
		bt.iterate(c, bt.Size(), nil, nil)
	}(ch)
	return ch
}

func (bt *btree) Iterate(start, end Key) <-chan Node {
	ch := make(chan Node)
	go func(c chan Node) {
		defer close(c)
		bt.iterate(c, bt.Size(), start, end)
	}(ch)
	return ch
}

func (bt *btree) iterate(ch chan Node, max uint, start, end Key) {
	if bt.root == nil {
		return
	}

	if start == nil {
		start = bt.left().Key()
	}
	if end == nil {
		end = bt.right().Key()
	}

	total := uint(0)
	i, _ := bt.list.indexOf(start, bt.less)
	for ; i < len(bt.list) && total < max; i++ {
		if bt.comp(bt.list[i].Key(), end) == 1 {
			break
		}

		for j := 0; j < len(bt.list[i].overflow) && total < max; j++ {
			ch <- bt.list[i].overflow[j]
			total++
		}
	}
}
