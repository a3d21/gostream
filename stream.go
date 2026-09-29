package gostream

// Map applies a mapping function to each element, producing a new Stream.
func (s Stream[T]) Map[U any](mapper func(T) U) Stream[U] {
	return Stream[U]{
		Iterate: func() Iterator[U] {
			next := s.Iterate()

			return func() (item U, ok bool) {
				var it T
				it, ok = next()
				if ok {
					item = mapper(it)
				}
				return
			}
		},
	}
}

// FlatMap projects each element of a stream to a Stream, iterates and
// flattens the resulting streams into one stream.
func (s Stream[T]) FlatMap[U any](selector func(T) Stream[U]) Stream[U] {
	return Stream[U]{
		Iterate: func() Iterator[U] {
			outernext := s.Iterate()
			var hasInner bool
			var innernext Iterator[U]

			return func() (item U, ok bool) {
				for !ok {
					if !hasInner {
						var outer T
						outer, ok = outernext()
						if !ok {
							return
						}

						innernext = selector(outer).Iterate()
						hasInner = true
						ok = false
					}

					item, ok = innernext()
					if !ok {
						hasInner = false
					}
				}

				return
			}
		},
	}
}

// Filter returns a stream containing only elements matching the predicate.
func (s Stream[T]) Filter(predicate func(T) bool) Stream[T] {
	return Stream[T]{
		Iterate: func() Iterator[T] {
			next := s.Iterate()

			return func() (item T, ok bool) {
				for item, ok = next(); ok; item, ok = next() {
					if predicate(item) {
						return
					}
				}
				return
			}
		},
	}
}

// Peek applies a function to each element without modifying the stream.
func (s Stream[T]) Peek(fn func(T)) Stream[T] {
	return Stream[T]{
		Iterate: func() Iterator[T] {
			next := s.Iterate()

			return func() (item T, ok bool) {
				item, ok = next()
				if ok {
					fn(item)
				}
				return
			}
		},
	}
}

// Distinct removes duplicate elements from the stream.
// T must be comparable.
func (s Stream[T]) Distinct() Stream[T] {
	return Stream[T]{
		Iterate: func() Iterator[T] {
			next := s.Iterate()
			set := make(map[any]bool)

			return func() (item T, ok bool) {
				for item, ok = next(); ok; item, ok = next() {
					if _, has := set[item]; !has {
						set[item] = true
						return
					}
				}
				return
			}
		},
	}
}

// DistinctBy removes duplicates based on a key selector.
func (s Stream[T]) DistinctBy[K comparable](selector func(T) K) Stream[T] {
	return Stream[T]{
		Iterate: func() Iterator[T] {
			next := s.Iterate()
			set := make(map[K]bool)

			return func() (item T, ok bool) {
				for item, ok = next(); ok; item, ok = next() {
					key := selector(item)
					if _, has := set[key]; !has {
						set[key] = true
						return
					}
				}
				return
			}
		},
	}
}

// Drop skips the first n items in the stream.
func (s Stream[T]) Drop(n int) Stream[T] {
	return Stream[T]{
		Iterate: func() Iterator[T] {
			next := s.Iterate()
			c := n

			return func() (item T, ok bool) {
				for ; c > 0; c-- {
					item, ok = next()
					if !ok {
						return
					}
				}
				return next()
			}
		},
	}
}

// Limit truncates the stream to at most n items.
func (s Stream[T]) Limit(n int) Stream[T] {
	return Stream[T]{
		Iterate: func() Iterator[T] {
			next := s.Iterate()
			c := n

			return func() (item T, ok bool) {
				if c <= 0 {
					return
				}
				c--
				return next()
			}
		},
	}
}
