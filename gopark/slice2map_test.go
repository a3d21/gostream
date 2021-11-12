package gopark

import (
	"reflect"
	"testing"
)

func TestSlice2Map(t *testing.T) {
	type args struct {
		vs interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{"nil", args{vs: []int(nil)}, map[int]bool{}},
		{"i64", args{vs: []int64{1, 2, 3}}, map[int64]bool{1: true, 2: true, 3: true}},
		{"string", args{vs: []string{"1", "2", "3"}}, map[string]bool{"1": true, "2": true, "3": true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Slice2Map(tt.args.vs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Slice2Map() = %v, want %v", got, tt.want)
			}
		})
	}
}
