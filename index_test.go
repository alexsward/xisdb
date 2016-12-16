package xisdb

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/alexsward/xisdb/indexes"
)

func TestIndexName(t *testing.T) {
	fmt.Println("--- TestIndexName")
	i, _ := newIndex("name", KeyIndex, nil, nil)
	if i.String() != "name" {
		t.Errorf("Expected index String() 'name', got %s", i.String())
	}
}

func TestIndexMatchers(t *testing.T) {
	fmt.Println("--- TestIndexMatchers")
	tests := []struct {
		it      IndexType
		m       indexes.Matcher
		item    Item
		matches bool
	}{
		{KeyIndex, indexes.WildcardMatcher, Item{"key", "", nil}, true},
		{ValueIndex, indexes.WildcardMatcher, Item{"", "value", nil}, true},
		{KeyIndex, indexes.PrefixMatcher("k"), Item{"key", "", nil}, true},
		{KeyIndex, indexes.PrefixMatcher("key"), Item{"key", "", nil}, true},
		{ValueIndex, indexes.PrefixMatcher("v"), Item{"", "value", nil}, true},
		{ValueIndex, indexes.PrefixMatcher("value"), Item{"", "value", nil}, true},
		{KeyIndex, indexes.PrefixMatcher("e"), Item{"key", "", nil}, false},
		{ValueIndex, indexes.PrefixMatcher("a"), Item{"", "value", nil}, false},
	}
	for i, test := range tests {
		match := newIndexMatcher(test.it, test.m)
		if match(&test.item) != test.matches {
			t.Errorf("Test %d failed: expected match %t, got %t", i+1, test.matches, !test.matches)
		}
	}
}

func TestIndexIterateOrder(t *testing.T) {
	fmt.Println("--- TestIndexIterateOrder")
	tests := []struct {
		items, expected []string
		it              IndexType
		match           indexes.Matcher
		order           indexes.Order
	}{
		{[]string{}, []string{}, KeyIndex, indexes.WildcardMatcher, indexes.ASC},
		{[]string{"a", "b", "c", "d"}, []string{"a", "b", "c", "d"}, KeyIndex, indexes.WildcardMatcher, indexes.ASC},
		{[]string{"d", "c", "b", "a"}, []string{"a", "b", "c", "d"}, KeyIndex, indexes.WildcardMatcher, indexes.ASC},
	}
	for i, test := range tests {
		db := openTestDB()
		// TODO: move all this to a helper function?
		err := db.AddIndex("test-index", test.it, test.match, NaturalOrderKeyComparison)
		if err != nil {
			t.Errorf("Test %d failed: error creating index: '%s'", i+1, err)
			continue
		}
		index := db.indexes["test-index"]
		for j, raw := range test.items {
			item := createTestItemForIndex(test.it, raw, j)
			index.add(&item)
		}
		if index.tree.Size() != uint(len(test.expected)) {
			t.Errorf("Test %d failed: expeted index to have %d items, had %d", i+1, len(test.items), index.tree.Size())
		}
		err = assertIteration(t, index.iterate(), test.expected)
		if err != nil {
			t.Errorf("Test %d failed: got error: %s", i+1, err)
		}
	}
}

func assertIteration(t *testing.T, iterator <-chan Item, expected []string) error {
	tick := time.NewTicker(time.Millisecond * 20)
	defer tick.Stop()
	total := 0

	iterDone := make(chan error)
	var err error
	go func() {
		for item := range iterator {
			if item.Key != expected[total] {
				t.Errorf("Item failed at index:%d, got %s and expected %s", total, item.Key, expected[total])
				err = fmt.Errorf("Item failed at index:%d, got %s and expected %s", total, item.Key, expected[total])
				return
			}
			total++
		}
		iterDone <- nil
	}()
	select {
	case <-tick.C:
		return fmt.Errorf("Test timed out, iterated %d items, expected: %d", total, len(expected))
	case <-iterDone:
		if total != len(expected) {
			return fmt.Errorf("Iteration completed, Expected %d items, got %d", len(expected), total)
		}
		return err
	}
}

func createTestItemForIndex(it IndexType, data string, num int) Item {
	if it == ValueIndex {
		// TODO: not doing this yet
		return Item{}
	}
	return Item{data, strconv.Itoa(num), nil}
}
