package gopark

import (
	"reflect"
	"testing"
)

func TestPartitionBy(t *testing.T) {
	type args struct {
		vs   interface{}
		size int
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{"empty slice", args{[]int{}, 1}, [][]int{}},
		{"by 1", args{[]int{1, 2, 3}, 1}, [][]int{{1}, {2}, {3}}},
		{"by 2", args{[]int{1, 2, 3}, 2}, [][]int{{1, 2}, {3}}},
		{"by 3", args{[]int{1, 2, 3}, 3}, [][]int{{1, 2, 3}}},
		{"by 4", args{[]int{1, 2, 3}, 4}, [][]int{{1, 2, 3}}},
		{"by 5", args{[]int{1, 2, 3}, 5}, [][]int{{1, 2, 3}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PartitionBy(tt.args.vs, tt.args.size); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PartitionBy() = %v, want %v", got, tt.want)
			}
		})
	}
}
