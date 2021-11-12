package gopark

import (
	"reflect"
)

var (
	boolTyp = reflect.TypeOf(true)
	boolVal = reflect.ValueOf(true)
)

// Slice2Map 装slice[T]转成成map[T]bool，且值为true
func Slice2Map(vs interface{}) interface{} {
	t := reflect.TypeOf(vs)
	if t.Kind() != reflect.Slice {
		panic("typ should be slice")
	}

	mtyp := reflect.MapOf(t.Elem(), boolTyp)
	mval := reflect.Indirect(reflect.MakeMap(mtyp))

	v := reflect.ValueOf(vs)
	vlen := v.Len()
	for i := 0; i < vlen; i++ {
		mval.SetMapIndex(v.Index(i), boolVal)
	}

	return mval.Interface()
}
