package iterator

import (
	"testing"

	"github.com/mikhasd/fluent"
	"github.com/stretchr/testify/assert"
)

func Test_FuncIterator_Next(t *testing.T) {
	called := false
	value := 79844654
	it := Func(func() fluent.Option[int] {
		if !called {
			called = true
			return fluent.Present(value)
		} else {
			return fluent.Empty[int]()
		}
	})

	o := it.Next()
	assert.NotNil(t, o, "option")
	assert.True(t, o.Present(), "present")
	assert.Equal(t, value, o.Get(), "value")

	o = it.Next()

	assert.NotNil(t, o, "option")
	assert.False(t, o.Present(), "present")
}

func Test_FuncIterator_Size(t *testing.T) {
	it := Func(func() fluent.Option[int] {
		return fluent.Empty[int]()
	})

	size := Size(it)
	assert.NotNil(t, size, "option")
	assert.False(t, size.Present(), "present")
}
