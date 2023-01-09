package ggopark

// Keys return map's key slice
func Keys[K comparable, V any](m map[K]V) []K {
	var ks []K
	for k, _ := range m {
		ks = append(ks, k)
	}
	return ks
}

// Values return map's value slice
func Values[K comparable, V any](m map[K]V) []V {
	var vs []V
	for _, v := range m {
		vs = append(vs, v)
	}
	return vs
}
