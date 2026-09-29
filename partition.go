package gostream

// Partition divides the stream into chunks of the given size.
func Partition[T any](s Stream[T], size int) Stream[[]T] {
	if size < 1 {
		panic("invalid partition size")
	}

	return Stream[[]T]{
		Iterate: func() Iterator[[]T] {
			next := s.Iterate()

			return func() ([]T, bool) {
				batch := make([]T, 0, size)

				for len(batch) < size {
					if it, ok := next(); ok {
						batch = append(batch, it)
					} else {
						break
					}
				}

				if len(batch) > 0 {
					return batch, true
				}

				return nil, false
			}
		},
	}
}
