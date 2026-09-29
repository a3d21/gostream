package gostream

// Iterator is a generic alias for function to iterate over data.
type Iterator[T any] func() (item T, ok bool)

// Stream is the generic stream type. It can be iterated manually
// as shown in the example.
type Stream[T any] struct {
	Iterate func() Iterator[T]
}

// KeyValue is a generic key-value pair.
type KeyValue[K any, V any] struct {
	Key   K
	Value V
}

// Iterable is a generic interface for types that can produce an Iterator.
type Iterable[T any] interface {
	Iterate() Iterator[T]
}

// From initializes a Stream from a slice.
func From[T any](source []T) Stream[T] {
	length := len(source)
	return Stream[T]{
		Iterate: func() Iterator[T] {
			index := 0
			return func() (item T, ok bool) {
				ok = index < length
				if ok {
					item = source[index]
					index++
				}
				return
			}
		},
	}
}

// FromMap initializes a Stream from a map, producing KeyValue elements.
func FromMap[K comparable, V any](source map[K]V) Stream[KeyValue[K, V]] {
	return Stream[KeyValue[K, V]]{
		Iterate: func() Iterator[KeyValue[K, V]] {
			// Collect keys upfront for deterministic iteration.
			keys := make([]K, 0, len(source))
			for k := range source {
				keys = append(keys, k)
			}
			length := len(keys)
			index := 0

			return func() (item KeyValue[K, V], ok bool) {
				ok = index < length
				if ok {
					key := keys[index]
					item = KeyValue[K, V]{
						Key:   key,
						Value: source[key],
					}
					index++
				}
				return
			}
		},
	}
}

// FromChannel initializes a Stream from a channel, iterating until the channel is closed.
func FromChannel[T any](source <-chan T) Stream[T] {
	return Stream[T]{
		Iterate: func() Iterator[T] {
			return func() (item T, ok bool) {
				item, ok = <-source
				return
			}
		},
	}
}

// FromString initializes a Stream from a string, iterating over its runes.
func FromString(source string) Stream[rune] {
	runes := []rune(source)
	length := len(runes)

	return Stream[rune]{
		Iterate: func() Iterator[rune] {
			index := 0

			return func() (item rune, ok bool) {
				ok = index < length
				if ok {
					item = runes[index]
					index++
				}
				return
			}
		},
	}
}

// FromIterable initializes a Stream from a type implementing the Iterable interface.
func FromIterable[T any](source Iterable[T]) Stream[T] {
	return Stream[T]{
		Iterate: source.Iterate,
	}
}

// Range make a int-range stream from `start` to `end`, aka [start, end).
func Range(start, end int) Stream[int] {
	return Stream[int]{
		Iterate: func() Iterator[int] {
			current := start
			return func() (item int, ok bool) {
				if current >= end {
					return
				}
				item, ok = current, true
				current++
				return
			}
		},
	}
}

// Repeat creates a stream that repeats `value` for `count` times.
func Repeat[T any](value T, count int) Stream[T] {
	return Stream[T]{
		Iterate: func() Iterator[T] {
			index := 0
			return func() (item T, ok bool) {
				if index >= count {
					return
				}
				item, ok = value, true
				index++
				return
			}
		},
	}
}
