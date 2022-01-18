package iterator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var singleTestData = 794621687964

func Test_singleItemIterator_Next(t *testing.T) {
	it := singleItemIterator[int]{
		consumed: false,
		item:     singleTestData,
	}

	o := it.Next()
	assert.NotNil(t, o, "option")
	assert.True(t, o.IsPresent(), "present")
	assert.Equal(t, singleTestData, o.Get())

	o = it.Next()
	assert.NotNil(t, o, "option")
	assert.False(t, o.IsPresent(), "present")
}

func Test_singleItemIterator_Size(t *testing.T) {
	var it Iterator[int] = &singleItemIterator[int]{
		consumed: false,
		item:     singleTestData,
	}

	o := Size(it)
	assert.NotNil(t, o, "option")
	assert.True(t, o.IsPresent(), "present")
	assert.Equal(t, 1, o.Get())

}
