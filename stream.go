package gostream

import (
	"github.com/ahmetb/go-linq/v3"
)

// Stream 流
type Stream linq.Query

// Linq 转换成go-linq Query
func (s Stream) Linq() linq.Query {
	return linq.Query(s)
}

// From ...
func From(source interface{}) Stream {
	return Stream(linq.From(source))
}

// Range make a int-range stream from `start` to `end`, aka [start, end).
func Range(start, end int) Stream {
	return Stream{
		Iterate: func() linq.Iterator {
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
	return Stream(linq.Repeat(value, count))
}

// Concat ...
func (s Stream) Concat(s2 Stream) Stream {
	return Stream(s.Linq().Concat(s2.Linq()))
}

// Append ...
func (s Stream) Append(item interface{}) Stream {
	return Stream(s.Linq().Append(item))
}

// Map ...
func (s Stream) Map(mapper normalizedFn) Stream {
	return Stream(linq.Query(s).Select(mapper))
}

// FlatMap ...
func (s Stream) FlatMap(mapper func(interface{}) Stream) Stream {
	selector := func(it interface{}) linq.Query {
		return linq.Query(mapper(it))
	}
	return Stream(linq.Query(s).SelectMany(selector))
}

// Filter ...
func (s Stream) Filter(predicate func(interface{}) bool) Stream {
	return Stream(linq.Query(s).Where(predicate))
}

// SortedBy ...
func (s Stream) SortedBy(fn normalizedFn) Stream {
	return Stream(linq.Query(s).OrderBy(fn).Query)
}

// SortedDescBy ...
func (s Stream) SortedDescBy(fn normalizedFn) Stream {
	return Stream(linq.Query(s).OrderByDescending(fn).Query)
}

// Peek 对经过的每一项item应用fn函数
func (s Stream) Peek(fn func(interface{})) Stream {
	return Stream{
		Iterate: func() linq.Iterator {
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

// Distinct 除重。item必须为值类型
func (s Stream) Distinct() Stream {
	return Stream(s.Linq().Distinct())
}

// Drop 丢弃前n项
func (s Stream) Drop(n int) Stream {
	return Stream(s.Linq().Skip(n))
}

// Limit 限制长多n项
func (s Stream) Limit(n int) Stream {
	return Stream(s.Linq().Take(n))
}
