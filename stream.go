package gostream

// Range make a int-range stream from `start` to `end`, aka [start, end).
func Range(start, end int) Stream {
	return Stream{
		Iterate: func() Iterator {
			current := start
			return func() (item interface{}, ok bool) {
				if current >= end {
					return nil, false
				}
				item, ok = current, true

				current++
				return
			}
		},
	}
}

// Repeat ...
func Repeat(value interface{}, count int) Stream {
	return Stream{
		Iterate: func() Iterator {
			index := 0

			return func() (item interface{}, ok bool) {
				if index >= count {
					return nil, false
				}

				item, ok = value, true

				index++
				return
			}
		},
	}
}

// Map ...
func (s Stream) Map(mapper normalizedFn) Stream {
	return Stream{
		Iterate: func() Iterator {
			next := s.Iterate()

			return func() (item interface{}, ok bool) {
				var it interface{}
				it, ok = next()
				if ok {
					item = mapper(it)
				}

				return
			}
		},
	}
}

// FlatMap projects each element of a collection to a Query, iterates and
// flattens the resulting collection into one collection.
func (s Stream) FlatMap(selector func(interface{}) Stream) Stream {
	return Stream{
		Iterate: func() Iterator {
			outernext := s.Iterate()
			var inner interface{}
			var innernext Iterator

			return func() (item interface{}, ok bool) {
				for !ok {
					if inner == nil {
						inner, ok = outernext()
						if !ok {
							return
						}

						innernext = selector(inner).Iterate()
					}

					item, ok = innernext()
					if !ok {
						inner = nil
					}
				}

				return
			}
		},
	}
}

// Filter ...
func (s Stream) Filter(predicate func(interface{}) bool) Stream {
	return Stream{
		Iterate: func() Iterator {
			next := s.Iterate()

			return func() (item interface{}, ok bool) {
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

// Peek ...
func (s Stream) Peek(fn func(interface{})) Stream {
	return Stream{
		Iterate: func() Iterator {
			next := s.Iterate()

			return func() (item interface{}, ok bool) {
				item, ok = next()
				if ok {
					fn(item)
				}
				return
			}
		},
	}
}

// Distinct ...
func (s Stream) Distinct() Stream {
	return Stream{
		Iterate: func() Iterator {
			next := s.Iterate()
			set := make(map[interface{}]bool)

			return func() (item interface{}, ok bool) {
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

// DistinctBy 除重
func (s Stream) DistinctBy(selector normalizedFn) Stream {
	return Stream{
		Iterate: func() Iterator {
			next := s.Iterate()
			set := make(map[interface{}]bool)

			return func() (item interface{}, ok bool) {
				for item, ok = next(); ok; item, ok = next() {
					s := selector(item)
					if _, has := set[s]; !has {
						set[s] = true
						return
					}
				}

				return
			}
		},
	}
}

// Drop 丢弃前n项
func (s Stream) Drop(n int) Stream {
	return Stream{
		Iterate: func() Iterator {
			next := s.Iterate()
			c := n

			return func() (item interface{}, ok bool) {
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

// Limit 限制长多n项
func (s Stream) Limit(n int) Stream {
	return Stream{
		Iterate: func() Iterator {
			next := s.Iterate()
			c := n

			return func() (item interface{}, ok bool) {
				if c <= 0 {
					return
				}

				c--
				return next()
			}
		},
	}
}
