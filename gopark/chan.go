package gopark

import (
	"time"
)

// BufferChan 将chan数据按数量和超时条件缓存，用于批量操作优化
func BufferChan[T any](in chan T, size int, timeout time.Duration) (out chan []T) {
	out = make(chan []T)

	if size <= 0 {
		panic("size should gt 0")
	}
	if timeout <= 0 {
		panic("timeout should gt 0")
	}

	go func() {
		var vs []T

		var flush = func() {
			out <- vs
			vs = nil
		}

		for {
			select {
			case v, ok := <-in:
				if ok {
					vs = append(vs, v)
					if len(vs) == size {
						flush()
					}
				} else {
					if len(vs) > 0 {
						flush()
					}
					close(out)
					return
				}
			case <-time.After(timeout):
				if len(vs) > 0 {
					flush()
				}
			}
		}
	}()

	return
}
