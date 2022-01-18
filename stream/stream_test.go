package stream

import (
	"testing"

	"github.com/mikhasd/fluent/set"
	"github.com/stretchr/testify/assert"
)

func Test_Map(t *testing.T) {
	s := Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	expected := []bool{false, true, false, true, false, true, false, true, false, true}
	mapped := Map(s, func(i int) bool {
		return i%2 == 0
	})

	actual := mapped.Array()
	assert.Equal(t, expected, actual)
}

func Test_MapArray(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	toFloat := func(val int) float32 {
		return float32(val)
	}
	even := func(val float32) bool {
		return int(val)%2 == 0
	}

	computed := set.FromIterable[float32](
		MapArray(arr, toFloat).Filter(even),
	)

	expected := set.FromArray([]float32{2, 4, 6, 8, 10})

	assert.True(t, computed.ContainsAll(expected), "content")
	assert.Equal(t, expected.Size(), computed.Size(), "size")
}
