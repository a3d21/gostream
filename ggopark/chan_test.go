package ggopark

import (
	. "github.com/a3d21/gostream/core"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBufferChanBySize(t *testing.T) {
	in := make(chan int)
	go func() {
		for i := 0; i < 5; i++ {
			in <- i
		}
		close(in)
	}()

	want := [][]int{{0, 1, 2}, {3, 4}}
	out := BufferChan(in, 3, time.Second)
	got := From(out).Collect(ToSlice([][]int{}))
	assert.Equal(t, want, got)
}

func TestBufferChanByTimeout(t *testing.T) {
	in := make(chan int)
	go func() {
		in <- 0
		time.Sleep(time.Millisecond * 500)
		in <- 1
		in <- 2
		time.Sleep(time.Millisecond * 500)
		in <- 3
		in <- 4
		close(in)
	}()
	want := [][]int{{0}, {1, 2}, {3, 4}}
	out := BufferChan(in, 100, time.Millisecond*300)
	got := From(out).Collect(ToSlice([][]int{}))
	assert.Equal(t, want, got)
}
