package smap

import "sync"

// TODO: possibly shard this and lock per shard?
type smap struct {
	entries map[string]string
	mutex   sync.RWMutex
}

func new() *smap {
	return &smap{
		entries: make(map[string]string),
	}
}

func (m *smap) get(key string) (string, bool) {
	m.lock(false)
	defer m.unlock(false)
	value, exists := m.entries[key]
	return value, exists
}

func (m *smap) delete(key string) bool {
	m.lock(true)
	defer m.unlock(true)
	contains := m.contains(key)
	delete(m.entries, key)
	return contains
}

func (m *smap) put(key, value string) bool {
	m.lock(true)
	defer m.unlock(true)

	updated := m.contains(key)
	m.entries[key] = value
	return updated
}

func (m *smap) replace(key, value string) bool {
	if m.contains(key) {
		return m.put(key, value)
	}
	return false
}

// alter will apply fn to the key's value, if it exists
func (m *smap) alter(key string, fn func(string) string) bool {
	v, contains := m.get(key)
	if !contains {
		return contains
	}
	m.put(key, fn(v))
	return contains
}

func (m *smap) contains(key string) bool {
	_, contains := m.entries[key]
	return contains
}

func (m *smap) containsValue(search string) bool {
	m.lock(false)
	defer m.unlock(false)

	for _, value := range m.entries {
		if search == value {
			return true
		}
	}
	return false
}

func (m *smap) size() int {
	return len(m.entries)
}

func (m *smap) isEmpty() bool {
	return len(m.entries) == 0
}

func (m *smap) forEach(fn func(string, string) bool) {
	for key, value := range m.entries {
		if fn(key, value) {
			return
		}
	}
}

// merge will combine m2 into the smap, and return a slice of keys that were updated
func (m *smap) merge(m2 *smap) []string {
	var additions []string
	for key, value := range m2.entries {
		if m.contains(key) {
			additions = append(additions, key)
		}
		m.put(key, value)
	}
	return additions
}

func (m *smap) transform(fn func(string) string) {
	m.forEach(func(key, value string) bool {
		m.put(key, fn(value))
		return false
	})
}

func (m *smap) lock(write bool) {
	if write {
		m.mutex.Lock()
		return
	}

	m.mutex.RLock()
}

func (m *smap) unlock(write bool) {
	if write {
		m.mutex.Unlock()
		return
	}

	m.mutex.RUnlock()
}
