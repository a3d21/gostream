package gostream

import (
	"reflect"

	"github.com/ahmetb/go-linq/v3"
)

var identity = func(it interface{}) interface{} { return it }

type collector func(stream Stream) interface{}

// Collector custom collector
func Collector(supplier func() interface{}, accumulator accumulatorFn) collector {
	return func(s Stream) interface{} {
		return s.ReduceWith(supplier(), accumulator)
	}
}

// Count ...
func Count() collector {
	return func(s Stream) interface{} {
		return s.Count()
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

	return func(stream Stream) interface{} {
		v := reflect.New(t)
		container := v.Interface()
		stream.Map(mapper).Linq().ToSlice(container)
		return v.Elem().Interface()
	}
}

// ToMap ...
func ToMap(typ interface{}, keyMapper, valueMapper normalizedFn) collector {
	t := reflect.TypeOf(typ)
	if t.Kind() != reflect.Map {
		panic("typ should be map")
	}

	return func(stream Stream) interface{} {
		v := reflect.New(reflect.MapOf(t.Key(), t.Elem()))
		v.Elem().Set(reflect.MakeMap(t))
		container := v.Interface()
		stream.Linq().ToMapBy(container, keyMapper, valueMapper)
		return v.Elem().Interface()
	}
}

// ToSet 收集器。收集为map[T]bool
func ToSet(typ interface{}) collector {
	t := reflect.TypeOf(typ)
	if t.Kind() != reflect.Map || t.Elem().Kind() != reflect.Bool {
		panic("typ should be map[T]bool")
	}

	return func(stream Stream) interface{} {
		v := reflect.New(reflect.MapOf(t.Key(), t.Elem()))
		v.Elem().Set(reflect.MakeMap(t))
		container := v.Interface()
		truly := func(_ interface{}) interface{} { return true }
		stream.Linq().ToMapBy(container, identity, truly)
		return v.Elem().Interface()
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

	return func(stream Stream) interface{} {
		v := reflect.New(reflect.MapOf(t.Key(), t.Elem()))
		v.Elem().Set(reflect.MakeMap(t))
		container := v.Interface()
		stream.Linq().GroupBy(classifier, identity).Select(func(g interface{}) interface{} {
			return KeyValue{
				Key:   g.(linq.Group).Key,
				Value: From(g.(linq.Group).Group).Collect(downstream),
			}
		}).ToMap(container)
		return v.Elem().Interface()
	}
}

// Collect ...
func (s Stream) Collect(c collector) interface{} {
	return c(s)
}
