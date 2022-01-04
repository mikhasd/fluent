package array

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Map(t *testing.T) {
	in := []int{1, 2, 4, 8}
	expected := []float32{1.1, 2.2, 4.4, 8.8}
	mapper := func(i int) float32 {
		return float32(i) * 1.1
	}

	actual := Map(in, mapper)

	assert.Equal(t, expected, actual)
}

func Test_Filter(t *testing.T) {
	in := []int{1, 2, 3, 4, 5, 6}
	expected := []int{2, 4, 6}
	condition := func(i int) bool {
		return i%2 == 0
	}

	actual := Filter(in, condition)

	assert.Equal(t, expected, actual)
}

func Test_Filter_empty(t *testing.T) {
	in := []int{}
	expected := []int{}
	condition := func(i int) bool {
		return i%2 == 0
	}

	actual := Filter(in, condition)

	assert.Equal(t, expected, actual)
	assert.Len(t, actual, 0)
}
