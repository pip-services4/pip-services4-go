package test_convert

import (
	"reflect"
	"testing"
	"time"

	"github.com/pip-services4/pip-services4-commons-go/convert"
	//"github.com/pip-services4/pip-services4-commons-go/data"
	"github.com/stretchr/testify/assert"
)

func TestToTypeCode(t *testing.T) {
	assert.Equal(t, convert.String, convert.TypeConverter.ToTypeCode("123"))
	assert.Equal(t, convert.Integer, convert.TypeConverter.ToTypeCode(123))
	assert.Equal(t, convert.Long, convert.TypeConverter.ToTypeCode(int64(123)))
	assert.Equal(t, convert.Double, convert.TypeConverter.ToTypeCode(123.456))
	assert.Equal(t, convert.DateTime, convert.TypeConverter.ToTypeCode(time.Now()))
	assert.Equal(t, convert.Duration, convert.TypeConverter.ToTypeCode(time.Microsecond*10))
	assert.Equal(t, convert.Array, convert.TypeConverter.ToTypeCode([]int{}))
	assert.Equal(t, convert.Map, convert.TypeConverter.ToTypeCode(map[string]string{}))
	//assert.Equal(t, convert.Object, convert.TypeConverter.ToTypeCode(*data.NewEmptyAnyValue()))

	assert.Equal(t, convert.String, convert.TypeConverter.ToTypeCode(reflect.TypeOf("123")))
	assert.Equal(t, convert.Integer, convert.TypeConverter.ToTypeCode(reflect.TypeOf(123)))
	assert.Equal(t, convert.Long, convert.TypeConverter.ToTypeCode(reflect.TypeOf(int64(123))))
	assert.Equal(t, convert.Double, convert.TypeConverter.ToTypeCode(reflect.TypeOf(123.456)))
	assert.Equal(t, convert.DateTime, convert.TypeConverter.ToTypeCode(reflect.TypeOf((*time.Time)(nil))))
	assert.Equal(t, convert.Duration, convert.TypeConverter.ToTypeCode(reflect.TypeOf((*time.Duration)(nil))))
	assert.Equal(t, convert.Array, convert.TypeConverter.ToTypeCode(reflect.TypeOf([]int{})))
	//assert.Equal(t, convert.Map, convert.TypeConverter.ToTypeCode(reflect.TypeOf((*config.ConfigParams)(nil))))
	//assert.Equal(t, convert.Object, convert.TypeConverter.ToTypeCode(reflect.TypeOf((*data.AnyValue)(nil))))
}

func TestToNullableType(t *testing.T) {
	tp, ok := convert.TypeConverter.ToNullableType(convert.String, 123)
	assert.True(t, ok)
	assert.Equal(t, "123", tp)
	//assert.Equal(t, 123, *convert.TypeConverter.ToNullableType(convert.Integer, "123").(*int))
	//assert.Equal(t, int64(123), *convert.TypeConverter.ToNullableType(convert.Long, 123.456).(*int64))
	//assert.True(t, 123-*convert.TypeConverter.ToNullableType(convert.Float, 123).(*float32) < 0.001)
	//assert.True(t, 123-*convert.TypeConverter.ToNullableType(convert.Double, 123).(*float64) < 0.001)
	//assert.Equal(t, convert.DateTimeConverter.ToDateTime("1975-04-08T17:30:00.00Z"),
	//	*convert.TypeConverter.ToNullableType(convert.DateTime, "1975-04-08T17:30:00.00Z").(*time.Time))
	//assert.Equal(t, 1, len(*convert.TypeConverter.ToNullableType(convert.Array, 123).(*[]interface{})))
	//assert.Equal(t, 1, convert.TypeConverter .toNullableType<any>(convert.Map, StringValueMap.fromString("abc=123")).length)
}

func TestToType(t *testing.T) {
	assert.Equal(t, "123", convert.TypeConverter.ToType(convert.String, 123))
	assert.Equal(t, 123, convert.TypeConverter.ToType(convert.Integer, "123"))
	assert.Equal(t, int64(123), convert.TypeConverter.ToType(convert.Long, 123.456))
	assert.True(t, 123-convert.TypeConverter.ToType(convert.Float, 123).(float32) < 0.001)
	assert.True(t, 123-convert.TypeConverter.ToType(convert.Double, 123).(float64) < 0.001)
	assert.Equal(t, convert.DateTimeConverter.ToDateTime("1975-04-08T17:30:00.00Z"),
		convert.TypeConverter.ToType(convert.DateTime, "1975-04-08T17:30:00.00Z"))
	assert.Equal(t, 1, len(convert.TypeConverter.ToType(convert.Array, 123).([]interface{})))
	//assert.Equal(t, 1, convert.TypeConverter.ToType<any>(convert.Map, StringValueMap.fromString("abc=123")).length)
}

func TestToTypeWithDefault(t *testing.T) {
	assert.Equal(t, "123", convert.TypeConverter.ToTypeWithDefault(convert.String, nil, "123"))
	assert.Equal(t, 123, convert.TypeConverter.ToTypeWithDefault(convert.Integer, nil, 123))
	assert.Equal(t, 123, convert.TypeConverter.ToTypeWithDefault(convert.Long, nil, 123))
	assert.True(t, 123-convert.TypeConverter.ToTypeWithDefault(convert.Float, nil, float32(123)).(float32) < 0.001)
	assert.True(t, 123-convert.TypeConverter.ToTypeWithDefault(convert.Double, nil, float64(123.)).(float64) < 0.001)
	assert.Equal(t, convert.DateTimeConverter.ToDateTime("1975-04-08T17:30:00.00Z"),
		convert.TypeConverter.ToTypeWithDefault(convert.DateTime, "1975-04-08T17:30:00.00Z", time.Time{}))
	assert.Equal(t, 1, len(convert.TypeConverter.ToTypeWithDefault(convert.Array, 123, []interface{}{}).([]interface{})))
	//assert.Equal(t, 1, convert.TypeConverter.ToTypeWithDefault<any>(convert.Map, StringValueMap.fromString("abc=123"), null)).length)
}
