package core

import (
	"github.com/a3d21/gostream/v2"
)

// Core Fn for dot-import

type KeyValue[K any, V any] = gostream.KeyValue[K, V]
type Stream[T any] = gostream.Stream[T]
type Group[K comparable, V any] = gostream.Group[K, V]
type Collector[T any, R any] = gostream.Collector[T, R]

// GTuple is a generic tuple, useful for multi-value comparison and sorting.
type GTuple = gostream.GTuple

func From[T any](source []T) Stream[T] { return gostream.From(source) }
func FromMap[K comparable, V any](source map[K]V) Stream[KeyValue[K, V]] {
	return gostream.FromMap(source)
}
func FromChannel[T any](source <-chan T) Stream[T] { return gostream.FromChannel(source) }
func Range(start, end int) Stream[int]             { return gostream.Range(start, end) }
func Repeat[T any](value T, count int) Stream[T]   { return gostream.Repeat(value, count) }

// Collectors
func ToSliceCollector[T any]() Collector[T, []T] {
	return gostream.ToSliceCollector[T]()
}

func ToMapCollector[K comparable, V any]() Collector[KeyValue[K, V], map[K]V] {
	return gostream.ToMapCollector[K, V]()
}

func ToMapByCollector[T any, K comparable, V any](keyMapper func(T) K, valueMapper func(T) V) Collector[T, map[K]V] {
	return gostream.ToMapByCollector(keyMapper, valueMapper)
}

func ToSetCollector[T comparable]() Collector[T, map[T]bool] {
	return gostream.ToSetCollector[T]()
}

func CountCollector[T any]() Collector[T, int] {
	return gostream.CountCollector[T]()
}

func CollectBy[T any, R any](supplier func() R, accumulator func(acc R, item T) R) Collector[T, R] {
	return gostream.CollectBy(supplier, accumulator)
}

func GroupByCollector[T any, K comparable, R any](classifier func(T) K, downstream Collector[T, R]) Collector[T, map[K]R] {
	return gostream.GroupByCollector(classifier, downstream)
}

func Zip2By[L any, R any, O any](left Stream[L], right Stream[R], fn func(L, R) O) Stream[O] {
	return gostream.Zip2By(left, right, fn)
}

func Zip3By[A any, B any, C any, O any](first Stream[A], second Stream[B], third Stream[C], fn func(A, B, C) O) Stream[O] {
	return gostream.Zip3By(first, second, third, fn)
}
