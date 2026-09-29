package gopark

// Slice2Map converts []T to map[T]bool with all values set to true.
func Slice2Map[T comparable](vs []T) map[T]bool {
	m := map[T]bool{}
	for _, v := range vs {
		m[v] = true
	}
	return m
}
