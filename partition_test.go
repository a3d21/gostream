package gostream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPartition(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	want := [][]int{{1, 2, 3}, {4, 5}}

	got := Partition(From(input), 3).ToSlice()
	assert.Equal(t, want, got)
}

func TestPartitionSpec(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	got := Partition(From(input), 3).ToSlice()
	want := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10}}
	assert.Equal(t, want, got)
}
