package gostream

import (
	"testing"

	"github.com/ahmetb/go-linq/v3"
)

const (
	size = 100000
)

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
		Stream(linq.Range(1, size)).ForEach(func(it interface{}) {
			c = append(c, it.(int))
		})
	}
}

func BenchmarkCollectToSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Stream(linq.Range(1, size)).Collect(ToSlice([]int(nil)))
	}
}

func BenchmarkCollectToSliceV2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Stream(linq.Range(1, size)).CollectV2(ToSliceV2([]int{}))
	}
}

func BenchmarkLinqToSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := make([]int, 0)
		linq.Range(1, size).ToSlice(&c)
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
		Stream(linq.Range(1, size)).Collect(ToMap(map[int]int(nil), identity, identity))
	}
}

func BenchmarkCollectToMapV2(b *testing.B) {
	identity := func(it interface{}) interface{} { return it }
	for i := 0; i < b.N; i++ {
		Stream(linq.Range(1, size)).CollectV2(ToMapV2(map[int]int(nil), identity, identity))
	}
}

func BenchmarkLinqToMap(b *testing.B) {
	identity := func(it interface{}) interface{} { return it }
	for i := 0; i < b.N; i++ {
		c := make(map[int]int)
		linq.Range(1, size).ToMapBy(&c, identity, identity)
	}
}

////// GroupBy

func BenchmarkGroupByRaw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := make(map[int][]int)
		for j := 0; j < size; j++ {
			down, ok := c[j]
			if !ok {
				down = make([]int, 0)
			}
			down = append(down, j)
			c[j] = down
		}
	}
}

func BenchmarkGroupBy(b *testing.B) {
	identity := func(it interface{}) interface{} { return it }
	for i := 0; i < b.N; i++ {
		Stream(linq.Range(1, size)).Collect(GroupBy(map[int][]int(nil), identity, ToSlice([]int(nil))))
	}
}

func BenchmarkGroupByV2(b *testing.B) {
	identity := func(it interface{}) interface{} { return it }
	for i := 0; i < b.N; i++ {
		Stream(linq.Range(1, size)).CollectV2(GroupByV2(map[int][]int(nil), identity, ToSliceV2([]int(nil))))
	}
}

func BenchmarkLinqGroupBy(b *testing.B) {
	identity := func(it interface{}) interface{} { return it }
	for i := 0; i < b.N; i++ {
		c := make(map[int]interface{})
		linq.Range(1, size).GroupBy(identity, identity).Select(func(it interface{}) interface{} {
			return KeyValue{
				Key:   it.(linq.Group).Key,
				Value: it.(linq.Group).Group,
			}
		}).ToMap(&c)
	}
}

func BenchmarkPartition(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Stream(linq.Range(1, size)).Partition([]int(nil), 3).Last()
	}
}

func BenchmarkCountRaw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Stream(linq.Range(1, size)).Count()
	}
}

func BenchmarkCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Stream(linq.Range(1, size)).Collect(Count())
	}
}

func BenchmarkCountV2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Stream(linq.Range(1, size)).CollectV2(CountV2())
	}
}

func BenchmarkGroupCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Stream(linq.Range(1, size)).Collect(GroupBy(map[int]int(nil), identity, Count()))
	}
}

func BenchmarkGroupCountV2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Stream(linq.Range(1, size)).CollectV2(GroupByV2(map[int]int(nil), identity, CountV2()))
	}
}
