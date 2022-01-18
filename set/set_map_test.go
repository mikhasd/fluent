package set

import (
	"testing"

	"github.com/mikhasd/fluent"
	"github.com/mikhasd/fluent/iterator"
	"github.com/stretchr/testify/assert"
)

func Test_MapSet_Contains_True(t *testing.T) {
	value := 654956
	s := New[int]()
	s.Add(value)

	assert.True(t, s.Contains(value))
}

func Test_MapSet_Contains_False(t *testing.T) {
	value := 654956
	s := New[int]()
	s.Add(value)

	assert.False(t, s.Contains(value+1))
}

func Test_MapSet_ContainsAll_True(t *testing.T) {
	values := []int{654956, 987987, 64654, 8686, 45211}
	s := New[int]()
	for _, val := range values {
		s.Add(val)
	}

	assert.True(t, s.ContainsAll(iterator.ArrayIterable(values)))
}

func Test_MapSet_ContainsAll_False(t *testing.T) {
	values := []int{654956, 987987, 64654, 8686, 45211}
	s := New[int]()
	for _, val := range values {
		s.Add(val)
	}

	test := []int{1}

	assert.False(t, s.ContainsAll(iterator.ArrayIterable(test)))
}

func Test_MapSet_AddAll(t *testing.T) {
	data := []int{654956, 987987, 64654, 8686, 45211}
	values := iterator.ArrayIterable(data)
	s := New[int]()

	assert.Equal(t, 0, s.Size(), "initial size")
	assert.True(t, s.Empty(), "empty")

	s.AddAll(values)

	assert.Equal(t, len(data), s.Size(), "final size")
	assert.False(t, s.Empty(), "empty")
}

func Test_MapSet_Iterator(t *testing.T) {
	data := []int{0, 1, 2, 3, 4}
	s := FromArray(data)
	it := s.Iterator()

	var o fluent.Option[int]

	for i := 0; i < len(data); i++ {
		o = it.Next()
		assert.True(t, o.IsPresent(), "present")
		assert.Contains(t, data, o.Get(), "value")
	}

	o = it.Next()

	assert.False(t, o.IsPresent(), "last present")
}

func Test_MapSet_ForEach(t *testing.T) {
	data := []int{0, 1, 2, 3, 4}

	counter := new(int)
	*counter = 0

	FromArray(data).ForEach(func(actual int) {
		assert.Contains(t, data, actual)
		*counter = *counter + 1
	})

	assert.Equal(t, len(data), *counter, "iterations")
}

func Test_MapSet_Remove(t *testing.T) {
	value := 8454654
	data := []int{value}
	s := FromArray(data)
	size := s.Size()

	assert.True(t, s.Contains(value), "contains")

	s.Remove(value)

	assert.False(t, s.Contains(value), "contains")
	assert.Equal(t, size-1, s.Size(), "size")
}
