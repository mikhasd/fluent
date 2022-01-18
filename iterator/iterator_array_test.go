package iterator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var arrayTestData = []int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}

func Test_arrayIterator_Next(t *testing.T) {
	it := FromArray(arrayTestData)

	for _, val := range arrayTestData {
		o := it.Next()
		assert.NotNil(t, o, "option")
		assert.True(t, o.IsPresent(), "present")
		assert.Equal(t, val, o.Get())
	}

	o := it.Next()
	assert.NotNil(t, o, "option")
	assert.False(t, o.IsPresent(), "present")
}

func Test_arrayIterator_Size(t *testing.T) {
	it := FromArray(arrayTestData)

	o := Size(it)
	assert.NotNil(t, o, "option")
	assert.True(t, o.IsPresent(), "present")
	assert.Equal(t, len(arrayTestData), o.Get())

}

func Test_arrayIterable_Iterator(t *testing.T) {
	iter := ArrayIterable(arrayTestData)

	it := iter.Iterator()

	for _, val := range arrayTestData {
		o := it.Next()
		assert.NotNil(t, o, "option")
		assert.True(t, o.IsPresent(), "present")
		assert.Equal(t, val, o.Get())
	}

	o := it.Next()
	assert.NotNil(t, o, "option")
	assert.False(t, o.IsPresent(), "present")
}
