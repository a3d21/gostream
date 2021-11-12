# GoStream

gostream 是一个数据流式处理库。它可以声明式地对数据进行转换、过滤、排序、分组、收集，而无需关心操作细节。


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
$ go test -bench=.
goos: darwin
goarch: amd64
pkg: github.com/a3d21/gostream
BenchmarkToSliceRaw-12                      8216            132435 ns/op
BenchmarkToSliceStreamForeach-12             639           1819347 ns/op
BenchmarkCollectToSlice-12                    88          13006774 ns/op
BenchmarkCollectToSliceV2-12                 308           3847322 ns/op
BenchmarkLinqToSlice-12                      340           3513771 ns/op
BenchmarkToMapRaw-12                         184           6195945 ns/op
BenchmarkCollectToMap-12                      97          12378528 ns/op
BenchmarkCollectToMapV2-12                    98          12409819 ns/op
BenchmarkLinqToMap-12                        100          12293144 ns/op
BenchmarkGroupByRaw-12                        86          12673383 ns/op
BenchmarkGroupBy-12                           22          48737249 ns/op
BenchmarkGroupByV2-12                          9         116926833 ns/op
BenchmarkLinqGroupBy-12                       19          61038649 ns/op
BenchmarkPartition-12                        171           6949116 ns/op
BenchmarkCountRaw-12                         810           1468941 ns/op
BenchmarkCount-12                            170           7044074 ns/op
BenchmarkCountV2-12                          813           1468029 ns/op
BenchmarkGroupCount-12                        33          35121443 ns/op
BenchmarkGroupCountV2-12                      15          70305638 ns/op

```
