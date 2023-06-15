package convert

import (
	"time"
)

// DurationConverter Converts arbitrary values into time.Duration values.
//
// Example:
//
//  value1, ok1 := convert.DurationConverter.ToNullableDuration("123")
//  value2, ok2 := convert.DurationConverter.ToNullableDuration(123)
//  value3, ok3 := convert.DurationConverter.ToNullableDuration(123 * time.Second)
//  fmt.Println(value1, ok1) // 123ms, true
//  fmt.Println(value2, ok2) // 123ms, true
//  fmt.Println(value3, ok3) // 2m3s, true
var DurationConverter = &_TDurationConverter{}

type _TDurationConverter struct{}

// ToNullableDuration converts value into time.Duration or returns null when conversion is not possible.
// Parameters: "value" - the value to convert.
// Returns: time.Duration value and true or 0 and false when conversion is not supported.
func (c *_TDurationConverter) ToNullableDuration(value any) (time.Duration, bool) {
	return toNullableDuration(value)
}

// ToDuration converts value into time.Duration or returns current when conversion is not possible.
// Parameters: "value" - the value to convert.
// Returns: time.Duration value or current when conversion is not supported.
func (c *_TDurationConverter) ToDuration(value any) time.Duration {
	return toDuration(value)
}

// ToDurationWithDefault converts value into time.Duration or returns default when conversion is not possible.
// Parameters:
//  "value" - the value to convert.
//  "defaultValue" - the default value.
// Returns: time.Duration value or default when conversion is not supported.
func (c *_TDurationConverter) ToDurationWithDefault(value any, defaultValue time.Duration) time.Duration {
	return toDurationWithDefault(value, defaultValue)
}

// ToNullableDuration converts value into time.Duration or returns null when conversion is not possible.
// Parameters: "value" - the value to convert.
// Returns: time.Duration value and true or 0 and false when conversion is not supported.
func toNullableDuration(value any) (time.Duration, bool) {
	if value == nil {
		return 0, false
	}

	var r time.Duration

	switch value.(type) {
	case int8:
		r = (time.Duration)(value.(int8)) * time.Millisecond
		break
	case uint8:
		r = (time.Duration)(value.(uint8)) * time.Millisecond
		break
	case int:
		r = (time.Duration)(value.(int)) * time.Millisecond
		break
	case int16:
		r = (time.Duration)(value.(int16)) * time.Millisecond
		break
	case uint16:
		r = (time.Duration)(value.(uint16)) * time.Millisecond
		break
	case int32:
		r = (time.Duration)(value.(int32)) * time.Millisecond
		break
	case uint32:
		r = (time.Duration)(value.(uint32)) * time.Millisecond
		break
	case int64:
		r = (time.Duration)(value.(int64)) * time.Millisecond
		break
	case uint64:
		r = (time.Duration)(value.(uint64)) * time.Millisecond
		break
	case float32:
		r = (time.Duration)(value.(float32)) * time.Millisecond
		break
	case float64:
		r = (time.Duration)(value.(float64)) * time.Millisecond
		break

	case time.Duration:
		r = value.(time.Duration)
		break

	case string:
		v := value.(string)
		var err error
		r, err = time.ParseDuration(v)
		if err != nil {
			r = (time.Duration)(LongConverter.ToLong(value)) * time.Millisecond
		}
		break

	default:
		return 0, false
	}

	return r, true
}

// ToDuration converts value into time.Duration or returns current when conversion is not possible.
// Parameters: "value" - the value to convert.
// Returns: time.Duration value or current when conversion is not supported.
func toDuration(value any) time.Duration {
	return toDurationWithDefault(value, 0*time.Millisecond)
}

// ToDurationWithDefault converts value into time.Duration or returns default when conversion is not possible.
// Parameters:
//  "value" - the value to convert.
//  "defaultValue" - the default value.
// Returns: time.Duration value or default when conversion is not supported.
func toDurationWithDefault(value any, defaultValue time.Duration) time.Duration {
	if r, ok := toNullableDuration(value); ok {
		return r
	}
	return defaultValue
}
