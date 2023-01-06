package gg

// Slice2Map 将[]T转成map[T]bool，且值为true
func Slice2Map[T comparable](vs []T) map[T]bool {
	m := map[T]bool{}
	for _, v := range vs {
		m[v] = true
	}
	return m
}
