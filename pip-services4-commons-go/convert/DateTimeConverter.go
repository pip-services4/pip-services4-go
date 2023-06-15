package convert

import (
	"time"
)

// DateTimeConverter converts arbitrary values into Date values using extended conversion rules:
// - Strings: converted using ISO time format
// - Numbers: converted using milliseconds since unix epoch
//
// Example:
//
//  value1, ok1 := convert.DateTimeConverter.ToNullableDateTime("ABC")
//  value2, ok2 := convert.DateTimeConverter.ToNullableDateTime("2019-01-01T11:30:00.0Z")
//  value3, ok3 := convert.DateTimeConverter.ToNullableDateTime(123)
//  fmt.Println(value1, ok1) // 0001-01-01 00:00:00 +0000 UTC, false
//  fmt.Println(value2, ok2) // 2019-01-01 11:30:00 +0000 UTC, true
//  fmt.Println(value3, ok3) // 1970-01-01 02:02:03 +0200 EET, true
var DateTimeConverter = &_TDateTimeConverter{}

type _TDateTimeConverter struct{}

// ToNullableDateTime converts value into Date or returns null when conversion is not possible.
// Parameters: "value" - the value to convert.
// Returns: Date value and true or zero time and false when conversion is not supported.
func (c *_TDateTimeConverter) ToNullableDateTime(value any) (time.Time, bool) {
	return toNullableDateTime(value)
}

// ToDateTime converts value into Date or returns current when conversion is not possible.
// Parameters: "value" - the value to convert.
// Returns: Date value or current when conversion is not supported.
func (c *_TDateTimeConverter) ToDateTime(value any) time.Time {
	return toDateTime(value)
}

// ToDateTimeWithDefault converts value into Date or returns default when conversion is not possible.
// Parameters:
// "value" - the value to convert.
// "defaultValue" - the default value.
// Returns: Date value or default when conversion is not supported.
func (c *_TDateTimeConverter) ToDateTimeWithDefault(value any, defaultValue time.Time) time.Time {
	return toDateTimeWithDefault(value, defaultValue)
}

// ToNullableDateTime converts value into Date or returns null when conversion is not possible.
// Parameters: "value" - the value to convert.
// Returns: Date value and true or zero time and false when conversion is not supported.
func toNullableDateTime(value any) (time.Time, bool) {
	if value == nil {
		return time.Time{}, false
	}

	var r time.Time

	switch value.(type) {
	case int8:
		r = time.Unix((int64)(value.(int8)), 0)
		break
	case uint8:
		r = time.Unix((int64)(value.(uint8)), 0)
		break
	case int:
		r = time.Unix((int64)(value.(int)), 0)
		break
	case int16:
		r = time.Unix((int64)(value.(int16)), 0)
		break
	case uint16:
		r = time.Unix((int64)(value.(uint16)), 0)
		break
	case int32:
		r = time.Unix((int64)(value.(int32)), 0)
		break
	case uint32:
		r = time.Unix((int64)(value.(uint32)), 0)
		break
	case int64:
		r = time.Unix(value.(int64), 0)
		break
	case uint64:
		r = time.Unix((int64)(value.(uint64)), 0)
		break
	case float32:
		r = time.Unix((int64)(value.(float32)), 0)
		break
	case float64:
		r = time.Unix((int64)(value.(float64)), 0)
		break

	case time.Time:
		r = value.(time.Time)
		break

	case string:
		v := value.(string)
		var err error
		r, err = time.Parse(time.RFC3339, v)
		if err != nil {
			r, err = time.Parse(time.RFC3339Nano, v)
			break
		}
		if err != nil {
			return time.Time{}, false
		}

	default:
		return time.Time{}, false
	}

	return r, true
}

// ToDateTime converts value into Date or returns current when conversion is not possible.
// Parameters: "value" - the value to convert.
// Returns: Date value or current when conversion is not supported.
func toDateTime(value any) time.Time {
	return toDateTimeWithDefault(value, time.Time{})
}

// ToDateTimeWithDefault converts value into Date or returns default when conversion is not possible.
// Parameters:
//  "value" - the value to convert.
//  "defaultValue" - the default value.
// Returns: Date value or default when conversion is not supported.
func toDateTimeWithDefault(value any, defaultValue time.Time) time.Time {
	if r, ok := toNullableDateTime(value); ok {
		return r
	}
	return defaultValue
}
