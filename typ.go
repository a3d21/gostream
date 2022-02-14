package gostream

import (
	"github.com/ahmetb/go-linq/v3"
)

// KeyValue ...
type KeyValue = linq.KeyValue

type normalizedFn func(interface{}) interface{}
type accumulatorFn func(acc interface{}, item interface{}) interface{}
type zip2Fn func(left, right interface{}) interface{}
type zip3Fn func(first, second, third interface{}) interface{}
type lessFn func(a, b interface{}) bool
