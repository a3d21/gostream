package gostream

// fork from go-linq

// Comparable is an interface that has to be implemented by a custom collection
// elements in order to work with gostream.
//
// Example:
//
//	func (f foo) CompareTo(c Comparable) int {
//		a, b := f.f1, c.(foo).f1
//
//		if a < b {
//			return -1
//		} else if a > b {
//			return 1
//		}
//
//		return 0
//	}
type Comparable interface {
	CompareTo(Comparable) int
}

func getComparer[T comparable](data T) func(T, T) int {
	switch any(data).(type) {
	case int:
		return func(x, y T) int {
			a, b := any(x).(int), any(y).(int)
			switch {
			case a > b:
				return 1
			case b > a:
				return -1
			default:
				return 0
			}
		}
	case int8:
		return func(x, y T) int {
			a, b := any(x).(int8), any(y).(int8)
			switch {
			case a > b:
				return 1
			case b > a:
				return -1
			default:
				return 0
			}
		}
	case int16:
		return func(x, y T) int {
			a, b := any(x).(int16), any(y).(int16)
			switch {
			case a > b:
				return 1
			case b > a:
				return -1
			default:
				return 0
			}
		}
	case int32:
		return func(x, y T) int {
			a, b := any(x).(int32), any(y).(int32)
			switch {
			case a > b:
				return 1
			case b > a:
				return -1
			default:
				return 0
			}
		}
	case int64:
		return func(x, y T) int {
			a, b := any(x).(int64), any(y).(int64)
			switch {
			case a > b:
				return 1
			case b > a:
				return -1
			default:
				return 0
			}
		}
	case uint:
		return func(x, y T) int {
			a, b := any(x).(uint), any(y).(uint)
			switch {
			case a > b:
				return 1
			case b > a:
				return -1
			default:
				return 0
			}
		}
	case uint8:
		return func(x, y T) int {
			a, b := any(x).(uint8), any(y).(uint8)
			switch {
			case a > b:
				return 1
			case b > a:
				return -1
			default:
				return 0
			}
		}
	case uint16:
		return func(x, y T) int {
			a, b := any(x).(uint16), any(y).(uint16)
			switch {
			case a > b:
				return 1
			case b > a:
				return -1
			default:
				return 0
			}
		}
	case uint32:
		return func(x, y T) int {
			a, b := any(x).(uint32), any(y).(uint32)
			switch {
			case a > b:
				return 1
			case b > a:
				return -1
			default:
				return 0
			}
		}
	case uint64:
		return func(x, y T) int {
			a, b := any(x).(uint64), any(y).(uint64)
			switch {
			case a > b:
				return 1
			case b > a:
				return -1
			default:
				return 0
			}
		}
	case float32:
		return func(x, y T) int {
			a, b := any(x).(float32), any(y).(float32)
			switch {
			case a > b:
				return 1
			case b > a:
				return -1
			default:
				return 0
			}
		}
	case float64:
		return func(x, y T) int {
			a, b := any(x).(float64), any(y).(float64)
			switch {
			case a > b:
				return 1
			case b > a:
				return -1
			default:
				return 0
			}
		}
	case string:
		return func(x, y T) int {
			a, b := any(x).(string), any(y).(string)
			switch {
			case a > b:
				return 1
			case b > a:
				return -1
			default:
				return 0
			}
		}
	case bool:
		return func(x, y T) int {
			a, b := any(x).(bool), any(y).(bool)
			switch {
			case a == b:
				return 0
			case a:
				return 1
			default:
				return -1
			}
		}
	default:
		return func(x, y T) int {
			a, b := any(x).(Comparable), any(y).(Comparable)
			return a.CompareTo(b)
		}
	}
}

// GTuple is a generic comparable tuple.
// It compares elements in sequence. An empty GTuple is considered the smallest.
type GTuple []any

func (l1 GTuple) CompareTo(l2 Comparable) int {
	a, b := l1, l2.(GTuple)
	alen, blen := len(a), len(b)
	if alen == 0 && blen == 0 {
		return 0
	}

	for i := 0; i < alen && i < blen; i++ {
		x, y := a[i], b[i]
		c := getComparer(x)
		res := c(x, y)
		if res < 0 {
			return -1
		} else if res > 0 {
			return 1
		}
	}

	if alen < blen {
		return -1
	} else if alen == blen {
		return 0
	} else {
		return 1
	}
}
