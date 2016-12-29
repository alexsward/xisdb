package smap

import (
	"fmt"
	"testing"
)

func TestSmapSize(t *testing.T) {
	fmt.Println("-- TestSmapSize")
	m := new()
	assertSmapSize(t, m, 0)
	if !m.isEmpty() {
		t.Error("Map size 0 assertion passed, but isEmpty didnt?")
	}
	m.put("key", "val")
	assertSmapSize(t, m, 1)
}

func TestSmapContainsValue(t *testing.T) {
	fmt.Println("-- TestSmapContainsValue")
	m := new()
	m.put("1", "a")
	contains := m.containsValue("a")
	if !contains {
		t.Error("Expected map to contain value 'a', it didn't")
	}
	notContains := m.containsValue("whatever")
	if notContains {
		t.Error("Did not expect to contain value that map shouldn't")
	}
}

func TestSmapDelete(t *testing.T) {
	fmt.Println("-- TestSmapDelete")
	m := new()
	m.put("1", "a")
	assertSmapValue(t, m, "1", "a")
	m.delete("1")
	val, got := m.get("1")
	if got {
		t.Errorf("Did not expect a value for key %s, but got %s", "1", val)
	}
}

func TestSmapReplace(t *testing.T) {
	fmt.Println("-- TestSmapReplace")
	m := new()
	replaced := m.replace("1", "a")
	if replaced {
		t.Error("Did not expect to be able to replace a non-existent key")
	}
	m.put("1", "a")
	m.replace("1", "b")
	assertSmapValue(t, m, "1", "b")
}

func TestSmapMerge(t *testing.T) {
	fmt.Println("-- TestSmapMerge")
	m1 := new()
	m1.put("1", "a")
	m1.put("2", "b")
	m1.put("3", "c")
	m2 := new()
	m2.put("3", "d")
	m2.put("4", "e")
	m2.put("5", "f")
	m1.merge(m2)
	assertSmapValue(t, m1, "1", "a")
	assertSmapValue(t, m1, "2", "b")
	assertSmapValue(t, m1, "3", "d")
	assertSmapValue(t, m1, "4", "e")
	assertSmapValue(t, m1, "5", "f")
}

func TestSmapTransform(t *testing.T) {
	fmt.Println("-- TestSmapTransform")
	m := new()
	m.put("1", "a")
	m.put("2", "a")
	m.put("3", "a")
	m.transform(func(arg1 string) string {
		return "b"
	})
	assertSmapValue(t, m, "1", "b")
	assertSmapValue(t, m, "2", "b")
	assertSmapValue(t, m, "3", "b")
	i := 1
	m.forEach(func(k, v string) bool {
		if i == 3 {
			return true
		}
		i++
		m.put(k, "c")
		return false
	})
	assertSmapValue(t, m, "1", "c")
	assertSmapValue(t, m, "2", "c")
	assertSmapValue(t, m, "3", "b")
}

func TestSmapAlter(t *testing.T) {
	fmt.Println("-- TestSmapAlter")
	m := new()
	alter := func(x string) string {
		return "b"
	}
	altered1 := m.alter("1", alter)
	if altered1 {
		t.Error("Didn't expect to be able to alter key")
	}
	m.put("1", "a")
	altered2 := m.alter("1", alter)
	if !altered2 {
		t.Error("Expected to be able to alter key")
	}
}

func assertSmapSize(t *testing.T, m *smap, expected int) {
	if m.size() != expected {
		t.Errorf("Expected map size to be %d, instead got %d", expected, m.size())
	}
}

func assertSmapValue(t *testing.T, m *smap, key, value string) {
	got, found := m.get(key)
	if !found {
		t.Errorf("Expected to find key %s, didn't", key)
	}
	if got != value {
		t.Errorf("Expected key:%s to contain value %s, instead got %s", key, value, got)
	}
}
