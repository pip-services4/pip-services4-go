package convert

import (
	"fmt"
	"strings"
	"time"
)

// BooleanConverter converts arbitrary values to boolean values using extended conversion rules:
// 	- Numbers: above 0, less much 0 are true; equal to 0 are false
// 	- Strings: "true", "yes", "T", "Y", "1" are true, "false", "no", "F", "N" are false
// 	- DateTime: above 0, less much 0 total milliseconds are true, equal to 0 are false
//
// Example:
//
//  value1, ok1 := convert.BooleanConverter.ToNullableBoolean(true)
//  value2, ok2 := convert.BooleanConverter.ToNullableBoolean("yes")
//  value3, ok3 := convert.BooleanConverter.ToNullableBoolean(1)
//  value4, ok4 := convert.BooleanConverter.ToNullableBoolean(struct{}{})
//  fmt.Println(value1, ok1) // true, true
//  fmt.Println(value2, ok2) // true, true
//  fmt.Println(value3, ok3) // true, true
//  fmt.Println(value4, ok4)  // false, false
var BooleanConverter = &_TBooleanConverter{}

type _TBooleanConverter struct{}

// ToNullableBoolean converts value into boolean or returns null when conversion is not possible.
// Parameters: "value" - the value to convert.
// Returns: boolean value and true or false and false when conversion is not supported.
func (c *_TBooleanConverter) ToNullableBoolean(value any) (bool, bool) {
	return toNullableBoolean(value)
}

// ToBoolean converts value into boolean or returns false when conversion is not possible.
// Parameters: "value" - the value to convert.
// Returns: boolean value or false when conversion is not supported.
func (c *_TBooleanConverter) ToBoolean(value any) bool {
	return toBoolean(value)
}

// ToBooleanWithDefault converts value into boolean or returns default value when conversion is not possible
// Parameters:"value" - the value to convert.
//  "defaultValue" - the default value
// Returns: boolean value or default when conversion is not supported.
func (c *_TBooleanConverter) ToBooleanWithDefault(value any, defaultValue bool) bool {
	return toBooleanWithDefault(value, defaultValue)
}

// ToNullableBoolean converts value into boolean or returns null when conversion is not possible.
// Parameters: "value" - the value to convert.
// Returns: boolean value and true or false and false when conversion is not supported.
func toNullableBoolean(value any) (bool, bool) {
	if value == nil {
		return false, false
	}

	var v string

	switch value.(type) {
	case bool:
		r, ok := value.(bool)
		return r, ok
	case string:
		if _v, ok := value.(string); ok {
			v = strings.ToLower(_v)
		}
		break
	case time.Duration:
		if d, ok := value.(time.Duration); ok {
			return d.Nanoseconds() > 0, true
		}
		break
	default:
		v = strings.ToLower(fmt.Sprint(value))
	}

	if v == "1" || v == "true" || v == "t" || v == "yes" || v == "y" {
		return true, true
	}

	if v == "0" || v == "false" || v == "f" || v == "no" || v == "n" {
		return false, true
	}

	return false, false
}

// ToBoolean converts value into boolean or returns false when conversion is not possible.
// Parameters: "value" - the value to convert.
// Returns: boolean value or false when conversion is not supported.
func toBoolean(value any) bool {
	return toBooleanWithDefault(value, false)
}

// ToBooleanWithDefault converts value into boolean or returns default value when conversion is not possible
// Parameters:
//  "value" - the value to convert.
//  "defaultValue" - the default value.
// Returns: boolean value or default when conversion is not supported.
func toBooleanWithDefault(value any, defaultValue bool) bool {
	if r, ok := toNullableBoolean(value); ok {
		return r
	}
	return defaultValue
}
