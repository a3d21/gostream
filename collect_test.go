package gostream

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyCollectToSliceShouldBeNil(t *testing.T) {
	var input []int
	got := From(input).Collect(ToSlice([]int{})).([]int)
	assert.Nil(t, got)
	assert.Empty(t, got)
}

func TestEmptyCollectToSliceV2ShouldNotBeNil(t *testing.T) {
	var input []int
	got := From(input).CollectV2(ToSliceV2([]int{})).([]int)
	assert.Nil(t, got)
	assert.Empty(t, got)
}

func TestEmptyCollectToMapShouldNotBeNil(t *testing.T) {
	var input []int
	got := From(input).Collect(ToMap(map[int]int{}, func(v interface{}) interface{} {
		return v
	}, func(v interface{}) interface{} {
		return v
	})).(map[int]int)

	assert.NotNil(t, got)
	assert.Empty(t, got)
}

func TestEmptyCollectToMapV2ShouldNotBeNil(t *testing.T) {
	var input []int
	got := From(input).CollectV2(ToMapV2(map[int]int{}, func(v interface{}) interface{} {
		return v
	}, func(v interface{}) interface{} {
		return v
	})).(map[int]int)

	assert.NotNil(t, got)
	assert.Empty(t, got)
}

func TestMakeSlice(t *testing.T) {
	var a []int
	tt := reflect.TypeOf(a)
	t.Log(tt)
	ss := reflect.MakeSlice(tt, 0, 0)
	t.Log(ss)
	t.Log(ss.Interface())
	t.Log(ss.Type())

}

func TestCollectToSlice(t *testing.T) {
	intput := []int{1, 2, 3, 4, 5}
	got := From(intput).Collect(ToSlice([]int{}))

	if !reflect.DeepEqual(intput, got) {
		t.Errorf("%v != %v", got, intput)
	}
}

func TestCollectToSliceV2(t *testing.T) {
	intput := []int{1, 2, 3, 4, 5}
	got := From(intput).CollectV2(ToSliceV2([]int{}))

	if !reflect.DeepEqual(intput, got) {
		t.Errorf("%v != %v", got, intput)
	}
}

func TestCollectToMap(t *testing.T) {
	intput := []int{1, 2, 3, 4, 5}
	want := map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5}
	identity := func(it interface{}) interface{} { return it }
	got := From(intput).Collect(ToMap(map[int]int{}, identity, identity))

	if !reflect.DeepEqual(got, want) {
		t.Errorf("%v != %v", got, want)
	}
}

func TestCollectToMapV2(t *testing.T) {
	intput := []int{1, 2, 3, 4, 5}
	want := map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5}
	identity := func(it interface{}) interface{} { return it }
	got := From(intput).CollectV2(ToMapV2(map[int]int{}, identity, identity))

	if !reflect.DeepEqual(got, want) {
		t.Errorf("%v != %v", got, want)
	}
}

type Cargo struct {
	ID       int
	Name     string
	Location string
	Status   int
}

func TestCollectGroupBy(t *testing.T) {
	input := []*Cargo{{
		ID:       1,
		Name:     "foo",
		Location: "shenzhen",
		Status:   1,
	}, {
		ID:       2,
		Name:     "bar",
		Location: "shenzhen",
		Status:   0,
	}, {
		ID:       3,
		Name:     "a3d21",
		Location: "guangzhou",
		Status:   1,
	}}
	want := map[string][]*Cargo{
		"shenzhen": {{
			ID:       1,
			Name:     "foo",
			Location: "shenzhen",
			Status:   1,
		}, {
			ID:       2,
			Name:     "bar",
			Location: "shenzhen",
			Status:   0,
		}},
		"guangzhou": {{
			ID:       3,
			Name:     "a3d21",
			Location: "guangzhou",
			Status:   1,
		}},
	}

	getLocation := func(it interface{}) interface{} {
		return it.(*Cargo).Location
	}
	got := From(input).Collect(
		GroupBy(map[string][]*Cargo{}, getLocation,
			ToSlice([]*Cargo{}))).(map[string][]*Cargo)

	if !reflect.DeepEqual(want, got) {
		t.Errorf("%v != %v", got, want)
	}
}

