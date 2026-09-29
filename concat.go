package gostream

// Concat concatenates two streams.
//
// The Concat method differs from the Union method because the Concat method
// returns all the original elements in the input sequences. The Union method
// returns only unique elements.
func (s Stream[T]) Concat(q2 Stream[T]) Stream[T] {
	return Stream[T]{
		Iterate: func() Iterator[T] {
			next := s.Iterate()
			next2 := q2.Iterate()
			use1 := true

			return func() (item T, ok bool) {
				if use1 {
					item, ok = next()
					if ok {
						return
					}
					use1 = false
				}
				return next2()
			}
		},
	}
}

// Append inserts an item to the end of a collection, so it becomes the last
// item.
func (s Stream[T]) Append(item T) Stream[T] {
	return Stream[T]{
		Iterate: func() Iterator[T] {
			next := s.Iterate()
			appended := false

			return func() (T, bool) {
				i, ok := next()
				if ok {
					return i, ok
				}

				if !appended {
					appended = true
					return item, true
				}

				var zero T
				return zero, false
			}
		},
	}
}
