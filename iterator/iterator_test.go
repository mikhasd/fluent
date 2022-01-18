package iterator

import (
	"fmt"
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

func Test_Of_single(t *testing.T) {
	it := Of(1)
	assert.IsType(t, &singleItemIterator[int]{}, it)
}

func Test_Of_many(t *testing.T) {
	it := Of(1, 2)
	assert.IsType(t, &arrayIterator[int]{}, it)
}

func Test_FromArray_empty(t *testing.T) {
	it := FromArray([]int{})
	assert.IsType(t, &emptyIterator[int]{}, it)
}

func Test_Parallel(t *testing.T) {
	it := FromArray([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30})
	stop := new(bool)
	*stop = false

	for !*stop {
		fmt.Println("invoking")
		go func() {
			o := it.Next()
			if o.Present() {
				fmt.Println("parallel", o.Get())
			} else {
				*stop = true
			}
		}()
	}
}
