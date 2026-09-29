package gostream

// Zip2By combines two streams using the provided function.
func Zip2By[L any, R any, O any](left Stream[L], right Stream[R], fn func(L, R) O) Stream[O] {
	return Stream[O]{
		Iterate: func() Iterator[O] {
			leftNext := left.Iterate()
			rightNext := right.Iterate()

			return func() (O, bool) {
				if l, ok1 := leftNext(); ok1 {
					if r, ok2 := rightNext(); ok2 {
						return fn(l, r), true
					}
				}

				var zero O
				return zero, false
			}
		},
	}
}

// Zip3By combines three streams using the provided function.
func Zip3By[A any, B any, C any, O any](first Stream[A], second Stream[B], third Stream[C], fn func(A, B, C) O) Stream[O] {
	return Stream[O]{
		Iterate: func() Iterator[O] {
			firstNext := first.Iterate()
			secondNext := second.Iterate()
			thirdNext := third.Iterate()

			return func() (O, bool) {
				if it1, ok1 := firstNext(); ok1 {
					if it2, ok2 := secondNext(); ok2 {
						if it3, ok3 := thirdNext(); ok3 {
							return fn(it1, it2, it3), true
						}
					}
				}
				var zero O
				return zero, false
			}
		},
	}
}
