package gostream

import (
	"github.com/ahmetb/go-linq/v3"
)

// All 判断所有item满足条件，当stream为nil时默认返回true
func (s Stream) All(predicate func(interface{}) bool) bool {
	return linq.Query(s).All(predicate)
}

// Any 判断stream是否非空
func (s Stream) Any() bool {
	return linq.Query(s).Any()
}

// AnyWith 判断是否有item满足条件
func (s Stream) AnyWith(predicate func(interface{}) bool) bool {
	return linq.Query(s).AnyWith(predicate)
}

// Contains 判断是否包含value。注意value必须为*值类型*
func (s Stream) Contains(value interface{}) bool {
	return linq.Query(s).Contains(value)
}

// Count 统计item数量
func (s Stream) Count() int {
	return linq.Query(s).Count()
}

// First 返回第一个item，若stream为空，返回 (nil, false)
func (s Stream) First() (interface{}, bool) {
	return s.Iterate()()
}

// Last 返回最后一个item，若stream为空，返回（nil, false）
func (s Stream) Last() (interface{}, bool) {
	next := s.Iterate()

	if r, ok := next(); ok {
		for item, ok2 := next(); ok2; item, ok2 = next() {
			r = item
		}
		return r, true
	}
	return nil, false
}

// OutSlice 将item输出，out类型为*[]T
func (s Stream) OutSlice(out interface{}) {
	linq.Query(s).ToSlice(out)
}

// OutMap 将item输出，out类型为 *map[K]V
func (s Stream) OutMap(out interface{}) {
	linq.Query(s).ToMap(out)
}

// ForEach 遍历item
func (s Stream) ForEach(action func(interface{})) {
	linq.Query(s).ForEach(action)
}

// ForEachIndexed 带index遍历item
func (s Stream) ForEachIndexed(action func(int, interface{})) {
	linq.Query(s).ForEachIndexed(action)
}

// Reduce ...
func (s Stream) Reduce(accumulator accumulatorFn) interface{} {
	return linq.Query(s).Aggregate(accumulator)
}

// ReduceWith ...
func (s Stream) ReduceWith(seed interface{}, accumulator accumulatorFn) interface{} {
	return linq.Query(s).AggregateWithSeed(seed, accumulator)
}

// Process 对每个item应用函数，如果其中一项返回err，中断并返回err
func (s Stream) Process(fn func(interface{}) error) error {
	if err, ok := s.Map(func(v interface{}) interface{} {
		return fn(v)
	}).Filter(func(v interface{}) bool {
		return v != nil
	}).First(); ok && err != nil {
		return err.(error)
	}
	return nil
}
