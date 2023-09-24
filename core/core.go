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
func ToSliceBy[T any](typ []T, mapper func(any) any) gostream.C_ {
	return gostream.ToSliceBy(typ, mapper)
}

func ToMap[K comparable, V any](typ map[K]V) gostream.C_ {
	return gostream.ToMap(typ)
}

func ToMapBy[K comparable, V any](typ map[K]V, keyMapper, valueMapper func(any) any) gostream.C_ {
	return gostream.ToMapBy(typ, keyMapper, valueMapper)
}

func ToSet[K comparable](typ map[K]bool) gostream.C_ {
	return gostream.ToSet(typ)
}

func GroupBy[K comparable, S any](typ map[K]S, classifier func(any) any, downstream gostream.C_) gostream.C_ {
	return gostream.GroupBy(typ, classifier, downstream)
}
