# GoStream

gostream 是一个数据流式处理库。它可以声明式地对数据进行转换、过滤、排序、分组、收集，而无需关心操作细节。

## Changelog
2021-11-18
- add ToSet() collector

## Roadmap
- [ ] 移除go-linq依赖

## Get GoStream

```
go get github.com/a3d21/gostream
```

## Example

See [walkthrough.go](./example/walkthrough.go)

### Base Example

```go
package main

import (
	"fmt"
	"reflect"

	. "github.com/a3d21/gostream/core"
)

func main() {
	input := []int{4, 3, 2, 1}
	want := []int{6, 8}

	got := From(input).Map(func(it interface{}) interface{} {
		return 2 * it.(int)
	}).Filter(func(it interface{}) bool {
		return it.(int) > 5
	}).SortedBy(func(it interface{}) interface{} {
		return it
	}).Collect(ToSlice([]int(nil)))

	if !reflect.DeepEqual(got, want) {
		panic(fmt.Sprintf("%v != %v", got, want))
	}

	// walkthrough()
}
```


### Map & FlatMap

Map和FlatMap的差别在于：FlatMap的mapper返回一个Stream。

```go
input := [][]int{{3, 2, 1}, {6, 5, 4}, {9, 8, 7}}
want := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
got := From(input).FlatMap(func(it interface{}) Stream {
   return From(it)
}).SortedBy(func(it interface{}) interface{} {
   return it
}).Collect(ToSlice([]int{}))
```

### Collect ToSlice & ToMap

Collect将数据收集起来。受限于go的类型系统，需要显示传类型参数——一个目标类型的实例，可以为nil。

```go
intput := []int{1, 2, 3, 4, 5}
identity := func(it interface{}) interface{} { return it }

// []int{1, 2, 3, 4, 5}
gotSlice := From(intput).Collect(ToSlice([]int{}))
// map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5}
gotMap := From(intput).Collect(ToMap(map[int]int(nil), identity, identity))
```

### Collect GroupBy
GroupBy定义一个分组收集器，参数依序分别为 类型参数、分类方法、下游收集器。
GroupBy可以和ToSlice、ToMap组合，GroupBy也可以多级嵌套，实现多级分组。
```go
GroupBy(typ interface{}, classifier normalFn, downstream collector) collector
```

假设一组货物，需要按Status,Location进行分组，目标类型为 map[int]map[string][]*Cargo。

```go
// Cargo 货物实体
type Cargo struct {
   ID       int
   Name     string
   Location string
   Status   int
}

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

```

```go
getStatus := func(it interface{}) interface{} { return it.(*Cargo).Status }
getLocation := func(it interface{}) interface{} { return it.(*Cargo).Location }

// result type: map[int]map[string][]*Cargo
got := From(input).Collect(
   GroupBy(map[int]map[string][]*Cargo(nil), getStatus,
      GroupBy(map[string][]*Cargo(nil), getLocation,
         ToSlice([]*Cargo(nil)))))
```

### Flatten Group

这个示例演示如何将多级分组Map转成Slice。`map[int]map[string][]*Cargo => []*Cargo`

```go
From(cargoByLocationByTo).FlatMap(func(it interface{}) Stream {
   return From(it.(KeyValue).Value).FlatMap(func(it2 interface{}) Stream {
      return From(it2.(KeyValue).Value)
   })
}).Collect(ToSlice([]*Cargo{}))
```


## Benchmark
```
$ go test -bench .
goos: darwin
goarch: amd64
pkg: github.com/a3d21/gostream
cpu: Intel(R) Core(TM) i7-1068NG7 CPU @ 2.30GHz
BenchmarkToSliceRaw-8                       8953            135354 ns/op
BenchmarkToSliceStreamForeach-8              727           1646809 ns/op
BenchmarkCollectToSlice-8                    100          11056603 ns/op
BenchmarkCollectToSliceV2-8                  343           3488931 ns/op
BenchmarkLinqToSlice-8                       370           3218450 ns/op
BenchmarkToMapRaw-8                          165           7093403 ns/op
BenchmarkCollectToMap-8                      207           5789365 ns/op
BenchmarkCollectToMapV2-8                    218           5541738 ns/op
BenchmarkLinqToMap-8                         223           5376933 ns/op
BenchmarkToSetRaw-8                          196           6045459 ns/op
BenchmarkCollectToSet-8                      100          11328736 ns/op
BenchmarkCollectToSetV2-8                    100          11786037 ns/op
BenchmarkLinqToSet-8                          87          11689963 ns/op
BenchmarkGroupByRaw-8                        487           2366660 ns/op
BenchmarkGroupBy-8                            42          28509248 ns/op
BenchmarkGroupByV2-8                         100          11627525 ns/op
BenchmarkLinqGroupBy-8                       153           7730886 ns/op
BenchmarkPartition-8                         196           6017234 ns/op
BenchmarkCountRaw-8                          865           1374016 ns/op
BenchmarkCount-8                             198           6071583 ns/op
BenchmarkCountV2-8                           868           1378172 ns/op
BenchmarkGroupCount-8                         61          19945508 ns/op
BenchmarkGroupCountV2-8                      130           9101270 ns/op
BenchmarkSumRaw-8                          40366             30390 ns/op
BenchmarkCustomSumCollector-8                189           6361963 ns/op
BenchmarkCustomSumCollectorV2-8              435           2735501 ns/op
BenchmarkGroupSumRaw-8                      1059           1128347 ns/op
BenchmarkGroupSum-8                           58          20542984 ns/op
BenchmarkGroupSumV2-8                        100          10555309 ns/op
PASS
ok      github.com/a3d21/gostream       44.039s

```
