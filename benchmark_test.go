package gostream

import (
	"testing"
)

const (
	size   = 100000
	groups = 100
)

// intGroupBy group by it%count
func intGroupBy(count int) func(it interface{}) interface{} {
	return func(it interface{}) interface{} {
		return it.(int) % count
	}
}

////// ToSlice

func BenchmarkToSliceRaw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := make([]int, 0, size)
		for j := 0; j < size; j++ {
			c = append(c, j)
		}
	}
}

func BenchmarkToSliceStreamForeach(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := make([]int, 0, size)
		Range(0, size).ForEach(func(it interface{}) {
			c = append(c, it.(int))
		})
	}
}

func BenchmarkCollectToSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Range(0, size).Collect(ToSlice([]int(nil)))
	}
}

////// ToMap

func BenchmarkToMapRaw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := make(map[int]int)
		for j := 0; j < size; j++ {
			c[j] = j
		}
	}
}

func BenchmarkCollectToMap(b *testing.B) {
	identity := func(it interface{}) interface{} { return it }
	for i := 0; i < b.N; i++ {
		Range(0, size).Collect(ToMapBy(map[int]int(nil), intGroupBy(groups), identity))
	}
}

////// ToSet

func BenchmarkToSetRaw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := make(map[int]bool)
		for j := 0; j < size; j++ {
			c[j] = true
		}
	}
}

func BenchmarkCollectToSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Range(0, size).Collect(ToSet(map[int]bool(nil)))
	}
}

////// GroupBy

func BenchmarkGroupByRaw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := make(map[int][]int)
		for j := 0; j < size; j++ {
			k := j % groups
			down, ok := c[k]
			if !ok {
				down = make([]int, 0)
			}
			down = append(down, j)
			c[k] = down
		}
	}
}

func BenchmarkGroupBy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Range(0, size).Collect(GroupBy(map[int][]int(nil), intGroupBy(groups), ToSlice([]int(nil))))
	}
}

func BenchmarkPartition(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Range(0, size).Partition([]int(nil), 3).Last()
	}
}

func BenchmarkCountRaw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Range(0, size).Count()
	}
}

func BenchmarkCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Range(0, size).Collect(Count())
	}
}

func BenchmarkGroupCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Range(0, size).Collect(GroupBy(map[int]int(nil), intGroupBy(groups), Count()))
	}
}

func BenchmarkSumRaw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sum := 0
		for j := 0; j < size; j++ {
			sum += j
		}
	}
}

func BenchmarkCustomSumCollector(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Range(0, size).Collect(CollectBy(func() interface{} {
			return 0
		}, func(acc interface{}, item interface{}) interface{} {
			return acc.(int) + item.(int)
		}))
	}
}

func BenchmarkGroupSumRaw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		got := map[int]int{}
		for j := 0; j < size; j++ {
			key := j % groups
			got[key] += j
		}
	}
}

func BenchmarkGroupSum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Range(0, size).Collect(GroupBy(map[int]int{},
			intGroupBy(groups),
			CollectBy(func() interface{} {
				return 0
			}, func(acc interface{}, item interface{}) interface{} {
				return acc.(int) + item.(int)
			})))
	}
}
