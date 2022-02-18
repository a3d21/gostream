package gopark

import (
	"reflect"
	"testing"
)

func TestKeys(t *testing.T) {
	type args struct {
		m interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{"nil map", args{map[string]int(nil)}, []string{}},
		{"empty map", args{map[string]int{}}, []string{}},
		{"string-int map", args{map[string]int{"foo": 1}}, []string{"foo"}},
		{"int64-bool map", args{map[int64]bool{1: true}}, []int64{1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Keys(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValues(t *testing.T) {
	type args struct {
		m interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{"nil map", args{map[string]int(nil)}, []int{}},
		{"empty map", args{map[string]int{}}, []int{}},
		{"string-int map", args{map[string]int{"foo": 1}}, []int{1}},
		{"int64-bool map", args{map[int64]bool{1: true}}, []bool{true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Values(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Values() = %v, want %v", got, tt.want)
			}
		})
	}
}
