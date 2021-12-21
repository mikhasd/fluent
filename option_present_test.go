package fluent

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

const value = 987654231

func Test_OptionPresent_Present(t *testing.T) {
	o := Present(value)
	assert.True(t, o.Present())
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

	assert.True(t, o.Present(), "present")
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

	assert.True(t, actual.Present(), "present")
	assert.Equal(t, value, actual.Get())
}

func Test_OptionPresent_OrError(t *testing.T) {
	o := Present(value)

	err := errors.New("err")

	actual := o.OrError(err)

	assert.True(t, actual.IsOk(), "IsOk")
	assert.Equal(t, value, actual.Get(), "value")
}

func Test_OptionPresent_String(t *testing.T) {
	o := Present(value)
	assert.NotEmpty(t, o.String())
}
