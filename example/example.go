package main

import (
	"fmt"
	"reflect"

	. "github.com/a3d21/gostream"
)

func main() {
	input := []int{4, 3, 2, 1}
	want := []int{6, 8}

	got := From(input).Map(func(it interface{}) interface{} {
		return 2 * it.(int)
	}).Filter(func(it interface{}) bool {
		return it.(int) > 5
	}).SortedBy(func(it interface{}) interface{} {
		return it
	}).Collect(ToSlice([]int(nil)))

	if !reflect.DeepEqual(got, want) {
		panic(fmt.Sprintf("%v != %v", got, want))
	}

	// walkthrough()
}
