package gostream

import (
	"testing"
)

const (
	size   = 100000
	groups = 100
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
		Range(0, size).ForEach(func(it int) {
			c = append(c, it)
		})
	}
}

func BenchmarkCollectToSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Range(0, size).ToSlice()
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
	for i := 0; i < b.N; i++ {
		Range(0, size).Collect(ToMapByCollector(
			func(it int) int { return it % groups },
			func(it int) int { return it },
		))
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
		Range(0, size).Collect(ToSetCollector[int]())
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
		Range(0, size).Collect(GroupByCollector(
			func(it int) int { return it % groups },
			ToSliceCollector[int](),
		))
	}
}

func BenchmarkPartition(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Partition(Range(0, size), 3).Last()
	}
}

func BenchmarkCountRaw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Range(0, size).Count()
	}
}

func BenchmarkCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Range(0, size).Collect(CountCollector[int]())
	}
}

func BenchmarkGroupCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Range(0, size).Collect(GroupByCollector(
			func(it int) int { return it % groups },
			CountCollector[int](),
		))
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
		Range(0, size).Collect(CollectBy(
			func() int { return 0 },
			func(acc int, item int) int { return acc + item },
		))
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
		Range(0, size).Collect(GroupByCollector(
			func(it int) int { return it % groups },
			CollectBy(
				func() int { return 0 },
				func(acc int, item int) int { return acc + item },
			),
		))
	}
}
