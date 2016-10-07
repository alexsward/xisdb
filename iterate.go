package xisdb

// Each -- Returns a channel that can be iterated upon
func (tx *Tx) Each() chan Item {
	ch := make(chan Item)
	go func() {
		for _, item := range tx.db.data {
			ch <- item
		}
		close(ch)
	}()

	return ch
}
