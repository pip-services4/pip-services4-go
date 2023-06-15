package test_convert

import (
	"testing"
	"time"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	"github.com/stretchr/testify/assert"
)

func TestToDateTime(t *testing.T) {
	val, ok := convert.DateTimeConverter.ToNullableDateTime(nil)
	assert.False(t, ok)
	assert.True(t, val.IsZero())

	date1 := time.Date(1975, time.April, 8, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, date1, convert.DateTimeConverter.ToDateTimeWithDefault(nil, date1))
	assert.Equal(t, date1, convert.DateTimeConverter.ToDateTime(date1))
	assert.Equal(t, date1, convert.DateTimeConverter.ToDateTime("1975-04-08T00:00:00Z"))
	assert.Equal(t, date1, convert.DateTimeConverter.ToDateTime("1975-04-08T00:00:00.00Z"))

	date2 := time.Unix(123, 0)
	assert.Equal(t, date2, convert.DateTimeConverter.ToDateTime(123))
	assert.Equal(t, date2, convert.DateTimeConverter.ToDateTime(123.456))
}
