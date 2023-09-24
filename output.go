package gostream

import (
	"reflect"
)

// OutSlice 将item输出，out类型为*[]T
func (s Stream) OutSlice(out interface{}) {
	res := reflect.ValueOf(out)
	slice := reflect.Indirect(res)

	cap := slice.Cap()
	res.Elem().Set(slice.Slice(0, cap)) // make len(slice)==cap(slice) from now on

	next := s.Iterate()
	index := 0
	for item, ok := next(); ok; item, ok = next() {
		if index >= cap {
			slice, cap = grow(slice)
		}
		slice.Index(index).Set(reflect.ValueOf(item))
		index++
	}

	// reslice the len(res)==cap(res) actual res size
	res.Elem().Set(slice.Slice(0, index))
}

// OutMap 将item输出，out类型为 *map[K]V
func (s Stream) OutMap(out interface{}) {
	s.OutMapBy(
		out,
		func(i interface{}) interface{} {
			return i.(KeyValue).Key
		},
		func(i interface{}) interface{} {
			return i.(KeyValue).Value
		})
}

// OutMapBy iterates over a collection and populates the result map with
// elements. Functions keySelector and valueSelector are executed for each
// element of the collection to generate key and value for the map. Generated
// key and value types must be assignable to the map's key and value types.
// OutMapBy doesn't empty the result map before populating it.
func (s Stream) OutMapBy(out interface{},
	keySelector func(interface{}) interface{},
	valueSelector func(interface{}) interface{}) {
	res := reflect.ValueOf(out)
	m := reflect.Indirect(res)
	next := s.Iterate()

	for item, ok := next(); ok; item, ok = next() {
		key := reflect.ValueOf(keySelector(item))
		value := reflect.ValueOf(valueSelector(item))

		m.SetMapIndex(key, value)
	}

	res.Elem().Set(m)
}

// OutChan 将item输出chan
func (s Stream) OutChan(ch chan<- interface{}) {
	next := s.Iterate()

	for item, ok := next(); ok; item, ok = next() {
		ch <- item
	}

	close(ch)
}

// OutChanT 将item输出chan
func (s Stream) OutChanT(ch interface{}) {
	r := reflect.ValueOf(ch)
	next := s.Iterate()

	for item, ok := next(); ok; item, ok = next() {
		r.Send(reflect.ValueOf(item))
	}

	r.Close()
}
