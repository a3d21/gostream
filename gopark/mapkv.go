package gopark

import (
	"reflect"
)

// Keys return map's key slice
func Keys(m interface{}) interface{} {
	mval := reflect.ValueOf(m)
	if mval.Kind() != reflect.Map {
		panic("typ should be map")
	}

	styp := reflect.SliceOf(mval.Type().Key())
	slen := mval.Len()
	sval := reflect.MakeSlice(styp, slen, slen)
	for i, k := range mval.MapKeys() {
		sval.Index(i).Set(k)
	}

	return sval.Interface()
}

// Values return map's value slice
func Values(m interface{}) interface{} {
	mval := reflect.ValueOf(m)
	if mval.Kind() != reflect.Map {
		panic("typ should be map")
	}

	styp := reflect.SliceOf(mval.Type().Elem())
	slen := mval.Len()
	sval := reflect.MakeSlice(styp, slen, slen)
	for i, k := range mval.MapKeys() {
		v := mval.MapIndex(k)
		sval.Index(i).Set(v)
	}
	return sval.Interface()
}
