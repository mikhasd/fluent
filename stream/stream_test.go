package stream

import (
	"testing"

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
