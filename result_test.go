package fluent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ResultMap_Ok(t *testing.T) {
	value := "a"
	r := Ok(value)
	expected := len(value)
	called := new(bool)
	*called = false
	actual := MapResult(r, func(val string) int {
		*called = true
		return len(val)
	})
	assert.True(t, *called, "mapper called")
	assert.Equal(t, expected, actual.Get())
}

func Test_ResultMap_Err(t *testing.T) {
	r := Err[string](testErr)
	called := new(bool)
	*called = false
	actual := MapResult(r, func(val string) int {
		*called = true
		return len(val)
	})
	assert.False(t, *called, "mapper called")
	assert.True(t, actual.IsErr())
}

func Test_ResultFromCall_Ok(t *testing.T) {
	expected := 951357897
	fn := func() (int, error) {
		return expected, nil
	}
	r := CallResult(fn)
	assert.True(t, r.IsOk(), "IsOk")
	assert.Equal(t, expected, r.Get())
}

func Test_ResultFromCall_Err(t *testing.T) {
	expected := testErr
	fn := func() (int, error) {
		return 0, expected
	}
	r := CallResult(fn)
	assert.True(t, r.IsErr(), "IsErr")
	assert.Equal(t, expected, r.GetErr())
}
