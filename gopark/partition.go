package gopark

import (
	"reflect"
)

// PartitionBy 对slice按size分区
func PartitionBy(vs interface{}, size int) interface{} {
	if size < 1 {
		panic("illegal size")
	}
	t := reflect.TypeOf(vs)
	if t.Kind() != reflect.Slice {
		panic("typ should be slice")
	}

	v := reflect.ValueOf(vs)
	vlen := v.Len()
	length := (vlen + size - 1) / size

	resultValue := reflect.MakeSlice(reflect.SliceOf(t), length, length)

	for i := 0; i < length; i++ {
		begin := i * size
		end := begin + size
		if end > vlen {
			end = vlen
		}
		resultValue.Index(i).Set(v.Slice(begin, end))
	}

	return resultValue.Interface()
}
