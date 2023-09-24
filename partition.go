package gostream

import (
	"reflect"
)

// Partition 将Stream按size分区
func (s Stream) Partition(typ interface{}, size int) Stream {
	if size < 1 {
		panic("invalid partition size")
	}
	t := reflect.TypeOf(typ)
	if t.Kind() != reflect.Slice {
		panic("typ should be slice")
	}

	return Stream{
		Iterate: func() Iterator {
			next := s.Iterate()

			return func() (interface{}, bool) {
				sv := reflect.MakeSlice(t, size, size)

				idx := 0
				for idx < size {
					if it, ok := next(); ok {
						sv.Index(idx).Set(reflect.ValueOf(it))
						idx++
					} else {
						break
					}
				}

				if idx > 0 {
					return sv.Slice(0, idx).Interface(), true
				}

				return nil, false
			}
		},
	}
}
