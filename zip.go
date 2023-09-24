package gostream

// Zip2By ...
func Zip2By(left, right Stream, fn zip2Fn) Stream {
	return Stream{
		Iterate: func() Iterator {
			leftNext := left.Iterate()
			rightNext := right.Iterate()

			return func() (interface{}, bool) {
				if l, ok1 := leftNext(); ok1 {
					if r, ok2 := rightNext(); ok2 {
						return fn(l, r), true
					}
				}

				return nil, false
			}
		},
	}
}

// Zip3By ...
func Zip3By(first, second, third Stream, fn zip3Fn) Stream {
	return Stream{
		Iterate: func() Iterator {
			firstNext := first.Iterate()
			secondNext := second.Iterate()
			thirdNext := third.Iterate()

			return func() (interface{}, bool) {
				if it1, ok1 := firstNext(); ok1 {
					if it2, ok2 := secondNext(); ok2 {
						if it3, ok3 := thirdNext(); ok3 {
							return fn(it1, it2, it3), true
						}
					}
				}
				return nil, false
			}
		},
	}
}
