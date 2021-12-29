package set

import (
	"testing"

	"github.com/mikhasd/fluent/iterator"
	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	s := New[int]()
	assert.NotNil(t, s)
	assert.Zero(t, s.Size())
}

func Test_WithSize(t *testing.T) {
	s := WithSize[int](128)
	assert.NotNil(t, s)
	assert.Zero(t, s.Size())
}

func Test_FromArray(t *testing.T) {
	data := []int{951, 753, 654, 185}
	s := FromArray(data)
	for _, val := range data {
		assert.True(t, s.Contains(val))
	}
}

func Test_FromIterable(t *testing.T) {
	data := []int{951, 753, 654, 185}
	it := iterator.ArrayIterable(data)
	s := FromIterable(it)
	for _, val := range data {
		assert.True(t, s.Contains(val))
	}
}
