package gostream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStreamDistinct(t *testing.T) {
	input := []int64{1, 2, 2, 3, 3, 3, 4, 4, 5}
	want := []int64{1, 2, 3, 4, 5}

	got := From(input).Distinct().Collect(ToSlice([]int64{}))
	assert.Equal(t, want, got)
}

func TestStreamConcat(t *testing.T) {
	input1 := []int{1, 2, 3, 4}
	input2 := []int{5, 6, 7, 8}
	want := []int{1, 2, 3, 4, 5, 6, 7, 8}

	got := From(input1).Concat(From(input2)).Collect(ToSlice([]int{}))
	assert.Equal(t, want, got)
}

func TestStreamAppend(t *testing.T) {
	input := []int{1, 2, 3}
	want := []int{1, 2, 3, 4, 5}
	got := From(input).Append(4).Append(5).Collect(ToSlice([]int{}))

	assert.Equal(t, want, got)
}

func TestRange(t *testing.T) {
	want := []int{1, 2, 3, 4}
	got := Range(1, 5).Collect(ToSlice([]int{}))

	assert.Equal(t, want, got)
}

func TestRepeat(t *testing.T) {
	want := []int{1, 1, 1, 1, 1}
	got := Repeat(1, 5).Collect(ToSlice([]int{}))

	assert.Equal(t, want, got)
}

func TestDrop(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	want := []int{4, 5}
	got := From(input).Drop(3).Collect(ToSlice([]int{})).([]int)

	assert.Equal(t, want, got)
}

func TestLimit(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	want := []int{1, 2, 3}
	got := From(input).Limit(3).Collect(ToSlice([]int{})).([]int)

	assert.Equal(t, want, got)
}

func TestPeek(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	want := []int{1, 2, 3, 4, 5}
	var got []int

	From(input).Peek(func(it interface{}) {
		got = append(got, it.(int))
	}).Last()
	assert.Equal(t, want, got)
}
