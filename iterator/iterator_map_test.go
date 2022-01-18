package iterator

import (
	"testing"

	"github.com/mikhasd/fluent"
	"github.com/stretchr/testify/assert"
)

var mapTestKeys = []string{
	"a",
	"b",
	"c",
	"d",
	"e",
	"f",
	"g",
	"h",
	"i",
	"j",
}

var mapTestValues = []int{
	2, 4, 6, 8, 10, 12, 14, 16, 18, 20,
}

var mapTestData map[string]int

func init() {
	mapTestData = make(map[string]int)
	for index := range mapTestKeys {
		mapTestData[mapTestKeys[index]] = mapTestValues[index]
	}
}

func Test_MapKeys(t *testing.T) {
	keys := MapKeys(mapTestData)

	count := 0
	var o fluent.Option[string]
	for o = keys.Next(); o.IsPresent(); o = keys.Next() {
		count++
		val := o.Get()
		assert.Contains(t, mapTestKeys, val, "value")
	}

	assert.Equal(t, len(mapTestKeys), count, "size")
	assert.NotNil(t, o, "option")
	assert.False(t, o.IsPresent(), "present")
}

func Test_MapValues(t *testing.T) {
	values := MapValues(mapTestData)

	count := 0

	var o fluent.Option[int]
	for o = values.Next(); o.IsPresent(); o = values.Next() {
		count++
		val := o.Get()
		assert.Contains(t, mapTestValues, val, "value")
	}

	assert.Equal(t, len(mapTestValues), count, "size")
	assert.NotNil(t, o, "option")
	assert.False(t, o.IsPresent(), "present")
}

func Test_mapIterator_Next(t *testing.T) {
	it := FromMap(mapTestData)

	count := 0
	var o fluent.Option[MapEntry[string, int]]
	for o = it.Next(); o.IsPresent(); o = it.Next() {
		count++
		kv := o.Get()
		val, found := mapTestData[kv.Key]
		assert.True(t, found, "found")
		assert.Equal(t, val, kv.Value)
	}

	assert.Equal(t, len(mapTestValues), count, "size")
	assert.NotNil(t, o, "option")
	assert.False(t, o.IsPresent(), "present")

}
