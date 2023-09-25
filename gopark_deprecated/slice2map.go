package gopark_deprecated

import (
	"reflect"
)

var (
	trueTyp = reflect.TypeOf(true)
	trueVal = reflect.ValueOf(true)
)

// Slice2Map 将[]T转成map[T]bool，且值为true
func Slice2Map(vs interface{}) interface{} {
	t := reflect.TypeOf(vs)
	if t.Kind() != reflect.Slice {
		panic("typ should be slice")
	}

	mtyp := reflect.MapOf(t.Elem(), trueTyp)
	mval := reflect.Indirect(reflect.MakeMap(mtyp))

	v := reflect.ValueOf(vs)
	vlen := v.Len()
	for i := 0; i < vlen; i++ {
		mval.SetMapIndex(v.Index(i), trueVal)
	}

	return mval.Interface()
}
