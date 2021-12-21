package fluent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testValue = 951753852

func Test_ResultOk_IsOk(t *testing.T) {
	r := Ok(testValue)

	assert.True(t, r.IsOk())
}

func Test_ResultOk_IsErr(t *testing.T) {
	r := Ok(testValue)

	assert.False(t, r.IsErr())
}

func Test_ResultOk_Ok(t *testing.T) {
	r := Ok(testValue)

	ok := r.Ok()

	assert.True(t, ok.Present(), "present")
}

func Test_ResultOk_Err(t *testing.T) {
	r := Ok(testValue)

	e := r.Err()

	assert.False(t, e.Present(), "present")
}

func Test_ResultOk_Map(t *testing.T) {
	r := Ok(testValue)
	called := new(bool)
	*called = false
	actual := r.Map(func(val int) int {
		*called = true
		return val * 2
	})

	assert.True(t, actual.IsOk(), "IsErr")
	assert.True(t, *called, "mapper called")
	assert.Equal(t, testValue*2, actual.Get())
}

func Test_ResultOk_MapErr(t *testing.T) {
	r := Ok(testValue)
	called := new(bool)
	*called = false
	expected := 987654231
	actual := r.MapErr(func(e error) int {
		*called = true
		return expected
	})

	assert.False(t, actual.IsErr(), "IsErr")
	assert.False(t, *called, "mapper called")
}

func Test_ResultOk_Get(t *testing.T) {
	r := Ok(testValue)

	actual := r.Get()

	assert.Equal(t, testValue, actual)
}

func Test_ResultOk_GetErr(t *testing.T) {
	r := Ok(testValue)
	defer func() {
		recover()
	}()
	r.GetErr()
	t.Error("should panic")
}

func Test_ResultOk_OrElse(t *testing.T) {
	r := Ok(testValue)
	actual := r.OrElse(1)
	assert.Equal(t, testValue, actual)
}

func Test_ResultOk_OrElseGet(t *testing.T) {
	r := Ok(testValue)
	actual := r.OrElseGet(func() int {
		return 1
	})
	assert.Equal(t, testValue, actual)
}

func Test_ResultOk_Or(t *testing.T) {
	r := Ok(testValue)
	actual := r.Or(func() Result[int] {
		return Ok(1)
	})
	assert.Equal(t, testValue, actual.Get())
}

func Test_ResultOk_String(t *testing.T) {
	r := Ok(testValue)
	assert.NotEmpty(t, r.String())
}
