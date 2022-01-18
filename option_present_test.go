package fluent

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

const value = 987654231

func Test_OptionPresent_Present(t *testing.T) {
	o := Present(value)
	assert.True(t, o.IsPresent())
}

func Test_OptionPresent_Get(t *testing.T) {
	o := Present(value)

	actual := o.Get()

	assert.Equal(t, value, actual)
}

func Test_OptionPresent_Map(t *testing.T) {
	o := Present(value)

	called := new(bool)
	*called = false

	o = o.Map(func(a int) int {
		*called = true
		return a * 2
	})

	actual := o.Get()

	assert.True(t, o.IsPresent(), "present")
	assert.True(t, *called, "mapper called")
	assert.Equal(t, value*2, actual)
}

func Test_OptionPresent_OrElse(t *testing.T) {
	o := Present(value)
	actual := o.OrElse(987654321)

	assert.Equal(t, value, actual)
}

func Test_OptionPresent_OrElseGet(t *testing.T) {
	o := Present(value)
	called := new(bool)
	*called = false
	actual := o.OrElseGet(func() int {
		*called = true
		return 98754654
	})

	assert.Equal(t, value, actual)
	assert.False(t, *called, "called")
}

func Test_OptionPresent_Or(t *testing.T) {
	o := Present(value)
	expected := 987654321
	actual := o.Or(func() Option[int] {
		return Present(expected)
	})

	assert.True(t, actual.IsPresent(), "present")
	assert.Equal(t, value, actual.Get())
}

func Test_OptionPresent_OrError(t *testing.T) {
	o := Present(value)

	err := errors.New("err")

	actual := o.OrError(err)

	assert.True(t, actual.IsOk(), "IsOk")
	assert.Equal(t, value, actual.Get(), "value")
}

func Test_OptionPresent_IfPresent(t *testing.T) {
	o := Present(value)
	called := new(bool)
	*called = false
	o.IfPresent(func(i int) {
		*called = true
		assert.Equal(t, value, i, "value")
	})
	assert.True(t, *called, "called")
}

func Test_OptionPresent_Filter_Match(t *testing.T) {
	o := Present(value)
	called := new(bool)
	*called = false
	o = o.Filter(func(i int) bool {
		*called = true
		assert.Equal(t, value, i, "value")
		return true
	})
	assert.True(t, o.IsPresent(), "present")
	assert.True(t, *called, "called")
}

func Test_OptionPresent_Filter_Mismatch(t *testing.T) {
	o := Present(value)
	called := new(bool)
	*called = false
	o = o.Filter(func(i int) bool {
		*called = true
		assert.Equal(t, value, i, "value")
		return false
	})
	assert.False(t, o.IsPresent(), "present")
	assert.True(t, *called, "called")
}

func Test_OptionPresent_String(t *testing.T) {
	o := Present(value)
	assert.NotEmpty(t, o.String())
}
