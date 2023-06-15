package test_convert

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	"github.com/stretchr/testify/assert"
)

func TestToDouble(t *testing.T) {
	val, ok := convert.DoubleConverter.ToNullableDouble(nil)
	assert.False(t, ok)
	assert.Equal(t, float64(0), val)

	assert.Equal(t, 123., convert.DoubleConverter.ToDouble(123))
	assert.Equal(t, 123.456, convert.DoubleConverter.ToDouble(123.456))
	assert.Equal(t, 123., convert.DoubleConverter.ToDouble("123"))
	assert.Equal(t, 123.456, convert.DoubleConverter.ToDouble("123.456"))

	assert.Equal(t, 123., convert.DoubleConverter.ToDoubleWithDefault(nil, 123))
	assert.Equal(t, 0., convert.DoubleConverter.ToDoubleWithDefault(false, 123))
	assert.Equal(t, 123., convert.DoubleConverter.ToDoubleWithDefault("ABC", 123))
}
