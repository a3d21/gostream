package gostream

// fork from go-linq

type comparer func(interface{}, interface{}) int

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

func getComparer(data interface{}) comparer {
	switch data.(type) {
	case int:
		return func(x, y interface{}) int {
			a, b := x.(int), y.(int)
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
		return func(x, y interface{}) int {
			a, b := x.(int8), y.(int8)
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
		return func(x, y interface{}) int {
			a, b := x.(int16), y.(int16)
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
		return func(x, y interface{}) int {
			a, b := x.(int32), y.(int32)
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
		return func(x, y interface{}) int {
			a, b := x.(int64), y.(int64)
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
		return func(x, y interface{}) int {
			a, b := x.(uint), y.(uint)
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
		return func(x, y interface{}) int {
			a, b := x.(uint8), y.(uint8)
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
		return func(x, y interface{}) int {
			a, b := x.(uint16), y.(uint16)
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
		return func(x, y interface{}) int {
			a, b := x.(uint32), y.(uint32)
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
		return func(x, y interface{}) int {
			a, b := x.(uint64), y.(uint64)
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
		return func(x, y interface{}) int {
			a, b := x.(float32), y.(float32)
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
		return func(x, y interface{}) int {
			a, b := x.(float64), y.(float64)
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
		return func(x, y interface{}) int {
			a, b := x.(string), y.(string)
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
		return func(x, y interface{}) int {
			a, b := x.(bool), y.(bool)
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
		return func(x, y interface{}) int {
			a, b := x.(Comparable), y.(Comparable)
			return a.CompareTo(b)
		}
	}
}

// GTuple 通用Tuple，可比较
// 依序比较每一个元素。空的GTuple最小
type GTuple []interface{}

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
