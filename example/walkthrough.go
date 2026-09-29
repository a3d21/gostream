package main

import (
	"fmt"
	"reflect"

	. "github.com/a3d21/gostream/v2"
)

func walkthrough() {
	// work through

	// 1. Slice example
	// output:
	// 6,7,8,
	fmt.Println("example 01")
	input1 := []int{1, 2, 3, 4, 5, 6, 7, 8}
	From(input1).Filter(func(it int) bool { return it > 5 }).
		ForEach(func(it int) {
			fmt.Printf("%v,", it)
		})
	fmt.Println()

	// 2. Map example
	// output:
	// c => 3,d => 4,
	fmt.Println("example 02")
	input2 := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	FromMap(input2).Filter(func(it KeyValue[string, int]) bool { return it.Value > 2 }).
		ForEach(func(it KeyValue[string, int]) {
			fmt.Printf("%v => %v,", it.Key, it.Value)
		})
	fmt.Println()

	// 3. Collect to Slice
	input3 := []int{1, 3, 5, 7, 9}
	got3 := From(input3).ToSlice()
	assertEqual(input3, got3)

	// 4. Collect to Map by Code, result type: map[int64]*Cargo
	// Assume a collection of cargos
	cargos := []*Cargo{{
		Code:     1000,
		Location: "shenzhen",
		From:     "shenzhen",
		To:       "beijing",
		Created:  "2021-02-01 10:00:00",
	}, {
		Code:     1001,
		Location: "guangzhou",
		From:     "shenzhen",
		To:       "beijing",
		Created:  "2021-02-01 13:00:00",
	}, {
		Code:     1002,
		Location: "shanghai",
		From:     "guangzhou",
		To:       "beijing",
		Created:  "2021-02-01 13:00:00",
	}, {
		Code:     1003,
		Location: "beijing",
		From:     "beijing",
		To:       "guangzhou",
		Created:  "2021-02-01 13:00:00",
	}, {
		Code:     1004,
		Location: "beijing",
		From:     "shanghai",
		To:       "guangzhou",
		Created:  "2021-02-01 13:00:00",
	}}

	wantCargoMap := map[int64]*Cargo{
		1000: {
			Code:     1000,
			Location: "shenzhen",
			From:     "shenzhen",
			To:       "beijing",
			Created:  "2021-02-01 10:00:00",
		},
		1001: {
			Code:     1001,
			Location: "guangzhou",
			From:     "shenzhen",
			To:       "beijing",
			Created:  "2021-02-01 13:00:00",
		},
		1002: {
			Code:     1002,
			Location: "shanghai",
			From:     "guangzhou",
			To:       "beijing",
			Created:  "2021-02-01 13:00:00",
		},
		1003: {
			Code:     1003,
			Location: "beijing",
			From:     "beijing",
			To:       "guangzhou",
			Created:  "2021-02-01 13:00:00",
		},
		1004: {
			Code:     1004,
			Location: "beijing",
			From:     "shanghai",
			To:       "guangzhou",
			Created:  "2021-02-01 13:00:00",
		},
	}

	code2CargoMap := From(cargos).Collect(ToMapByCollector(
		func(it *Cargo) int64 { return it.Code },
		func(it *Cargo) *Cargo { return it },
	))
	assertEqual(code2CargoMap, wantCargoMap)

	// 5. Group by Location, result type: map[string][]*Cargo
	wantCargoByLocation := map[string][]*Cargo{
		"shenzhen": {{
			Code:     1000,
			Location: "shenzhen",
			From:     "shenzhen",
			To:       "beijing",
			Created:  "2021-02-01 10:00:00",
		}},
		"guangzhou": {{
			Code:     1001,
			Location: "guangzhou",
			From:     "shenzhen",
			To:       "beijing",
			Created:  "2021-02-01 13:00:00",
		}},
		"shanghai": {{
			Code:     1002,
			Location: "shanghai",
			From:     "guangzhou",
			To:       "beijing",
			Created:  "2021-02-01 13:00:00",
		}},
		"beijing": {{
			Code:     1003,
			Location: "beijing",
			From:     "beijing",
			To:       "guangzhou",
			Created:  "2021-02-01 13:00:00",
		}, {
			Code:     1004,
			Location: "beijing",
			From:     "shanghai",
			To:       "guangzhou",
			Created:  "2021-02-01 13:00:00",
		}},
	}
	cargoByLocation := From(cargos).Collect(
		GroupByCollector(
			func(it *Cargo) string { return it.Location },
			ToSliceCollector[*Cargo](),
		),
	)
	assertEqual(cargoByLocation, wantCargoByLocation)

	// 6. Group into nested map, result type: map[string]map[int64]*Cargo
	location2code2cargomap := From(cargos).Collect(
		GroupByCollector(
			func(it *Cargo) string { return it.Location },
			ToMapByCollector(
				func(it *Cargo) int64 { return it.Code },
				func(it *Cargo) *Cargo { return it },
			),
		),
	)

	// 7. Multi-level grouping by Location and To, result type: map[string]map[string][]*Cargo
	cargoByLocationByTo := From(cargos).Collect(
		GroupByCollector(
			func(it *Cargo) string { return it.Location },
			GroupByCollector(
				func(it *Cargo) string { return it.To },
				ToSliceCollector[*Cargo](),
			),
		),
	)

	// Flatten grouped maps back into slice
	// 8. Map[int64]*Cargo => []*Cargo, use FromMap + Map
	assertEqual(cargos,
		FromMap(code2CargoMap).Map(func(kv KeyValue[int64, *Cargo]) *Cargo {
			return kv.Value
		}).SortedBy(func(it *Cargo) int64 { return it.Code }).ToSlice())

	// 9. map[string][]*Cargo => []*Cargo, use FromMap + FlatMap
	assertEqual(cargos,
		FromMap(cargoByLocation).FlatMap(func(kv KeyValue[string, []*Cargo]) Stream[*Cargo] {
			return From(kv.Value)
		}).SortedBy(func(it *Cargo) int64 { return it.Code }).ToSlice())

	// 10. map[string]map[string][]*Cargo => []*Cargo, double FlatMap
	assertEqual(cargos,
		FromMap(cargoByLocationByTo).FlatMap(func(kv KeyValue[string, map[string][]*Cargo]) Stream[*Cargo] {
			return FromMap(kv.Value).FlatMap(func(kv2 KeyValue[string, []*Cargo]) Stream[*Cargo] {
				return From(kv2.Value)
			})
		}).SortedBy(func(it *Cargo) int64 { return it.Code }).ToSlice())

	// 11. map[string]map[int64]*Cargo => []*Cargo, FlatMap + Map
	assertEqual(cargos,
		FromMap(location2code2cargomap).FlatMap(func(kv KeyValue[string, map[int64]*Cargo]) Stream[*Cargo] {
			return FromMap(kv.Value).Map(func(kv2 KeyValue[int64, *Cargo]) *Cargo {
				return kv2.Value
			})
		}).SortedBy(func(it *Cargo) int64 { return it.Code }).ToSlice())
}

// Cargo
type Cargo struct {
	Code     int64
	Name     string
	Location string
	From     string
	To       string
	Created  string
}

func assertEqual(a, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		panic(fmt.Sprintf("%v != %v", a, b))
	}
}
