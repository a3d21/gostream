package gostream

import (
	"reflect"

	"github.com/ahmetb/go-linq/v3"
)

type collector struct {
	supplier    func() interface{}
	accumulator func(a interface{}, it interface{}) interface{}
}

var identity = func(it interface{}) interface{} { return it }

// Count 收集器，统计数量
func Count() collector {
	t := reflect.TypeOf(int(0))
	supplier := func() interface{} { return reflect.Indirect(reflect.New(t)) }
	accumulator := func(a interface{}, it interface{}) interface{} {
		return reflect.ValueOf(a.(reflect.Value).Interface().(int) + 1)
	}
	return collector{
		supplier:    supplier,
		accumulator: accumulator,
	}
}

// ToSlice 收集器，将item收集为slice。
// typ为类型参数，允许为nil。 eg: []int{} or []int(nil)
func ToSlice(typ interface{}) collector {
	return ToSliceBy(typ, identity)
}

// ToSliceBy 收集器，将mapper应用于每一个item，再收集结果
func ToSliceBy(typ interface{}, mapper normalizedFn) collector {
	t := reflect.TypeOf(typ)
	if t.Kind() != reflect.Slice {
		panic("typ should be slice")
	}

	supplier := func() interface{} { return reflect.Indirect(reflect.New(t)) }
	accumulator := func(a interface{}, it interface{}) interface{} {
		return reflect.Append(a.(reflect.Value), reflect.ValueOf(mapper(it)))
	}
	return collector{
		supplier:    supplier,
		accumulator: accumulator,
	}
}

// ToMap 收集器
func ToMap(typ interface{}, keyMapper, valueMapper normalizedFn) collector {
	t := reflect.TypeOf(typ)
	if t.Kind() != reflect.Map {
		panic("typ should be map")
	}

	supplier := func() interface{} { return reflect.Indirect(reflect.MakeMap(t)) }
	accumulator := func(a interface{}, it interface{}) interface{} {
		key := keyMapper(it)
		val := valueMapper(it)
		a.(reflect.Value).SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
		return a
	}

	return collector{
		supplier:    supplier,
		accumulator: accumulator,
	}
}

// ToSet 收集器。收集为map[T]bool
func ToSet(typ interface{}) collector {
	t := reflect.TypeOf(typ)
	if t.Kind() != reflect.Map || t.Elem().Kind() != reflect.Bool {
		panic("typ should be map[T]bool")
	}

	trueVal := reflect.ValueOf(true)
	supplier := func() interface{} { return reflect.Indirect(reflect.MakeMap(t)) }
	accumulator := func(a interface{}, it interface{}) interface{} {
		a.(reflect.Value).SetMapIndex(reflect.ValueOf(it), trueVal)
		return a
	}
	return collector{
		supplier:    supplier,
		accumulator: accumulator,
	}
}

// GroupBy 分组收集器，将item分组收集。
// 参数说明：
//   classifier  分组函数
//   downstream  下游收集器
func GroupBy(typ interface{}, classifier normalizedFn, downstream collector) collector {
	t := reflect.TypeOf(typ)
	if t.Kind() != reflect.Map {
		panic("typ should be map")
	}
	supplier := func() interface{} { return reflect.Indirect(reflect.MakeMap(t)) }
	accumulator := func(a interface{}, it interface{}) interface{} {
		key := classifier(it)
		keyV := reflect.ValueOf(key)
		container := a.(reflect.Value).MapIndex(keyV)
		if !container.IsValid() {
			container = downstream.supplier().(reflect.Value)
		}
		a.(reflect.Value).SetMapIndex(keyV, downstream.accumulator(container, it).(reflect.Value))
		return a
	}
	return collector{
		supplier:    supplier,
		accumulator: accumulator,
	}
}

func (s Stream) collect(c collector) reflect.Value {
	return linq.Query(s).AggregateWithSeed(c.supplier(), c.accumulator).(reflect.Value)
}

// Collect ...
func (s Stream) Collect(c collector) interface{} {
	return s.collect(c).Interface()
}
