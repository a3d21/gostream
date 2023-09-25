package gostream

import (
	"github.com/a3d21/gostream/gopark"
	"reflect"
	"testing"
	"testing/quick"
)

func TestPartition(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	want := [][]int{{1, 2, 3}, {4, 5}}

	got := From(input).Partition([]int(nil), 3).Collect(ToSlice([][]int(nil)))

	if !reflect.DeepEqual(got, want) {
		t.Errorf("%v != %v", got, want)
	}
}

func TestPartitionSpec(t *testing.T) {
	assertion := func(vs []int, usize uint) bool {
		s := int(usize%20 + 1)
		v1 := From(vs).Partition([]int(nil), s).Collect(ToSlice(([][]int(nil)))).([][]int)
		v2 := gopark.PartitionBy(vs, s)
		return reflect.DeepEqual(v1, v2) || (len(v1) == 0 && len(v2) == 0)
	}

	if err := quick.Check(assertion, &quick.Config{MaxCount: 2000}); err != nil {
		t.Error(err)
	}
}
