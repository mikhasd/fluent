package iterator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_emptyIterator_Next(t *testing.T) {
	it := emptyIterator[int]{}
	o := it.Next()
	assert.NotNil(t, o, "option")
	assert.False(t, o.Present(), "present")
}

func Test_emptyIterator_Size(t *testing.T) {
	var it Iterator[int] = &emptyIterator[int]{}

	o := Size(it)
	assert.NotNil(t, o, "option")
	assert.True(t, o.Present(), "present")
	assert.Equal(t, 0, o.Get())

}
