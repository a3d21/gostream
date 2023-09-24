package ggopark

// PartitionBy 对slice按size分区
func PartitionBy[T any](vs []T, size int) [][]T {
	if size < 1 {
		panic("illegal size")
	}

	vlen := len(vs)
	resultLen := (vlen + size - 1) / size

	var result [][]T
	for i := 0; i < resultLen; i++ {
		begin := i * size
		end := begin + size
		if end > vlen {
			end = vlen
		}
		result = append(result, vs[begin:end])
	}

	return result
}
