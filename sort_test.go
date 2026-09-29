package gostream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortedByLess(t *testing.T) {
	want := []int{1, 2, 3, 4, 5}
	input := []int{5, 3, 1, 2, 4}
	got := From(input).Sorted(func(a, b int) bool {
		return a < b
	}).ToSlice()
	assert.Equal(t, want, got)
}

func TestSortedByLessOnStruct(t *testing.T) {
	type AStruct struct {
		Name string
		Age  int
	}

	want := []AStruct{{"aaa", 16}, {"aaa", 17}, {"bbb", 14}, {"bbb", 21}}
	input := []AStruct{{"bbb", 21}, {"bbb", 14}, {"aaa", 16}, {"aaa", 17}}
	got := From(input).Sorted(func(a, b AStruct) bool {
		return a.Name < b.Name || a.Age < b.Age
	}).ToSlice()
	assert.Equal(t, want, got)
}

func TestShuffle(t *testing.T) {
	input := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	got := From(input).Shuffle().ToSlice()
	sorted := From(got).SortedBy(func(it int) int { return it }).ToSlice()
	assert.Equal(t, len(input), len(got))
	assert.NotEqual(t, input, got)
	assert.Equal(t, input, sorted)
}
