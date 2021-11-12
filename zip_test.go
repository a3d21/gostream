package gostream

import (
	"reflect"
	"testing"
)

type Tuple2 struct {
	Left, Right interface{}
}

type Tuple3 struct {
	First, Second, Third interface{}
}

var tuple2ZipFn zip2Fn = func(left, right interface{}) interface{} {
	return Tuple2{Left: left, Right: right}
}

var tuple3ZipFn zip3Fn = func(first, second, third interface{}) interface{} {
	return Tuple3{First: first, Second: second, Third: third}
}

func TestZip2(t *testing.T) {
	input1 := []int{1, 2, 3, 4}
	input2 := []string{"a", "b", "c"}
	want := []Tuple2{{1, "a"}, {2, "b"}, {3, "c"}}

	got := Zip2By(From(input1), From(input2), tuple2ZipFn).CollectV2(ToSliceV2([]Tuple2(nil)))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%v != %v", got, want)
	}
}

func TestZip3(t *testing.T) {
	input1 := []int{1, 2, 3, 4}
	input2 := []string{"a", "b", "c"}
	input3 := []string{"foo", "bar"}
	want := []Tuple3{{1, "a", "foo"}, {2, "b", "bar"}}

	got := Zip3By(From(input1), From(input2), From(input3), tuple3ZipFn).CollectV2(ToSliceV2([]Tuple3(nil)))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%v != %v", got, want)
	}
}
