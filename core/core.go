package core

import (
	"github.com/a3d21/gostream"
)

// Core Fn for dot-import
var From = gostream.From
var ToSlice = gostream.ToSlice
var ToSliceBy = gostream.ToSliceBy
var ToMap = gostream.ToMap
var ToMapBy = gostream.ToMapBy
var ToSet = gostream.ToSet
var GroupBy = gostream.GroupBy
var Count = gostream.Count
var Collector = gostream.Collector

type KeyValue = gostream.KeyValue

// GTuple 通用Tuple, 可用于多数值比较排序
type GTuple = gostream.GTuple
