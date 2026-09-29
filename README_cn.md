# GoStream

[English](./README.md) | 简体中文

gostream 是一个数据流式处理库。它可以声明式地对数据进行转换、过滤、排序、分组、收集，而无需关心操作细节。

## Changelog

2026-09-29

- 支持 Go 泛型，重构为强类型安全流式处理 API
- 升级模块路径为 `github.com/a3d21/gostream/v2`
- 简化 Collector 设计，移除 dummy 实例传参

2023-09-24

- remove go-linq

2021-11-27

- upgrade collector to v2

2021-11-18

- add ToSet() collector

## Get GoStream

```
go get github.com/a3d21/gostream/v2
```

## Example

See [walkthrough.go](./example/walkthrough.go)

### Base Example

```go
package main

import (
	"fmt"
	"reflect"

	. "github.com/a3d21/gostream/v2"
)

func main() {
	input := []int{4, 3, 2, 1}
	want := []int{6, 8}

	got := From(input).Map(func(it int) int {
		return 2 * it
	}).Filter(func(it int) bool {
		return it > 5
	}).SortedBy(func(it int) int {
		return it
	}).ToSlice()

	if !reflect.DeepEqual(got, want) {
		panic(fmt.Sprintf("%v != %v", got, want))
	}

	// walkthrough()
}
```

### Map & FlatMap

`Map` 将流中的元素进行一对一转换；`FlatMap` 的 mapper 则返回一个 `Stream`，并将嵌套流扁平化展开。

```go
input := [][]int{{3, 2, 1}, {6, 5, 4}, {9, 8, 7}}
want := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

got := From(input).FlatMap(func(it []int) Stream[int] {
	return From(it)
}).SortedBy(func(it int) int {
	return it
}).ToSlice()
```

### FromMap (键值对流)

支持直接从 `map` 创建流，元素类型为 `KeyValue[K, V]`：

```go
input := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}

FromMap(input).
	Filter(func(it KeyValue[string, int]) bool { return it.Value > 2 }).
	ForEach(func(it KeyValue[string, int]) {
		fmt.Printf("%v => %v\n", it.Key, it.Value)
	})
```

### Collect ToSlice, ToMap & ToSet

在泛型支持下，收集操作具备完整的编译期类型安全，无需再传入 dummy 实例或执行类型断言：

- **ToSlice**：可直接调用 `.ToSlice()`，亦可通过 `.Collect(ToSliceCollector[T]())` 收集。
- **ToMap**：通过 `ToMapByCollector(keyMapper, valueMapper)` 自定义键值提取规则；若流中元素本身为 `KeyValue[K, V]`，可直接使用 `ToMapCollector[K, V]()`。
- **ToSet**：通过 `ToSetCollector[T]()` 收集为 `map[T]bool`。

```go
input := []int{1, 2, 3, 4, 5}

// 收集为切片: []int{1, 2, 3, 4, 5}
gotSlice := From(input).ToSlice()

// 收集为映射: map[int]int{1: 10, 2: 20, 3: 30, 4: 40, 5: 50}
gotMap := From(input).Collect(ToMapByCollector(
	func(it int) int { return it },
	func(it int) int { return it * 10 },
))

// 收集为集合: map[int]bool{1: true, 2: true, 3: true, 4: true, 5: true}
gotSet := From(input).Collect(ToSetCollector[int]())
```

### Collect GroupBy

`GroupByCollector` 定义分组收集器，参数依序为分类函数 `classifier` 与下游收集器 `downstream`。
得益于 Go 泛型类型推断，多数情况下无需显式指定类型参数。支持与 `ToSliceCollector`、`ToMapByCollector`、`CountCollector` 等自由组合，并支持多级嵌套分组。

```go
GroupByCollector[T any, K comparable, R any](classifier func(T) K, downstream Collector[T, R]) Collector[T, map[K]R]
```

假设一组货物：

```go
// Cargo 货物实体
type Cargo struct {
	ID       int
	Name     string
	Location string
	Status   int
}

input := []*Cargo{
	{ID: 1, Name: "foo", Location: "shenzhen", Status: 1},
	{ID: 2, Name: "bar", Location: "shenzhen", Status: 0},
	{ID: 3, Name: "a3d21", Location: "guangzhou", Status: 1},
}
```

单级分组与多级分组：

```go
// 1. 单级分组: 按 Location 分组 => map[string][]*Cargo
cargoByLocation := From(input).Collect(
	GroupByCollector(
		func(it *Cargo) string { return it.Location },
		ToSliceCollector[*Cargo](),
	),
)

// 2. 多级嵌套分组: 先按 Status 分组，再按 Location 分组 => map[int]map[string][]*Cargo
cargoByStatusByLocation := From(input).Collect(
	GroupByCollector(
		func(it *Cargo) int { return it.Status },
		GroupByCollector(
			func(it *Cargo) string { return it.Location },
			ToSliceCollector[*Cargo](),
		),
	),
)
```

### Flatten Group

结合 `FromMap` 与 `FlatMap`，可将多级分组 Map 结构展开回切片：`map[int]map[string][]*Cargo => []*Cargo`。
整个过程强类型安全，无任何运行时类型断言：

```go
cargos := FromMap(cargoByStatusByLocation).FlatMap(func(kv KeyValue[int, map[string][]*Cargo]) Stream[*Cargo] {
	return FromMap(kv.Value).FlatMap(func(kv2 KeyValue[string, []*Cargo]) Stream[*Cargo] {
		return From(kv2.Value)
	})
}).SortedBy(func(it *Cargo) int {
	return it.ID
}).ToSlice()
```

## Benchmark

```
$ go test -bench .
goos: darwin
goarch: amd64
pkg: github.com/a3d21/gostream
cpu: Intel(R) Core(TM) i7-7820HQ CPU @ 2.90GHz
BenchmarkToSliceRaw-8                       7677            140295 ns/op
BenchmarkToSliceStreamForeach-8              607           1895616 ns/op
BenchmarkCollectToSlice-8                    296           4038376 ns/op
BenchmarkToMapRaw-8                          164           7141694 ns/op
BenchmarkCollectToMap-8                      156           7556103 ns/op
BenchmarkToSetRaw-8                          177           6553303 ns/op
BenchmarkCollectToSet-8                       92          12908673 ns/op
BenchmarkGroupByRaw-8                        397           2951509 ns/op
BenchmarkGroupBy-8                            91          12943394 ns/op
BenchmarkPartition-8                         146           8094362 ns/op
BenchmarkCountRaw-8                          746           1527550 ns/op
BenchmarkCount-8                             720           1533764 ns/op
BenchmarkGroupCount-8                        100          10047602 ns/op
BenchmarkSumRaw-8                          21657             56056 ns/op
BenchmarkCustomSumCollector-8                385           3075337 ns/op
BenchmarkGroupSumRaw-8                      1020           1104203 ns/op
BenchmarkGroupSum-8                          100          12056522 ns/op
PASS
ok      github.com/a3d21/gostream       26.355s

```

**结论：**
1. 与原生操作相比，`ToSlice`、`GroupBy`、`Sum`性能相差较大
2. 因为一般业务系统数据规模小（< 1000），耗时主要在IO，所以使用gostream可感知的影响不大

