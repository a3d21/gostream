package gostream

import (
	"reflect"
	"time"
)

// BufferChan 对channel运行缓存
// 当接收消息达到`size`或超过`timeout`未收到新消息时，发送消息
// 参数说明
//   typ  slice类型参数
//   size  缓存数量
//   timeout  超时时间
func (s Stream) BufferChan(typ interface{}, size int, timeout time.Duration) Stream {
	t := reflect.TypeOf(typ)
	if t.Kind() != reflect.Slice {
		panic("typ should be slice")
	}
	if size <= 0 {
		panic("size should gt 0")
	}
	if timeout <= 0 {
		panic("timeout should gt 0")
	}

	in := make(chan interface{})
	out := make(chan interface{})
	go s.OutChan(in)

	go func() {
		sv := reflect.MakeSlice(t, size, size)
		idx := 0

		var flush = func() {
			out <- sv.Slice(0, idx).Interface()
			sv = reflect.MakeSlice(t, size, size)
			idx = 0
		}

		for {
			select {
			case v, ok := <-in:
				if ok {
					sv.Index(idx).Set(reflect.ValueOf(v))
					idx++
					if idx == size {
						flush()
					}
				} else {
					if idx > 0 {
						flush()
					}
					close(out)
					return
				}
			case <-time.After(timeout):
				if idx > 0 {
					flush()
				}
			}
		}
	}()

	return From(out)
}
