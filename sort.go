package gostream

import (
	"math/rand"
	"sort"

	"github.com/ahmetb/go-linq/v3"
)

// SortedBy ...
func (s Stream) SortedBy(fn normalizedFn) Stream {
	return Stream(linq.Query(s).OrderBy(fn).Query)
}

// SortedDescBy ...
func (s Stream) SortedDescBy(fn normalizedFn) Stream {
	return Stream(linq.Query(s).OrderByDescending(fn).Query)
}

// Sorted 按Less函数排序
// 参数说明
//    less函数。 a, b 为item，若a小于b(a排b前面)返回true
func (s Stream) Sorted(less lessFn) Stream {
	return Stream{
		Iterate: func() linq.Iterator {
			var items []interface{}
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

			return func() (item interface{}, ok bool) {
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

// Shuffle ...
func (s Stream) Shuffle() Stream {
	return Stream{
		Iterate: func() linq.Iterator {
			var items []interface{}
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

			return func() (item interface{}, ok bool) {
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
