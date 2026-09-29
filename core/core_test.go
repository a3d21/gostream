package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollectToSlice(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	got := From(input).ToSlice()
	assert.Equal(t, input, got)
}

func TestCollectToSliceWithMap(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	got := From(input).Map(func(it int) int {
		return it * 2
	}).ToSlice()
	want := []int{2, 4, 6, 8, 10}
	assert.Equal(t, want, got)
}

func TestCollectToMap(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	got := From(input).Filter(func(it int) bool {
		return it < 4
	}).Map(func(it int) KeyValue[int, int] {
		return KeyValue[int, int]{it, it}
	}).Collect(ToMapCollector[int, int]())
	want := map[int]int{1: 1, 2: 2, 3: 3}
	assert.Equal(t, want, got)
}

func TestCollectToMapBy(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	got := From(input).Filter(func(it int) bool {
		return it < 4
	}).Collect(ToMapByCollector(
		func(it int) int { return it + 1 },
		func(it int) int { return it * 2 },
	))
	want := map[int]int{2: 2, 3: 4, 4: 6}
	assert.Equal(t, want, got)
}

func TestCollectToSet(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	got := From(input).Collect(ToSetCollector[int]())
	want := map[int]bool{1: true, 2: true, 3: true, 4: true, 5: true}
	assert.Equal(t, want, got)
}

func TestCollectGroupBy(t *testing.T) {
	input := []int{11, 21, 31, 41, 12, 22, 32, 42, 13, 23, 33, 43, 14, 24, 34, 44}
	got := From(input).Collect(GroupByCollector(
		func(it int) int { return it / 10 },
		GroupByCollector(
			func(it int) int { return it % 10 },
			ToSliceCollector[int](),
		),
	))

	want := map[int]map[int][]int{
		1: {
			1: {11},
			2: {12},
			3: {13},
			4: {14},
		},
		2: {
			1: {21},
			2: {22},
			3: {23},
			4: {24},
		},
		3: {
			1: {31},
			2: {32},
			3: {33},
			4: {34},
		},
		4: {
			1: {41},
			2: {42},
			3: {43},
			4: {44},
		},
	}
	assert.Equal(t, want, got)
}
