package core

import (
	"github.com/a3d21/gostream"
)

// Core Fn for dot-import

type KeyValue = gostream.KeyValue

// GTuple 通用Tuple, 可用于多数值比较排序
type GTuple = gostream.GTuple

func From(source interface{}) gostream.Stream { return gostream.From(source) }

// generic collector
func ToSlice[T any](typ []T) gostream.C_ { return gostream.ToSlice(typ) }
func ToSliceBy[T, I any](typ []T, mapper func(I) T) gostream.C_ {
	return gostream.ToSliceBy(typ, func(it interface{}) interface{} {
		return mapper(it.(I))
	})
}

func ToMap[K comparable, V any](typ map[K]V) gostream.C_ {
	return gostream.ToMap(typ)
}

func ToMapBy[K comparable, V, I any](typ map[K]V, keyMapper func(I) K, valueMapper func(I) V) gostream.C_ {
	return gostream.ToMapBy(typ, func(it interface{}) interface{} {
		return keyMapper(it.(I))
	}, func(it interface{}) interface{} {
		return valueMapper(it.(I))
	})
}

func ToSet[K comparable](typ map[K]bool) gostream.C_ {
	return gostream.ToSet(typ)
}

func GroupBy[K comparable, S any, I any, O comparable](typ map[K]S, classifier func(I) O, downstream gostream.C_) gostream.C_ {
	return gostream.GroupBy(typ, func(it interface{}) interface{} {
		return classifier(it.(I))
	}, downstream)
}
