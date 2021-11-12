package main

import (
	"fmt"
	"reflect"

	. "github.com/a3d21/gostream"
)

func walkthrough() {
	// work through

	// 1. Slice example
	// output:
	// 6,7,8,
	fmt.Println("example 01")
	input1 := []int{1, 2, 3, 4, 5, 6, 7, 8}
	From(input1).Filter(func(it interface{}) bool { return it.(int) > 5 }).
		ForEach(func(it interface{}) {
			fmt.Printf("%v,", it)
		})
	fmt.Println()

	// 2. Map example
	// output:
	// c => 3,d => 4,
	fmt.Println("example 02")
	input2 := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	From(input2).Filter(func(it interface{}) bool { return it.(KeyValue).Value.(int) > 2 }).
		ForEach(func(it interface{}) {
			v := it.(KeyValue) // type: KeyValue
			fmt.Printf("%v => %v,", v.Key, v.Value)
		})
	fmt.Println()

	// 3. Collect to Slice
	// got3，anotherGot3 类型不同
	input3 := []int{1, 3, 5, 7, 9}
	got3 := From(input3).Collect(ToSlice([]int(nil))) // 说明：受限于Go类型系统，需要显示传类型参数`[]int(nil)`，即[]int类型的nil
	anotherGot4 := From(input3).Collect(ToSlice([]interface{}(nil)))
	_ = got3.([]int)
	_ = anotherGot4.([]interface{})
	assertEqual(input3, got3)

	// 4. Collect to Map by Code, result type: map[int64]*Cargo
	// 假设一组货物
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

	getCode := func(it interface{}) interface{} { return it.(*Cargo).Code }
	identity := func(it interface{}) interface{} { return it }
	code2CargoMap := From(cargos).Collect(ToMap(map[int64]*Cargo(nil), getCode, identity)).(map[int64]*Cargo)
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
	getLocation := func(it interface{}) interface{} { return it.(*Cargo).Location }
	cargoByLocation := From(cargos).Collect(
		GroupBy(map[string][]*Cargo(nil), getLocation,
			ToSlice([]*Cargo(nil))))

	assertEqual(cargoByLocation, wantCargoByLocation)

	// 6. 分组成Map result type: map[string]map[int64]*Cargo
	location2code2cargomap := From(cargos).Collect(
		GroupBy(map[string]map[int64]*Cargo(nil), getLocation,
			ToMap(map[int64]*Cargo(nil), getCode, identity)))

	// 7. 多重分组，by Location, To, result Type: map[string]map[string][]*Cargo
	getTo := func(it interface{}) interface{} { return it.(*Cargo).To }
	cargoByLocationByTo := From(cargos).Collect(
		GroupBy(map[string]map[string][]*Cargo(nil), getLocation,
			GroupBy(map[string][]*Cargo(nil), getTo,
				ToSlice([]*Cargo(nil)))))

	// okok, 看了这么多分组Collect后，让我们把数据转换成Slice吧
	// 8. Map[int64]*Cargo => []*Cargo, use Map
	assertEqual(cargos,
		From(code2CargoMap).Map(func(it interface{}) interface{} {
			return it.(KeyValue).Value
		}).SortedBy(getCode).Collect(ToSlice([]*Cargo{})))

	// 9. map[string][]*Cargo => []*Cargo， use FlatMap
	assertEqual(cargos,
		From(cargoByLocation).FlatMap(func(it interface{}) Stream {
			return From(it.(KeyValue).Value)
		}).SortedBy(getCode).Collect(ToSlice([]*Cargo{})))

	// 10. map[string]map[string][]*Cargo => []*Cargo, double FlatMap
	assertEqual(cargos,
		From(cargoByLocationByTo).FlatMap(func(it interface{}) Stream {
			return From(it.(KeyValue).Value).FlatMap(func(it2 interface{}) Stream {
				return From(it2.(KeyValue).Value)
			})
		}).SortedBy(getCode).Collect(ToSlice([]*Cargo{})))

	// 11. map[string]map[int64]*Cargo => []*Cargo, FlatMap + Map
	assertEqual(cargos,
		From(location2code2cargomap).FlatMap(func(it interface{}) Stream {
			return From(it.(KeyValue).Value).Map(func(it2 interface{}) interface{} {
				return it2.(KeyValue).Value
			})
		}).SortedBy(getCode).Collect(ToSlice([]*Cargo{})))
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
