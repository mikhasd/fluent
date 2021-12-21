package fluent

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_OptionEmpty_Present(t *testing.T) {
	o := OptionEmpty[int]()
	assert.False(t, o.Present())
}

func Test_OptionEmpty_Get(t *testing.T) {
	o := OptionEmpty[int]()

	defer func() {
		err := recover()
		assert.Equal(t, "empty option", err)
	}()

	o.Get()

	t.Error("should panic")
}

func Test_OptionEmpty_Map(t *testing.T) {
	o := OptionEmpty[int]()

	called := new(bool)
	*called = false

	o = o.Map(func(a int) int {
		*called = true
		return a * 2
	})

	assert.False(t, o.Present(), "present")
	assert.False(t, *called, "mapper called")
}

func Test_OptionEmpty_OrElse(t *testing.T) {
	o := OptionEmpty[int]()
	expected := 987654321
	actual := o.OrElse(expected)

	assert.Equal(t, expected, actual)
}

func Test_OptionEmpty_OrElseGet(t *testing.T) {
	o := OptionEmpty[int]()
	expected := 987654321
	actual := o.OrElseGet(func() int {
		return expected
	})

	assert.Equal(t, expected, actual)
}

func Test_OptionEmpty_Or(t *testing.T) {
	o := OptionEmpty[int]()
	expected := 987654321
	actual := o.Or(func() Option[int] {
		return OptionPresent(expected)
	})

	assert.True(t, actual.Present(), "present")
	assert.Equal(t, expected, actual.Get())
}

func Test_OptionError_OrError(t *testing.T) {
	o := OptionEmpty[int]()

	err := errors.New("err")

	actual := o.OrError(err)

	assert.True(t, actual.IsErr(), "IsErr")
	assert.Equal(t, err, actual.GetErr(), "error")
}

func Test_OptionEmpty_String(t *testing.T) {
	o := OptionEmpty[int]()
	assert.NotEmpty(t, o.String())
}
