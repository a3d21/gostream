package gostream

import (
	"testing"
	"testing/quick"
)

// TestPartitionMapDemo tests map partitioning and merging.
func TestPartitionMapDemo(t *testing.T) {

	assertion := func(src map[string]int) bool {
		if len(src) == 0 {
			return true
		}
		// Convert map to slice of KV pairs
		kvs := FromMap(src)

		// partition by size 3
		partitioned := Partition(kvs, 3).ToSlice()

		// merge partitions back
		merged := make(map[string]int)
		for _, part := range partitioned {
			for _, kv := range part {
				merged[kv.Key] = kv.Value
			}
		}

		if len(src) != len(merged) {
			return false
		}
		for k, v := range src {
			if merged[k] != v {
				return false
			}
		}
		return true
	}

	if err := quick.Check(assertion, &quick.Config{
		MaxCount: 2000,
	}); err != nil {
		t.Error(err)
	}
}
