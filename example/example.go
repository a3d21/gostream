package main

import (
	"fmt"
	"reflect"

	. "github.com/a3d21/gostream/v2"
)

func main() {
	input := []int{4, 3, 2, 1}
	want := []int{6, 8}

	got := From(input).Map(func(it int) int {
		return 2 * it
	}).Filter(func(it int) bool {
		return it > 5
	}).SortedBy(func(it int) int {
		return it
	}).ToSlice()

	if !reflect.DeepEqual(got, want) {
		panic(fmt.Sprintf("%v != %v", got, want))
	}

	// walkthrough()
}
