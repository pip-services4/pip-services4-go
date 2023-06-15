package convert

import (
	"strconv"
	"time"
)

// DoubleConverter converts arbitrary values into double using extended conversion rules:
// - Strings are converted to double values
// - DateTime: total number of milliseconds since unix epoсh
// - Boolean: 1 for true and 0 for false
//
// Example:
//
//  value1, ok1 := convert.DoubleConverter.ToNullableDouble("ABC")
//  value2, ok2 := convert.DoubleConverter.ToNullableDouble("123.456")
//  value3, ok3 := convert.DoubleConverter.ToNullableDouble(true)
//  value4, ok4 := convert.DoubleConverter.ToNullableDouble(time.Now())
//  fmt.Println(value1, ok1) // 0, false
//  fmt.Println(value2, ok2) // 123.456, true
//  fmt.Println(value3, ok3) // 1, true
//  fmt.Println(value4, ok4) // current milliseconds (e.g. 1.566333114e+09), true
var DoubleConverter = &_TDoubleConverter{}

type _TDoubleConverter struct{}

// ToNullableDouble converts value into doubles or returns null when conversion is not possible.
// Parameters: "value" - the value to convert.
// Returns: double value and true or 0 and false when conversion is not supported.
func (c *_TDoubleConverter) ToNullableDouble(value any) (float64, bool) {
	return toNullableDouble(value)
}

// ToDouble converts value into doubles or returns 0 when conversion is not possible.
// Parameters: "value" - the value to convert.
// Returns: double value or 0 when conversion is not supported.
func (c *_TDoubleConverter) ToDouble(value any) float64 {
	return toDouble(value)
}

// ToDoubleWithDefault converts value into doubles or returns default when conversion is not possible.
// Parameters:
//  "value" - the value to convert.
//  "defaultValue" - the default value
// Returns: double value or default when conversion is not supported.
func (c *_TDoubleConverter) ToDoubleWithDefault(value any, defaultValue float64) float64 {
	return toDoubleWithDefault(value, defaultValue)
}

// ToNullableDouble converts value into doubles or returns null when conversion is not possible.
// Parameters: "value" - the value to convert.
// Returns: double value and true or 0 and false when conversion is not supported.
func toNullableDouble(value any) (float64, bool) {
	if value == nil {
		return 0, false
	}

	switch value.(type) {
	case int8:
		r, ok := value.(int8)
		return float64(r), ok
	case uint8:
		r, ok := value.(uint8)
		return float64(r), ok
	case int:
		r, ok := value.(int)
		return float64(r), ok
	case int16:
		r, ok := value.(int16)
		return float64(r), ok
	case uint16:
		r, ok := value.(uint16)
		return float64(r), ok
	case int32:
		r, ok := value.(int32)
		return float64(r), ok
	case uint32:
		r, ok := value.(uint32)
		return float64(r), ok
	case int64:
		r, ok := value.(int64)
		return float64(r), ok
	case uint64:
		r, ok := value.(uint64)
		return float64(r), ok
	case float32:
		r, ok := value.(float32)
		return float64(r), ok
	case float64:
		r, ok := value.(float64)
		return r, ok

	case bool:
		if r, ok := value.(bool); ok {
			if r {
				return 1.0, true
			} else {
				return 0.0, true
			}
		}
		break

	case time.Time:
		if r, ok := value.(time.Time); ok {
			return float64(r.Unix()), true
		}
		break

	case time.Duration:
		if r, ok := value.(time.Duration); ok {
			return float64(r.Nanoseconds()) / float64(1000000), true
		}
		break

	case string:
		if r, ok := value.(string); ok {
			if v, err := strconv.ParseFloat(r, 0); err == nil {
				return v, true
			}
		}
		break
	}

	return 0, false
}

// ToDouble converts value into doubles or returns 0 when conversion is not possible.
// Parameters: "value" - the value to convert.
// Returns: double value or 0 when conversion is not supported.
func toDouble(value any) float64 {
	return toDoubleWithDefault(value, 0)
}

// ToDoubleWithDefault converts value into doubles or returns default when conversion is not possible.
// Parameters:
//  "value" - the value to convert.
//  "defaultValue" - the default value.
// Returns: double value or default when conversion is not supported.
func toDoubleWithDefault(value any, defaultValue float64) float64 {
	if r, ok := toNullableDouble(value); ok {
		return r
	}
	return defaultValue
}
