package gostream

import (
	"reflect"
	"sort"
	"testing"
	"testing/quick"

	"github.com/a3d21/gostream/gopark_deprecated"
)

func TestSlice2MapSpec(t *testing.T) {
	assertion := func(vs []int) bool {
		m1 := From(vs).Map(func(it interface{}) interface{} {
			return KeyValue{
				Key:   it,
				Value: true,
			}
		}).Collect(ToMap(map[int]bool{})).(map[int]bool)
		m2 := gopark_deprecated.Slice2Map(vs).(map[int]bool)

		return reflect.DeepEqual(m1, m2)
	}
	if err := quick.Check(assertion, &quick.Config{MaxCount: 2000}); err != nil {
		t.Error(err)
	}
}

func TestToSetSpec(t *testing.T) {
	assertion := func(vs []int) bool {
		m1 := From(vs).Collect(ToSet(map[int]bool{})).(map[int]bool)
		m2 := gopark_deprecated.Slice2Map(vs).(map[int]bool)

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

		s2 := gopark_deprecated.Keys(m).([]string)
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

		s2 := gopark_deprecated.Values(m).([]int)
		sort.Ints(s2)
		return reflect.DeepEqual(s1, s2) || (len(s1) == 0 && len(s2) == 0)
	}

	if err := quick.Check(assertion, &quick.Config{MaxCount: 2000}); err != nil {
		t.Error(err)
	}
}

/**
func TestMultiSortSpec(t *testing.T) {
	type foo struct {
		I    int
		I32  int32
		UI64 uint64
		F32  float32
		S    string
		B    bool
	}

	assertion := func(vs []foo) bool {

		vs = From(vs).Map(func(t interface{}) interface{} {
			f := t.(foo)
			return foo{
				I:    f.I % 20,
				I32:  f.I32 % 20,
				UI64: f.UI64 % 20,
				F32:  f.F32,
				S:    f.S,
				B:    f.B,
			}
		}).Collect(ToSlice([]foo{})).([]foo)

		got1 := From(vs).SortedBy(func(t interface{}) interface{} {
			f := t.(foo)
			return GTuple{f.I, f.I32, f.UI64, f.F32, f.S, f.B}
		}).Collect(ToSlice([]foo{})).([]foo)

		var got2 []foo
		linq.From(vs).OrderByT(func(f foo) interface{} {
			return f.I
		}).ThenByT(func(f foo) interface{} {
			return f.I32
		}).ThenByT(func(f foo) interface{} {
			return f.UI64
		}).ThenByT(func(f foo) interface{} {
			return f.F32
		}).ThenByT(func(f foo) interface{} {
			return f.S
		}).ThenByT(func(f foo) interface{} {
			return f.B
		}).ToSlice(&got2)

		return reflect.DeepEqual(got1, got2)
	}

	if err := quick.Check(assertion, &quick.Config{MaxCount: 5000}); err != nil {
		t.Error(err)
	}

}*/
