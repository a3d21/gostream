package gopark

import (
	"github.com/a3d21/gostream/gopark_deprecated"
	"reflect"
	"sort"
	"testing"
	"testing/quick"
)

func TestPartitionSpec(t *testing.T) {
	assertion := func(vs []int, usize uint) bool {
		size := int(usize%20 + 1)
		got := PartitionBy(vs, size)
		want := gopark_deprecated.PartitionBy(vs, size).([][]int)
		return reflect.DeepEqual(want, got) || (len(got) == 0 && len(want) == 0)
	}
	if err := quick.Check(assertion, &quick.Config{MaxCount: 2000}); err != nil {
		t.Error(err)
	}
}

func TestKeysSpec(t *testing.T) {
	assertion := func(m map[int]int) bool {
		got := Keys(m)
		want := gopark_deprecated.Keys(m).([]int)
		sort.Ints(got)
		sort.Ints(want)
		return reflect.DeepEqual(want, got) || (len(got) == 0 && len(want) == 0)
	}
	if err := quick.Check(assertion, &quick.Config{MaxCount: 2000}); err != nil {
		t.Error(err)
	}
}

func TestValuesSpec(t *testing.T) {
	assertion := func(m map[int]int) bool {
		got := Values(m)
		want := gopark_deprecated.Values(m).([]int)
		sort.Ints(got)
		sort.Ints(want)
		return reflect.DeepEqual(want, got) || (len(got) == 0 && len(want) == 0)
	}
	if err := quick.Check(assertion, &quick.Config{MaxCount: 2000}); err != nil {
		t.Error(err)
	}
}

func TestSlice2MapSpec(t *testing.T) {
	assertion := func(vs []int) bool {
		got := Slice2Map(vs)
		want := gopark_deprecated.Slice2Map(vs).(map[int]bool)
		return reflect.DeepEqual(want, got) || (len(got) == 0 && len(want) == 0)
	}
	if err := quick.Check(assertion, &quick.Config{MaxCount: 2000}); err != nil {
		t.Error(err)
	}
}
