package gostream

// OutChan sends stream items to the specified channel and closes it when done.
func (s Stream[T]) OutChan(ch chan<- T) {
	next := s.Iterate()

	for item, ok := next(); ok; item, ok = next() {
		ch <- item
	}

	close(ch)
}
