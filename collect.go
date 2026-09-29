package gostream

// Collector is a generic collector function.
type Collector[T any, R any] func(stream Stream[T]) R

// Collect applies a collector to the stream.
func (s Stream[T]) Collect[R any](c Collector[T, R]) R {
	return c(s)
}

// CollectBy creates a collector from a supplier and an accumulator.
//
//	supplier: supply the seed
//	accumulator: accumulate items
func CollectBy[T any, R any](supplier func() R, accumulator func(acc R, item T) R) Collector[T, R] {
	return func(s Stream[T]) R {
		result := supplier()
		next := s.Iterate()
		for current, ok := next(); ok; current, ok = next() {
			result = accumulator(result, current)
		}
		return result
	}
}

// CountCollector returns a collector that counts stream elements.
func CountCollector[T any]() Collector[T, int] {
	return func(s Stream[T]) int {
		return s.Count()
	}
}

// ToSliceCollector returns a collector that collects elements into a slice.
func ToSliceCollector[T any]() Collector[T, []T] {
	return func(s Stream[T]) []T {
		return s.ToSlice()
	}
}

// ToMapCollector returns a collector that collects KeyValue elements into a map.
func ToMapCollector[K comparable, V any]() Collector[KeyValue[K, V], map[K]V] {
	return func(s Stream[KeyValue[K, V]]) map[K]V {
		return ToMap(s)
	}
}

// ToMapByCollector returns a collector that applies key and value mappers and collects into a map.
func ToMapByCollector[T any, K comparable, V any](keyMapper func(T) K, valueMapper func(T) V) Collector[T, map[K]V] {
	return func(s Stream[T]) map[K]V {
		result := make(map[K]V)
		next := s.Iterate()
		for item, ok := next(); ok; item, ok = next() {
			result[keyMapper(item)] = valueMapper(item)
		}
		return result
	}
}

// ToSetCollector returns a collector that collects elements into a map[T]bool set.
func ToSetCollector[T comparable]() Collector[T, map[T]bool] {
	return func(s Stream[T]) map[T]bool {
		result := make(map[T]bool)
		next := s.Iterate()
		for item, ok := next(); ok; item, ok = next() {
			result[item] = true
		}
		return result
	}
}

// GroupByCollector creates a grouping collector.
// Parameters:
//
//	classifier  function to classify items into keys
//	downstream  collector for grouping downstream items
func GroupByCollector[T any, K comparable, R any](classifier func(T) K, downstream Collector[T, R]) Collector[T, map[K]R] {
	return func(s Stream[T]) map[K]R {
		// Group elements by key
		groups := make(map[K][]T)
		var orderedKeys []K
		next := s.Iterate()
		for item, ok := next(); ok; item, ok = next() {
			key := classifier(item)
			if _, exists := groups[key]; !exists {
				orderedKeys = append(orderedKeys, key)
			}
			groups[key] = append(groups[key], item)
		}

		// Apply downstream collector to each group
		result := make(map[K]R)
		for _, key := range orderedKeys {
			result[key] = downstream(From(groups[key]))
		}
		return result
	}
}
