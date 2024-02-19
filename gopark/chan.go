package gopark

import (
	"time"
)

// BufferChan 将chan数据按数量和超时条件缓存，用于批量操作优化。 session-window
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

// BufferChanInterval 将chan数据按数量和时间间隔缓存，用于批量操作优化。sliding-window
func BufferChanInterval[T any](in chan T, size int, interval time.Duration) (out chan []T) {
	out = make(chan []T)

	if size <= 0 {
		panic("size should gt 0")
	}
	if interval <= 0 {
		panic("interval should gt 0")
	}

	go func() {
		var vs []T

		// default long-interval
		var after = time.After(time.Hour)

		var resetAfter = func() {
			after = time.After(interval)
		}
		var flush = func() {
			out <- vs
			vs = nil
			after = time.After(time.Hour)
		}

		for {
			select {
			case v, ok := <-in:
				if ok {
					vs = append(vs, v)
					if len(vs) == 1 {
						resetAfter()
					}
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
			case <-after:
				if len(vs) > 0 {
					flush()
				}
			}
		}
	}()

	return
}
