package gostream

import (
	"reflect"

	"github.com/ahmetb/go-linq/v3"
)

type collectorV2 func(stream Stream) interface{}

// CountV2 ...
func CountV2() collectorV2 {
	return func(s Stream) interface{} {
		return s.Count()
	}
}

// ToSliceV2 ...
func ToSliceV2(typ interface{}) collectorV2 {
	return ToSliceByV2(typ, identity)
}

// ToSliceByV2 ...
func ToSliceByV2(typ interface{}, mapper normalizedFn) collectorV2 {
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

// ToMapV2 ...
func ToMapV2(typ interface{}, keyMapper, valueMapper normalizedFn) collectorV2 {
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

// ToSetV2 ...
func ToSetV2(typ interface{}) collectorV2 {
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

// GroupByV2 ...
func GroupByV2(typ interface{}, classifier normalizedFn, downstream collectorV2) collectorV2 {
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
				Value: From(g.(linq.Group).Group).CollectV2(downstream),
			}
		}).ToMap(container)
		return v.Elem().Interface()
	}
}

// CollectV2 Collect的v2版本，基于go-linq实现
func (s Stream) CollectV2(c collectorV2) interface{} {
	return c(s)
}
