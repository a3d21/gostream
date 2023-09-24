package gostream

import (
	"reflect"
)

// Iterator is an alias for function to iterate over data.
type Iterator func() (item interface{}, ok bool)

// Stream is the type returned from query functions. It can be iterated manually
// as shown in the example.
type Stream struct {
	Iterate func() Iterator
}

// KeyValue ...
type KeyValue struct {
	Key   interface{}
	Value interface{}
}

// Iterable ...
type Iterable interface {
	Iterate() Iterator
}

// From initializes a Stream with passed slice, array or map as the source.
// String, channel or struct implementing Iterable interface can be used as an
// input. In this case From delegates it to FromString, FromChannel and
// FromIterable internally.
func From(source interface{}) Stream {
	src := reflect.ValueOf(source)

	switch src.Kind() {
	case reflect.Slice, reflect.Array:
		len := src.Len()

		return Stream{
			Iterate: func() Iterator {
				index := 0

				return func() (item interface{}, ok bool) {
					ok = index < len
					if ok {
						item = src.Index(index).Interface()
						index++
					}

					return
				}
			},
		}
	case reflect.Map:
		len := src.Len()

		return Stream{
			Iterate: func() Iterator {
				index := 0
				keys := src.MapKeys()

				return func() (item interface{}, ok bool) {
					ok = index < len
					if ok {
						key := keys[index]
						item = KeyValue{
							Key:   key.Interface(),
							Value: src.MapIndex(key).Interface(),
						}

						index++
					}

					return
				}
			},
		}
	case reflect.String:
		return FromString(source.(string))
	case reflect.Chan:
		if _, ok := source.(chan interface{}); ok {
			return FromChannel(source.(chan interface{}))
		} else {
			return FromChannelT(source)
		}
	default:
		return FromIterable(source.(Iterable))
	}
}

// FromChannel initializes a Stream with passed channel, gostream iterates over
// channel until it is closed.
func FromChannel(source <-chan interface{}) Stream {
	return Stream{
		Iterate: func() Iterator {
			return func() (item interface{}, ok bool) {
				item, ok = <-source
				return
			}
		},
	}
}

// FromChannelT is the typed version of FromChannel.
//
//   - source is of type "chan TSource"
//
// NOTE: FromChannel has better performance than FromChannelT.
func FromChannelT(source interface{}) Stream {
	src := reflect.ValueOf(source)
	return Stream{
		Iterate: func() Iterator {
			return func() (interface{}, bool) {
				value, ok := src.Recv()
				return value.Interface(), ok
			}
		},
	}
}

// FromString initializes a Stream with passed string, gostream iterates over
// runes of string.
func FromString(source string) Stream {
	runes := []rune(source)
	len := len(runes)

	return Stream{
		Iterate: func() Iterator {
			index := 0

			return func() (item interface{}, ok bool) {
				ok = index < len
				if ok {
					item = runes[index]
					index++
				}

				return
			}
		},
	}
}

// FromIterable initializes a Stream with custom collection passed. This
// collection has to implement Iterable interface, gostream iterates over items,
// that has to implement Comparable interface or be basic types.
func FromIterable(source Iterable) Stream {
	return Stream{
		Iterate: source.Iterate,
	}
}
