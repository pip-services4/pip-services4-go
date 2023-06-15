package test_convert

import (
	"testing"

	"github.com/pip-services4/pip-services4-commons-go/convert"
	"github.com/stretchr/testify/assert"
)

func TestToBoolean(t *testing.T) {
	val, ok := convert.BooleanConverter.ToNullableBoolean(nil)
	assert.False(t, ok)
	assert.False(t, val)

	assert.True(t, convert.BooleanConverter.ToBoolean(true))
	assert.True(t, convert.BooleanConverter.ToBoolean(1))
	assert.True(t, convert.BooleanConverter.ToBoolean("True"))
	assert.True(t, convert.BooleanConverter.ToBoolean("yes"))
	assert.True(t, convert.BooleanConverter.ToBoolean("1"))
	assert.True(t, convert.BooleanConverter.ToBoolean("Y"))

	assert.False(t, convert.BooleanConverter.ToBoolean(false))
	assert.False(t, convert.BooleanConverter.ToBoolean(0))
	assert.False(t, convert.BooleanConverter.ToBoolean("False"))
	assert.False(t, convert.BooleanConverter.ToBoolean("no"))
	assert.False(t, convert.BooleanConverter.ToBoolean("0"))
	assert.False(t, convert.BooleanConverter.ToBoolean("N"))

	assert.False(t, convert.BooleanConverter.ToBoolean(123))
	assert.False(t, convert.BooleanConverter.ToBoolean(nil))
	assert.True(t, convert.BooleanConverter.ToBooleanWithDefault("XYZ", true))
}
