package gostream

import (
	"reflect"
	"testing"
	"testing/quick"
)

// TestPartitionMapDemo Map分区、合并测试
func TestPartitionMapDemo(t *testing.T) {

	assertion := func(src map[string]int, usize uint) bool {
		s := int(usize%20 + 1)
		// partition by size
		par := From(src).Partition([]KeyValue{}, s).Map(func(it interface{}) interface{} {
			return From(it).Collect(ToMap(map[string]int{}))
		}).Collect(ToSlice([]map[string]int{})).([]map[string]int)

		// merge partition
		merged := From(par).FlatMap(func(it interface{}) Stream {
			return From(it)
		}).Collect(ToMap(map[string]int{})).(map[string]int)

		return reflect.DeepEqual(src, merged)
	}

	if err := quick.Check(assertion, &quick.Config{
		MaxCount: 2000,
	}); err != nil {
		t.Error(err)
	}
}
