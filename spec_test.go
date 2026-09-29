package gostream

import (
	"github.com/a3d21/gostream/v2/gopark"
	"sort"
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/assert"
)

func TestSlice2MapSpec(t *testing.T) {
	assertion := func(vs []int) bool {
		m1 := From(vs).Collect(ToSetCollector[int]())
		m2 := gopark.Slice2Map(vs)

		if len(m1) != len(m2) {
			return false
		}
		for k, v := range m1 {
			if m2[k] != v {
				return false
			}
		}
		return true
	}
	if err := quick.Check(assertion, &quick.Config{MaxCount: 2000}); err != nil {
		t.Error(err)
	}
}

func TestToSetSpec(t *testing.T) {
	assertion := func(vs []int) bool {
		m1 := From(vs).Collect(ToSetCollector[int]())
		m2 := gopark.Slice2Map(vs)

		if len(m1) != len(m2) {
			return false
		}
		for k, v := range m1 {
			if m2[k] != v {
				return false
			}
		}
		return true
	}
	if err := quick.Check(assertion, &quick.Config{MaxCount: 2000}); err != nil {
		t.Error(err)
	}
}

func TestKeysSpec(t *testing.T) {
	assertion := func(m map[string]int64) bool {
		kvs := FromMap(m)
		s1 := kvs.Map(func(kv KeyValue[string, int64]) string {
			return kv.Key
		}).SortedBy(func(s string) string { return s }).ToSlice()

		s2 := gopark.Keys(m)
		sort.Strings(s2)
		return assert.ObjectsAreEqual(s1, s2) || (len(s1) == 0 && len(s2) == 0)
	}

	if err := quick.Check(assertion, &quick.Config{MaxCount: 2000}); err != nil {
		t.Error(err)
	}
}

func TestValuesSpec(t *testing.T) {
	assertion := func(m map[string]int) bool {
		kvs := FromMap(m)
		s1 := kvs.Map(func(kv KeyValue[string, int]) int {
			return kv.Value
		}).SortedBy(func(i int) int { return i }).ToSlice()

		s2 := gopark.Values(m)
		sort.Ints(s2)
		return assert.ObjectsAreEqual(s1, s2) || (len(s1) == 0 && len(s2) == 0)
	}

	if err := quick.Check(assertion, &quick.Config{MaxCount: 2000}); err != nil {
		t.Error(err)
	}
}
