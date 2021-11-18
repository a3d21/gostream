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
BenchmarkToSliceRaw-8                       9064            130558 ns/op
BenchmarkToSliceStreamForeach-8              724           1729650 ns/op
BenchmarkCollectToSlice-8                    100          10976960 ns/op
BenchmarkCollectToSliceV2-8                  356           3442700 ns/op
BenchmarkLinqToSlice-8                       369           3233192 ns/op
BenchmarkToMapRaw-8                          169           7029026 ns/op
BenchmarkCollectToMap-8                       96          12598103 ns/op
BenchmarkCollectToMapV2-8                     87          12358812 ns/op
BenchmarkLinqToMap-8                          84          12632280 ns/op
BenchmarkToSetRaw-8                          198           5979733 ns/op
BenchmarkCollectToSet-8                       91          11230164 ns/op
BenchmarkCollectToSetV2-8                    100          11303409 ns/op
BenchmarkLinqToSet-8                          96          11588965 ns/op
BenchmarkGroupByRaw-8                         90          13014776 ns/op
BenchmarkGroupBy-8                            22          47840903 ns/op
BenchmarkGroupByV2-8                           8         126779306 ns/op
BenchmarkLinqGroupBy-8                        18          63476018 ns/op
BenchmarkPartition-8                         198           6084481 ns/op
BenchmarkCountRaw-8                          850           1379226 ns/op
BenchmarkCount-8                             196           6036287 ns/op
BenchmarkCountV2-8                           865           1384256 ns/op
BenchmarkGroupCount-8                         31          34496447 ns/op
BenchmarkGroupCountV2-8                       14          87656078 ns/op
PASS
ok      github.com/a3d21/gostream       34.211s
```
