package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCollectToSlice(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	got := From(input).Collect(ToSlice([]int{}))
	assert.Equal(t, input, got)
}

func TestCollectToMap(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	got := From(input).Filter(func(it interface{}) bool {
		return it.(int) < 4
	}).Map(func(it interface{}) interface{} {
		return KeyValue{it, it}
	}).Collect(ToMap(map[int]int{}))
	want := map[int]int{1: 1, 2: 2, 3: 3}
	assert.Equal(t, want, got)
}

func TestCollectToSet(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	got := From(input).Collect(ToSet(map[int]bool{}))
	want := map[int]bool{1: true, 2: true, 3: true, 4: true, 5: true}
	assert.Equal(t, want, got)
}

func TestCollectGroupBy(t *testing.T) {
	input := []int{11, 21, 31, 41, 12, 22, 32, 42, 13, 23, 33, 43, 14, 24, 34, 44}
	got := From(input).Collect(GroupBy(map[int]map[int][]int{}, func(it any) any {
		return it.(int) / 10
	}, GroupBy(map[int][]int{}, func(it any) any {
		return it.(int) % 10
	}, ToSlice([]int{}))))

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
