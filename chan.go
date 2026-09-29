package gostream

import (
	"time"
)

// BufferChan buffers channel items based on batch size and inactivity timeout.
// It emits a batch when the buffer reaches `size` or when no new message is received within `timeout`.
// Parameters:
//
//	size     maximum batch size
//	timeout  inactivity timeout duration
func BufferChan[T any](s Stream[T], size int, timeout time.Duration) Stream[[]T] {
	if size <= 0 {
		panic("size should gt 0")
	}
	if timeout <= 0 {
		panic("timeout should gt 0")
	}

	in := make(chan T)
	out := make(chan []T)
	go s.OutChan(in)

	go func() {
		buf := make([]T, 0, size)

		flush := func() {
			out <- buf
			buf = make([]T, 0, size)
		}

		for {
			select {
			case v, ok := <-in:
				if ok {
					buf = append(buf, v)
					if len(buf) == size {
						flush()
					}
				} else {
					if len(buf) > 0 {
						flush()
					}
					close(out)
					return
				}
			case <-time.After(timeout):
				if len(buf) > 0 {
					flush()
				}
			}
		}
	}()

	return FromChannel(out)
}

// BufferChanInterval buffers channel items based on batch size and fixed interval.
// It emits a batch when the buffer reaches `size` or when `interval` elapses.
// Parameters:
//
//	size      maximum batch size
//	interval  time window duration
func BufferChanInterval[T any](s Stream[T], size int, interval time.Duration) Stream[[]T] {
	if size <= 0 {
		panic("size should gt 0")
	}
	if interval <= 0 {
		panic("interval should gt 0")
	}

	in := make(chan T)
	out := make(chan []T)
	go s.OutChan(in)

	go func() {
		buf := make([]T, 0, size)

		after := time.After(time.Hour)
		resetAfter := func() {
			after = time.After(interval)
		}

		flush := func() {
			out <- buf
			buf = make([]T, 0, size)
			after = time.After(time.Hour)
		}

		for {
			select {
			case v, ok := <-in:
				if ok {
					buf = append(buf, v)
					if len(buf) == 1 {
						resetAfter()
					}
					if len(buf) == size {
						flush()
					}
				} else {
					if len(buf) > 0 {
						flush()
					}
					close(out)
					return
				}
			case <-after:
				if len(buf) > 0 {
					flush()
				}
			}
		}
	}()

	return FromChannel(out)
}