func TestCollectGroupByV2(t *testing.T) {
	input := []*Cargo{{
		ID:       1,
		Name:     "foo",
		Location: "shenzhen",
		Status:   1,
	}, {
		ID:       2,
		Name:     "bar",
		Location: "shenzhen",
		Status:   0,
	}, {
		ID:       3,
		Name:     "a3d21",
		Location: "guangzhou",
		Status:   1,
	}}
	want := map[string][]*Cargo{
		"shenzhen": {{
			ID:       1,
			Name:     "foo",
			Location: "shenzhen",
			Status:   1,
		}, {
			ID:       2,
			Name:     "bar",
			Location: "shenzhen",
			Status:   0,
		}},
		"guangzhou": {{
			ID:       3,
			Name:     "a3d21",
			Location: "guangzhou",
			Status:   1,
		}},
	}

	getLocation := func(it interface{}) interface{} {
		return it.(*Cargo).Location
	}
	got := From(input).CollectV2(
		GroupByV2(map[string][]*Cargo{}, getLocation,
			ToSliceV2([]*Cargo{}))).(map[string][]*Cargo)

	if !reflect.DeepEqual(want, got) {
		t.Errorf("%v != %v", got, want)
	}
}

func TestMultiGroupBy(t *testing.T) {
	input := []*Cargo{{
		ID:       1,
		Name:     "foo",
		Location: "shenzhen",
		Status:   1,
	}, {
		ID:       2,
		Name:     "bar",
		Location: "shenzhen",
		Status:   0,
	}, {
		ID:       3,
		Name:     "a3d21",
		Location: "guangzhou",
		Status:   1,
	}}

	// group by status, city
	want := map[int]map[string][]*Cargo{
		1: {
			"shenzhen": {
				{
					ID:       1,
					Name:     "foo",
					Location: "shenzhen",
					Status:   1,
				},
			},
			"guangzhou": {
				{
					ID:       3,
					Name:     "a3d21",
					Location: "guangzhou",
					Status:   1,
				},
			},
		},
		0: {
			"shenzhen": {
				{
					ID:       2,
					Name:     "bar",
					Location: "shenzhen",
					Status:   0,
				},
			},
		},
	}

	getStatus := func(it interface{}) interface{} { return it.(*Cargo).Status }
	getLocation := func(it interface{}) interface{} { return it.(*Cargo).Location }

	// collect map, group by status,city
	// result type: map[int]map[string][]*Cargo
	got := From(input).Collect(
		GroupBy(map[int]map[string][]*Cargo(nil), getStatus,
			GroupBy(map[string][]*Cargo(nil), getLocation,
				ToSlice([]*Cargo(nil)))))

	if !reflect.DeepEqual(got, want) {
		t.Errorf("%v != %v", got, want)
	}
}

func TestMultiGroupByV2(t *testing.T) {
	input := []*Cargo{{
		ID:       1,
		Name:     "foo",
		Location: "shenzhen",
		Status:   1,
	}, {
		ID:       2,
		Name:     "bar",
		Location: "shenzhen",
		Status:   0,
	}, {
		ID:       3,
		Name:     "a3d21",
		Location: "guangzhou",
		Status:   1,
	}}

	// group by status, city
	want := map[int]map[string][]*Cargo{
		1: {
			"shenzhen": {
				{
					ID:       1,
					Name:     "foo",
					Location: "shenzhen",
					Status:   1,
				},
			},
			"guangzhou": {
				{
					ID:       3,
					Name:     "a3d21",
					Location: "guangzhou",
					Status:   1,
				},
			},
		},
		0: {
			"shenzhen": {
				{
					ID:       2,
					Name:     "bar",
					Location: "shenzhen",
					Status:   0,
				},
			},
		},
	}

	getStatus := func(it interface{}) interface{} { return it.(*Cargo).Status }
	getLocation := func(it interface{}) interface{} { return it.(*Cargo).Location }

	// collect map, group by status,city
	// result type: map[int]map[string][]*Cargo
	got := From(input).CollectV2(
		GroupByV2(map[int]map[string][]*Cargo(nil), getStatus,
			GroupByV2(map[string][]*Cargo(nil), getLocation,
				ToSliceV2([]*Cargo(nil)))))

	if !reflect.DeepEqual(got, want) {
		t.Errorf("%v != %v", got, want)
	}
}

func TestCollectorToMap(t *testing.T) {
	input := []*Cargo{{
		ID:       1,
		Name:     "foo",
		Location: "shenzhen",
		Status:   1,
	}, {
		ID:       2,
		Name:     "bar",
		Location: "shenzhen",
		Status:   0,
	}, {
		ID:       3,
		Name:     "a3d21",
		Location: "guangzhou",
		Status:   1,
	}}

	want := map[string]string{
		"foo":   "shenzhen",
		"bar":   "shenzhen",
		"a3d21": "guangzhou",
	}

	getLocation := func(it interface{}) interface{} { return it.(*Cargo).Location }
	getName := func(it interface{}) interface{} { return it.(*Cargo).Name }
	got := From(input).Collect(ToMap(map[string]string{}, getName, getLocation))

	if !reflect.DeepEqual(want, got) {
		t.Errorf("%v != %v", got, want)
	}
}

