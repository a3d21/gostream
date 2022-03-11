package gostream

import "testing"

func TestGetComparer(t *testing.T) {
	tests := []struct {
		x    interface{}
		y    interface{}
		want int
	}{
		{100, 500, -1},
		{-100, -500, 1},
		{256, 256, 0},
		{int8(100), int8(-100), 1},
		{int8(-100), int8(100), -1},
		{int8(100), int8(100), 0},
		{int16(100), int16(-100), 1},
		{int16(-100), int16(100), -1},
		{int16(100), int16(100), 0},
		{int32(100), int32(-100), 1},
		{int32(-100), int32(100), -1},
		{int32(100), int32(100), 0},
		{int64(100), int64(-100), 1},
		{int64(-100), int64(100), -1},
		{int64(100), int64(100), 0},
		{uint(100), uint(0), 1},
		{uint(0), uint(100), -1},
		{uint(100), uint(100), 0},
		{uint8(100), uint8(0), 1},
		{uint8(0), uint8(100), -1},
		{uint8(100), uint8(100), 0},
		{uint16(100), uint16(0), 1},
		{uint16(0), uint16(100), -1},
		{uint16(100), uint16(100), 0},
		{uint32(100), uint32(0), 1},
		{uint32(0), uint32(100), -1},
		{uint32(100), uint32(100), 0},
		{uint64(100), uint64(0), 1},
		{uint64(0), uint64(100), -1},
		{uint64(100), uint64(100), 0},
		{float32(5.), float32(1.), 1},
		{float32(1.), float32(5.), -1},
		{float32(0), float32(0), 0},
		{float64(5.), float64(1.), 1},
		{float64(1.), float64(5.), -1},
		{float64(0), float64(0), 0},
		{true, true, 0},
		{false, false, 0},
		{true, false, 1},
		{false, true, -1},
		{"foo", "foo", 0},
		{"foo", "bar", 1},
		{"bar", "foo", -1},
		{"FOO", "bar", -1},
		{foo{f1: 1}, foo{f1: 5}, -1},
		{foo{f1: 5}, foo{f1: 1}, 1},
		{foo{f1: 1}, foo{f1: 1}, 0},
	}

	for _, test := range tests {
		if r := getComparer(test.x)(test.x, test.y); r != test.want {
			t.Errorf("getComparer(%v)(%v,%v)=%v expected %v", test.x, test.x, test.y, r, test.want)
		}
	}
}

type foo struct {
	f1 int
	f2 bool
	f3 string
}

func (f foo) CompareTo(c Comparable) int {
	a, b := f.f1, c.(foo).f1

	if a < b {
		return -1
	} else if a > b {
		return 1
	}

	return 0
}

func TestComparableList_CompareTo(t *testing.T) {
	tests := []struct {
		name   string
		l1, l2 GTuple
		want   int
	}{
		{"nil = nil", nil, nil, 0},
		{"nil = {}", nil, GTuple{}, 0},
		{"{} < {1}", GTuple{}, GTuple{1}, -1},
		{"{1} = {1}", GTuple{1}, GTuple{1}, 0},
		{"{1,2} > {1,1}", GTuple{1, 2}, GTuple{1, 1}, 1},
		{"{1,a} < {1,b}", GTuple{1, "a"}, GTuple{1, "b"}, -1},
		{"multi type test", GTuple{1, int8(2), int16(3), uint(4), "5", true}, GTuple{1, int8(2), int16(3), uint(4), "5", true}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l1.CompareTo(tt.l2); got != tt.want {
				t.Errorf("CompareTo() = %v, want %v", got, tt.want)
			}
		})
	}
}
