package gostream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Tuple2[L any, R any] struct {
	Left  L
	Right R
}

type Tuple3[A any, B any, C any] struct {
	First  A
	Second B
	Third  C
}

func TestZip2(t *testing.T) {
	input1 := []int{1, 2, 3, 4}
	input2 := []string{"a", "b", "c"}
	want := []Tuple2[int, string]{{1, "a"}, {2, "b"}, {3, "c"}}

	got := Zip2By(From(input1), From(input2), func(l int, r string) Tuple2[int, string] {
		return Tuple2[int, string]{Left: l, Right: r}
	}).ToSlice()
	assert.Equal(t, want, got)
}

func TestZip3(t *testing.T) {
	input1 := []int{1, 2, 3, 4}
	input2 := []string{"a", "b", "c"}
	input3 := []string{"foo", "bar"}
	want := []Tuple3[int, string, string]{{1, "a", "foo"}, {2, "b", "bar"}}

	got := Zip3By(From(input1), From(input2), From(input3), func(a int, b string, c string) Tuple3[int, string, string] {
		return Tuple3[int, string, string]{First: a, Second: b, Third: c}
	}).ToSlice()
	assert.Equal(t, want, got)
}
