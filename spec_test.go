package gostream

import (
	"reflect"
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
		}).ToMap(map[int]bool{}).(map[int]bool)
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

func TestToSetV2Spec(t *testing.T) {
	assertion := func(vs []int) bool {
		m1 := From(vs).ToSet(map[int]bool{}).(map[int]bool)
		m2 := gopark.Slice2Map(vs).(map[int]bool)

		return reflect.DeepEqual(m1, m2)
	}
	if err := quick.Check(assertion, &quick.Config{MaxCount: 2000}); err != nil {
		t.Error(err)
	}
}
