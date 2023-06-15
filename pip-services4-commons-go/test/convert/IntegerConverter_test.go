package test_convert

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	"github.com/stretchr/testify/assert"
)

func TestToInteger(t *testing.T) {
	val, ok := convert.IntegerConverter.ToNullableInteger(nil)
	assert.False(t, ok)
	assert.Equal(t, 0, val)

	assert.Equal(t, int(123), convert.IntegerConverter.ToInteger(123))
	assert.Equal(t, int(123), convert.IntegerConverter.ToInteger(123.456))
	assert.Equal(t, int(123), convert.IntegerConverter.ToInteger("123"))
	assert.Equal(t, int(123), convert.IntegerConverter.ToInteger("123.456"))

	assert.Equal(t, int(123), convert.IntegerConverter.ToIntegerWithDefault(nil, 123))
	assert.Equal(t, int(0), convert.IntegerConverter.ToIntegerWithDefault(false, 123))
	assert.Equal(t, int(123), convert.IntegerConverter.ToIntegerWithDefault("ABC", 123))
}
