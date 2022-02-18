package gostream

import (
	"reflect"
	"sort"
	"testing"
	"testing/quick"

	"github.com/a3d21/gostream/gopark"
)

func TestSlice2MapSpec(t *testing.T) {
	assertion := func(vs []int) bool {
		m1 := From(vs).Map(func(it interface{}) interface{} {
			return KeyValue{
				Key:   it,
				Value: true,
			}
		}).Collect(ToMap(map[int]bool{})).(map[int]bool)
		m2 := gopark.Slice2Map(vs).(map[int]bool)

		return reflect.DeepEqual(m1, m2)
	}
	if err := quick.Check(assertion, &quick.Config{MaxCount: 2000}); err != nil {
		t.Error(err)
	}
}

func TestToSetSpec(t *testing.T) {
	assertion := func(vs []int) bool {
		m1 := From(vs).Collect(ToSet(map[int]bool{})).(map[int]bool)
		m2 := gopark.Slice2Map(vs).(map[int]bool)

		return reflect.DeepEqual(m1, m2)
	}
	if err := quick.Check(assertion, &quick.Config{MaxCount: 2000}); err != nil {
		t.Error(err)
	}
}

func TestKeysSpec(t *testing.T) {
	assertion := func(m map[string]int64) bool {
		s1 := From(m).Map(func(kv interface{}) interface{} {
			return kv.(KeyValue).Key
		}).SortedBy(identity).Collect(ToSlice([]string{})).([]string)

		s2 := gopark.Keys(m).([]string)
		sort.Strings(s2)
		return reflect.DeepEqual(s1, s2) || (len(s1) == 0 && len(s2) == 0)
	}

	if err := quick.Check(assertion, &quick.Config{MaxCount: 2000}); err != nil {
		t.Error(err)
	}
}

func TestValuesSpec(t *testing.T) {
	assertion := func(m map[string]int) bool {
		s1 := From(m).Map(func(kv interface{}) interface{} {
			return kv.(KeyValue).Value
		}).SortedBy(identity).Collect(ToSlice([]int{})).([]int)

		s2 := gopark.Values(m).([]int)
		sort.Ints(s2)
		return reflect.DeepEqual(s1, s2) || (len(s1) == 0 && len(s2) == 0)
	}

	if err := quick.Check(assertion, &quick.Config{MaxCount: 2000}); err != nil {
		t.Error(err)
	}
}
