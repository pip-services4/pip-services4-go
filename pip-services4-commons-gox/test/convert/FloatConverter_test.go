package test_convert

import (
	"testing"

	"github.com/pip-services4/pip-services4-commons-go/convert"
	"github.com/stretchr/testify/assert"
)

func TestToFloat(t *testing.T) {
	val, ok := convert.FloatConverter.ToNullableFloat(nil)
	assert.False(t, ok)
	assert.Equal(t, float32(0), val)

	assert.Equal(t, float32(123.), convert.FloatConverter.ToFloat(123))
	assert.Equal(t, float32(123.456), convert.FloatConverter.ToFloat(123.456))
	assert.Equal(t, float32(123.), convert.FloatConverter.ToFloat("123"))
	assert.Equal(t, float32(123.456), convert.FloatConverter.ToFloat("123.456"))

	assert.Equal(t, float32(123.), convert.FloatConverter.ToFloatWithDefault(nil, 123))
	assert.Equal(t, float32(0.), convert.FloatConverter.ToFloatWithDefault(false, 123))
	assert.Equal(t, float32(123.), convert.FloatConverter.ToFloatWithDefault("ABC", 123))
}
