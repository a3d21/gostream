package gostream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyCollectToSliceShouldBeNil(t *testing.T) {
	var input []int
	got := From(input).ToSlice()
	assert.Nil(t, got)
	assert.Empty(t, got)
}

func TestEmptyCollectToMapShouldNotBeNil(t *testing.T) {
	var input []int
	got := From(input).Collect(ToMapByCollector(
		func(v int) int { return v },
		func(v int) int { return v },
	))

	assert.NotNil(t, got)
	assert.Empty(t, got)
}

func TestCollectToSlice(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	got := From(input).ToSlice()
	assert.Equal(t, input, got)
}

func TestCollectToMap(t *testing.T) {
	input := []int{1, 2, 3}
	want := map[int]bool{1: true, 2: true, 3: true}
	got := From(input).Map(func(it int) KeyValue[int, bool] {
		return KeyValue[int, bool]{
			Key:   it,
			Value: true,
		}
	}).Collect(ToMapCollector[int, bool]())
	assert.Equal(t, want, got)
}

func TestCollectToMapBy(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	want := map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5}
	got := From(input).Collect(ToMapByCollector(
		func(it int) int { return it },
		func(it int) int { return it },
	))
	assert.Equal(t, want, got)
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

	got := From(input).Collect(
		GroupByCollector(
			func(it *Cargo) string { return it.Location },
			ToSliceCollector[*Cargo](),
		),
	)

	assert.Equal(t, want, got)
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

	got := From(input).Collect(
		GroupByCollector(
			func(it *Cargo) int { return it.Status },
			GroupByCollector(
				func(it *Cargo) string { return it.Location },
				ToSliceCollector[*Cargo](),
			),
		),
	)

	assert.Equal(t, want, got)
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

	got := From(input).Collect(ToMapByCollector(
		func(it *Cargo) string { return it.Name },
		func(it *Cargo) string { return it.Location },
	))
	assert.Equal(t, want, got)
}

func TestCollectToSet(t *testing.T) {
	input := []int{1, 2, 3}
	want := map[int]bool{1: true, 2: true, 3: true}
	got := From(input).Collect(ToSetCollector[int]())
	assert.Equal(t, want, got)
}

func TestFlatMap(t *testing.T) {
	input := [][]int{{3, 2, 1}, {6, 5, 4}, {9, 8, 7}}
	want := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	got := From(input).FlatMap(func(it []int) Stream[int] {
		return From(it)
	}).SortedBy(func(it int) int {
		return it
	}).ToSlice()
	assert.Equal(t, want, got)
}

func TestCount(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	want := 5
	got := From(input).Count()
	assert.Equal(t, want, got)
}

func TestGroupCount(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	want := map[bool]int{
		true:  3,
		false: 2,
	}

	got := From(input).Collect(GroupByCollector(
		func(v int) bool { return v < 4 },
		CountCollector[int](),
	))
	assert.Equal(t, want, got)
}

func TestCustomAddCollector(t *testing.T) {
	input := []int{1, 2, 3}
	want := 1 + 2 + 3

	got := From(input).Collect(CollectBy(
		func() int { return 0 },
		func(acc int, item int) int { return acc + item },
	))
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
	got := From(input).Collect(GroupByCollector(
		func(it AType) string { return it.Name },
		CollectBy(
			func() int { return 0 },
			func(acc int, item AType) int { return acc + item.Count },
		),
	))
	assert.Equal(t, want, got)
}
