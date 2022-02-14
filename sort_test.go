package gostream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortedByLess(t *testing.T) {
	want := []int{1, 2, 3, 4, 5}
	input := []int{5, 3, 1, 2, 4}
	got := From(input).Sorted(func(a, b interface{}) bool {
		return a.(int) < b.(int)
	}).Collect(ToSlice([]int{})).([]int)
	assert.Equal(t, want, got)
}

func TestSortedByLessOnStruct(t *testing.T) {
	type AStruct struct {
		Name string
		Age  int
	}

	want := []AStruct{{"aaa", 16}, {"aaa", 17}, {"bbb", 14}, {"bbb", 21}}
	input := []AStruct{{"bbb", 21}, {"bbb", 14}, {"aaa", 16}, {"aaa", 17}}
	got := From(input).Sorted(func(a, b interface{}) bool {
		aItem := a.(AStruct)
		bItem := b.(AStruct)
		return aItem.Name < bItem.Name || aItem.Age < bItem.Age
	}).Collect(ToSlice([]AStruct{})).([]AStruct)
	assert.Equal(t, want, got)
}
