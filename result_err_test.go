package fluent

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testErr = errors.New("err")

func Test_ResultErr_IsOk(t *testing.T) {
	r := ResultErr[int](testErr)

	assert.False(t, r.IsOk())
}

func Test_ResultErr_IsErr(t *testing.T) {
	r := ResultErr[int](testErr)

	assert.True(t, r.IsErr())
}

func Test_ResultErr_Ok(t *testing.T) {
	r := ResultErr[int](testErr)

	ok := r.Ok()

	assert.False(t, ok.Present(), "present")
}

func Test_ResultErr_Err(t *testing.T) {
	r := ResultErr[int](testErr)

	e := r.Err()

	assert.True(t, e.Present(), "present")
}

func Test_ResultErr_Map(t *testing.T) {
	r := ResultErr[int](testErr)
	called := new(bool)
	*called = false
	actual := r.Map(func(val int) int {
		*called = true
		return val * 2
	})

	assert.True(t, actual.IsErr(), "IsErr")
	assert.False(t, *called, "mapper called")
}

func Test_ResultErr_MapErr(t *testing.T) {
	r := ResultErr[int](testErr)
	called := new(bool)
	*called = false
	expected := 987654231
	actual := r.MapErr(func(e error) int {
		*called = true
		return expected
	})

	assert.False(t, actual.IsErr(), "IsErr")
	assert.True(t, *called, "mapper called")
	assert.Equal(t, expected, actual.Get())
}

func Test_ResultErr_Get(t *testing.T) {
	r := ResultErr[int](testErr)

	defer func() {
		recover()
	}()
	r.Get()
	t.Error("should panic")
}

func Test_ResultErr_GetErr(t *testing.T) {
	r := ResultErr[int](testErr)
	actual := r.GetErr()
	assert.Equal(t, testErr, actual)
}

func Test_ResultErr_OrElse(t *testing.T) {
	r := ResultErr[int](testErr)
	expected := 987654321
	actual := r.OrElse(expected)
	assert.Equal(t, expected, actual)
}

func Test_ResultErr_Or(t *testing.T) {
	r := ResultErr[int](testErr)
	expected := 987654321
	actual := r.Or(func() Result[int] {
		return ResultOk(expected)
	})
	assert.Equal(t, expected, actual.Get())
}

func Test_ResultErr_String(t *testing.T) {
	r := ResultErr[int](testErr)
	assert.NotEmpty(t, r.String())
}
