package gopark

import (
	"time"
)

// BufferChan buffers channel items based on batch size and inactivity timeout for batch operations (session window).
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

// BufferChanInterval buffers channel items based on batch size and fixed interval for batch operations (sliding/tumbling window).
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
