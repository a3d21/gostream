package gostream

// Group is a generic type that stores the result of GroupBy method.
type Group[K comparable, V any] struct {
	Key   K
	Group []V
}

// GroupBy method groups the elements of a collection according to a specified
// key selector function and projects the elements for each group by using a
// specified function.
func (q Stream[T]) GroupBy[K comparable, V any](keySelector func(T) K,
	elementSelector func(T) V) Stream[Group[K, V]] {
	return Stream[Group[K, V]]{
		func() Iterator[Group[K, V]] {
			next := q.Iterate()
			set := make(map[K][]V)
			// Preserve insertion order of keys.
			var orderedKeys []K

			for item, ok := next(); ok; item, ok = next() {
				key := keySelector(item)
				if _, exists := set[key]; !exists {
					orderedKeys = append(orderedKeys, key)
				}
				set[key] = append(set[key], elementSelector(item))
			}

			length := len(orderedKeys)
			groups := make([]Group[K, V], length)
			for i, k := range orderedKeys {
				groups[i] = Group[K, V]{k, set[k]}
			}

			index := 0

			return func() (item Group[K, V], ok bool) {
				ok = index < length
				if ok {
					item = groups[index]
					index++
				}
				return
			}
		},
	}
}
