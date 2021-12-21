package fluent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_OptionFromReference_Empty(t *testing.T) {
	var val *string
	option := OptionOfNillable(val)
	assert.False(t, option.Present())
}

func Test_OptionFromReference_Present(t *testing.T) {
	val := "test"
	option := OptionOfNillable(&val)
	assert.True(t, option.Present())
}

func Test_MapOption_Present(t *testing.T) {
	o := OptionPresent("a")
	var called *bool
	called = new(bool)
	*called = false
	actual := OptionMap(o, func(str string) int {
		*called = true
		return len(str)
	})

	assert.True(t, actual.Present(), "present")
	assert.True(t, *called, "mapper called")
	assert.Equal(t, 1, actual.Get())
}

func Test_MapOption_Empty(t *testing.T) {
	o := OptionEmpty[string]()
	var called *bool
	called = new(bool)
	*called = false
	actual := OptionMap(o, func(str string) int {
		*called = true
		return len(str)
	})

	assert.False(t, actual.Present(), "present")
	assert.False(t, *called, "mapper called")
}
