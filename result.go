package gostream

import (
	"reflect"
)

// All 判断所有item满足条件，当stream为nil时默认返回true
func (s Stream) All(predicate func(interface{}) bool) bool {
	next := s.Iterate()

	for item, ok := next(); ok; item, ok = next() {
		if !predicate(item) {
			return false
		}
	}

	return true
}

// Any 判断stream是否非空
func (s Stream) Any() bool {
	_, ok := s.Iterate()()
	return ok
}

// AnyWith 判断是否有item满足条件
func (s Stream) AnyWith(predicate func(interface{}) bool) bool {
	next := s.Iterate()

	for item, ok := next(); ok; item, ok = next() {
		if predicate(item) {
			return true
		}
	}

	return false
}

// Contains 判断是否包含value。注意value必须为*值类型*
func (s Stream) Contains(value interface{}) bool {
	next := s.Iterate()

	for item, ok := next(); ok; item, ok = next() {
		if item == value {
			return true
		}
	}

	return false
}

// Count 统计item数量
func (s Stream) Count() int {
	var r int
	next := s.Iterate()

	for _, ok := next(); ok; _, ok = next() {
		r++
	}

	return r
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

// ForEach 遍历item
func (s Stream) ForEach(action func(interface{})) {
	next := s.Iterate()

	for item, ok := next(); ok; item, ok = next() {
		action(item)
	}
}

// Reduce ...
func (s Stream) Reduce(accumulator accumulatorFn) interface{} {
	next := s.Iterate()

	result, any := next()
	if !any {
		return nil
	}

	for current, ok := next(); ok; current, ok = next() {
		result = accumulator(result, current)
	}

	return result
}

// ReduceWith ...
func (s Stream) ReduceWith(seed interface{}, accumulator accumulatorFn) interface{} {
	next := s.Iterate()
	result := seed

	for current, ok := next(); ok; current, ok = next() {
		result = accumulator(result, current)
	}

	return result
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

// grow grows the slice s by doubling its capacity, then it returns the new
// slice (resliced to its full capacity) and the new capacity.
func grow(s reflect.Value) (v reflect.Value, newCap int) {
	cap := s.Cap()
	if cap == 0 {
		cap = 1
	} else {
		cap *= 2
	}
	newSlice := reflect.MakeSlice(s.Type(), cap, cap)
	reflect.Copy(newSlice, s)
	return newSlice, cap
}
