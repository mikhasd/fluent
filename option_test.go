package fluent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_OptionFromReference_Empty(t *testing.T) {
	var val *string
	option := OfNillable(val)
	assert.False(t, option.Present())
}

func Test_OptionFromReference_Present(t *testing.T) {
	val := "test"
	option := OfNillable(&val)
	assert.True(t, option.Present())
}

func Test_MapOption_Present(t *testing.T) {
	o := Present("a")
	var called *bool
	called = new(bool)
	*called = false
	actual := MapOption(o, func(str string) int {
		*called = true
		return len(str)
	})

	assert.True(t, actual.Present(), "present")
	assert.True(t, *called, "mapper called")
	assert.Equal(t, 1, actual.Get())
}

func Test_MapOption_Empty(t *testing.T) {
	o := Empty[string]()
	var called *bool
	called = new(bool)
	*called = false
	actual := MapOption(o, func(str string) int {
		*called = true
		return len(str)
	})

	assert.False(t, actual.Present(), "present")
	assert.False(t, *called, "mapper called")
}
