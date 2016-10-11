package xisdb

import "strings"

type subscription struct {
	prefix   string
	channels []chan Item
}

func (s *subscription) add(capacity int) chan Item {
	ch := make(chan Item, capacity)
	s.channels = append(s.channels, ch)
	return ch
}

// Subscribe to all changes for a specific prefix
func (db *DB) Subscribe(prefix string, capacity int) (chan Item, error) {
	if _, exists := db.subscriptions[prefix]; !exists {
		db.subscriptions[prefix] = subscription{
			prefix:   prefix,
			channels: make([]chan Item, 0),
		}
	}

	s := db.subscriptions[prefix]
	ch := s.add(capacity)
	db.subscriptions[prefix] = s
	return ch, nil
}

// Unsubscribe a particular prefix/channel combination
// If ch is nil, will close and unsubscribe all channels for that pattern
func (db *DB) Unsubscribe(prefix string, ch chan Item) error {
	if _, exists := db.subscriptions[prefix]; !exists {
		return nil
	}

	for _, c := range db.subscriptions[prefix].channels {
		if ch == nil || c == ch {
			close(c)
		}
	}

	if len(db.subscriptions[prefix].channels) == 0 {
		delete(db.subscriptions, prefix)
	}

	return nil
}

func (db *DB) publish(items ...Item) error {
	for _, item := range items {
		for _, sub := range db.subscriptions {
			if strings.HasPrefix(item.Key, sub.prefix) {
				for _, ch := range sub.channels {
					ch <- item
				}
			}
		}
	}

	return nil
}