func TestCollectToSet(t *testing.T) {
	input := []int{1, 2, 3}
	want := map[int]bool{1: true, 2: true, 3: true}
	got := From(input).Collect(ToSet(map[int]bool{})).(map[int]bool)
	assert.Equal(t, want, got)
}

func TestCollectToSetV2(t *testing.T) {
	input := []int{1, 2, 3}
	want := map[int]bool{1: true, 2: true, 3: true}
	got := From(input).CollectV2(ToSetV2(map[int]bool{})).(map[int]bool)
	assert.Equal(t, want, got)
}

func TestFlatMap(t *testing.T) {
	input := [][]int{{3, 2, 1}, {6, 5, 4}, {9, 8, 7}}
	want := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	got := From(input).FlatMap(func(it interface{}) Stream {
		return From(it)
	}).SortedBy(func(it interface{}) interface{} {
		return it
	}).Collect(ToSlice([]int{}))

	if !reflect.DeepEqual(want, got) {
		t.Errorf("%v != %v", got, want)
	}
}

func TestCount(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	want := 5
	got := From(input).Collect(Count())

	if !reflect.DeepEqual(want, got) {
		t.Errorf("%v != %v", got, want)
	}
}

func TestGroupCount(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	want := map[bool]int{
		true:  3,
		false: 2,
	}

	got := From(input).Collect(GroupBy(want, func(v interface{}) interface{} {
		return v.(int) < 4
	}, Count()))
	if !reflect.DeepEqual(want, got) {
		t.Errorf("%v != %v", got, want)
	}
}

func TestCountV2(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	want := 5
	got := From(input).CollectV2(CountV2())

	if !reflect.DeepEqual(want, got) {
		t.Errorf("%v != %v", got, want)
	}
}

func TestGroupCountV2(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	want := map[bool]int{
		true:  3,
		false: 2,
	}

	got := From(input).CollectV2(GroupByV2(want, func(v interface{}) interface{} {
		return v.(int) < 4
	}, CountV2()))
	if !reflect.DeepEqual(want, got) {
		t.Errorf("%v != %v", got, want)
	}
}

func TestCustomAddCollector(t *testing.T) {
	input := []int{1, 2, 3}
	want := 1 + 2 + 3

	got := From(input).Collect(Collector(func() interface{} {
		return 0
	}, func(acc interface{}, item interface{}) interface{} {
		return acc.(int) + item.(int)
	}))
	assert.Equal(t, want, got)
}

func TestGroupSum(t *testing.T) {
	type AType struct {
		Name  string
		Count int
	}

	input := []AType{
		{"foo", 10},
		{"bar", 15},
		{"foo", 20},
		{"bar", 30},
	}
	want := map[string]int{"foo": 30, "bar": 45}
	got := From(input).Collect(GroupBy(map[string]int{},
		func(it interface{}) interface{} {
			return it.(AType).Name
		},
		Collector(func() interface{} {
			return 0
		}, func(acc interface{}, item interface{}) interface{} {
			return acc.(int) + item.(AType).Count
		}))).(map[string]int)
	assert.Equal(t, want, got)
}

func TestCustomAddCollectorV2(t *testing.T) {
	input := []int{1, 2, 3}
	want := 1 + 2 + 3

	got := From(input).CollectV2(CollectorV2(func() interface{} {
		return 0
	}, func(acc interface{}, item interface{}) interface{} {
		return acc.(int) + item.(int)
	}))
	assert.Equal(t, want, got)
}

func TestGroupSumV2(t *testing.T) {
	type AType struct {
		Name  string
		Count int
	}

	input := []AType{
		{"foo", 10},
		{"bar", 15},
		{"foo", 20},
		{"bar", 30},
	}
	want := map[string]int{"foo": 30, "bar": 45}
	got := From(input).CollectV2(GroupByV2(map[string]int{},
		func(it interface{}) interface{} {
			return it.(AType).Name
		},
		CollectorV2(func() interface{} {
			return 0
		}, func(acc interface{}, item interface{}) interface{} {
			return acc.(int) + item.(AType).Count
		}))).(map[string]int)
	assert.Equal(t, want, got)
}
