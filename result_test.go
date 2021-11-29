package gostream

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultAllShouldBeTrue(t *testing.T) {
	all := From([]int(nil)).All(func(v interface{}) bool {
		return false
	})
	assert.True(t, all)
}

func TestDefaultAnyShouldBeFalse(t *testing.T) {
	any := From([]int(nil)).Any()
	assert.False(t, any)
}

func TestDefaultAnyWithShouldBeFalse(t *testing.T) {
	anyWith := From([]int(nil)).AnyWith(func(v interface{}) bool {
		return true
	})
	assert.False(t, anyWith)
}

func TestNilStream_First(t *testing.T) {
	v, ok := From([]int(nil)).First()
	assert.False(t, ok)
	assert.Nil(t, v)
}

func TestStream_First(t *testing.T) {
	input := []int{1, 2, 3}
	v, ok := From(input).First()
	assert.True(t, ok)
	assert.Equal(t, input[0], v)
}

func TestNilStream_Last(t *testing.T) {
	v, ok := From([]int(nil)).Last()
	assert.False(t, ok)
	assert.Nil(t, v)
}

func TestStream_Last(t *testing.T) {
	input := []int{1, 2, 3}
	v, ok := From(input).Last()
	assert.True(t, ok)
	assert.Equal(t, input[2], v)
}

func TestProcessSucc(t *testing.T) {
	input := []int{1, 2, 3}
	err := From(input).Process(func(v interface{}) error {
		return nil
	})

	assert.Nil(t, err)
}

func TestProcessFail(t *testing.T) {
	input := []int{1, 2, 3}
	err := From(input).Process(func(v interface{}) error {
		if v.(int) > 2 {
			return errors.New("some err")
		}
		return nil
	})
	assert.NotNil(t, err)
}
