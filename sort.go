package gostream

import (
	"math/rand"
	"sort"
)

// SortedBy sorts the stream by comparing values selected by `selector`.
func (s Stream[T]) SortedBy[U comparable](selector func(T) U) Stream[T] {
	less := func(a, b T) bool {
		x, y := selector(a), selector(b)
		c := getComparer(x)
		res := c(x, y)
		return res < 0
	}
	return s.Sorted(less)
}

// SortedDescBy sorts the stream in descending order by comparing values selected by `selector`.
func (s Stream[T]) SortedDescBy[U comparable](selector func(T) U) Stream[T] {
	less := func(a, b T) bool {
		x, y := selector(a), selector(b)
		c := getComparer(x)
		res := c(x, y)
		return res > 0
	}
	return s.Sorted(less)
}

// Sorted sorts elements using the provided less function.
// Parameters:
//
//	less  comparison function. Returns true if a should be placed before b.
func (s Stream[T]) Sorted(less func(a, b T) bool) Stream[T] {
	return Stream[T]{
		Iterate: func() Iterator[T] {
			var items []T
			next := s.Iterate()
			for item, ok := next(); ok; item, ok = next() {
				items = append(items, item)
			}

			itemLen := len(items)
			index := 0

			if itemLen > 0 {
				sort.Slice(items, func(i, j int) bool {
					return less(items[i], items[j])
				})
			}

			return func() (item T, ok bool) {
				ok = index < itemLen
				if ok {
					item = items[index]
					index++
				}
				return
			}
		},
	}
}

// Shuffle randomly shuffles the elements in the stream.
func (s Stream[T]) Shuffle() Stream[T] {
	return Stream[T]{
		Iterate: func() Iterator[T] {
			var items []T
			next := s.Iterate()
			for item, ok := next(); ok; item, ok = next() {
				items = append(items, item)
			}

			itemLen := len(items)
			index := 0

			if itemLen > 0 {
				rand.Shuffle(itemLen, func(i, j int) {
					items[i], items[j] = items[j], items[i]
				})
			}

			return func() (item T, ok bool) {
				ok = index < itemLen
				if ok {
					item = items[index]
					index++
				}
				return
			}
		},
	}
}
