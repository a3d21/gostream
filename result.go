package gostream

// All checks if all items satisfy the predicate. Returns true by default if the stream is empty.
func (s Stream[T]) All(predicate func(T) bool) bool {
	next := s.Iterate()

	for item, ok := next(); ok; item, ok = next() {
		if !predicate(item) {
			return false
		}
	}

	return true
}

// Any checks if the stream contains any items.
func (s Stream[T]) Any() bool {
	_, ok := s.Iterate()()
	return ok
}

// AnyWith checks if any item satisfies the predicate.
func (s Stream[T]) AnyWith(predicate func(T) bool) bool {
	next := s.Iterate()

	for item, ok := next(); ok; item, ok = next() {
		if predicate(item) {
			return true
		}
	}

	return false
}

// Contains checks if the stream contains the value. Note that value must be comparable.
func (s Stream[T]) Contains(value T) bool {
	next := s.Iterate()

	for item, ok := next(); ok; item, ok = next() {
		if any(item) == any(value) {
			return true
		}
	}

	return false
}

// Count returns the number of items in the stream.
func (s Stream[T]) Count() int {
	var r int
	next := s.Iterate()

	for _, ok := next(); ok; _, ok = next() {
		r++
	}

	return r
}

// First returns the first item. If the stream is empty, returns (zero value, false).
func (s Stream[T]) First() (T, bool) {
	return s.Iterate()()
}

// Last returns the last item. If the stream is empty, returns (zero value, false).
func (s Stream[T]) Last() (T, bool) {
	next := s.Iterate()

	if r, ok := next(); ok {
		for item, ok2 := next(); ok2; item, ok2 = next() {
			r = item
		}
		return r, true
	}
	var zero T
	return zero, false
}

// ForEach applies an action to each item in the stream.
func (s Stream[T]) ForEach(action func(T)) {
	next := s.Iterate()

	for item, ok := next(); ok; item, ok = next() {
		action(item)
	}
}

// Reduce reduces the stream to a single value using the accumulator.
func (s Stream[T]) Reduce(accumulator func(acc T, item T) T) (T, bool) {
	next := s.Iterate()

	result, any := next()
	if !any {
		var zero T
		return zero, false
	}

	for current, ok := next(); ok; current, ok = next() {
		result = accumulator(result, current)
	}

	return result, true
}

// ReduceWith reduces the stream starting with a seed value.
func (s Stream[T]) ReduceWith[R any](seed R, accumulator func(acc R, item T) R) R {
	next := s.Iterate()
	result := seed

	for current, ok := next(); ok; current, ok = next() {
		result = accumulator(result, current)
	}

	return result
}

// Process applies a function to each item. If any call returns an error, processing stops and the error is returned.
func (s Stream[T]) Process(fn func(T) error) error {
	next := s.Iterate()

	for item, ok := next(); ok; item, ok = next() {
		if err := fn(item); err != nil {
			return err
		}
	}
	return nil
}

// ToSlice collects all elements of the stream into a slice.
func (s Stream[T]) ToSlice() []T {
	var result []T
	next := s.Iterate()
	for item, ok := next(); ok; item, ok = next() {
		result = append(result, item)
	}
	return result
}

// ToMap collects all KeyValue elements into a map.
func ToMap[K comparable, V any](s Stream[KeyValue[K, V]]) map[K]V {
	result := make(map[K]V)
	next := s.Iterate()
	for item, ok := next(); ok; item, ok = next() {
		result[item.Key] = item.Value
	}
	return result
}
