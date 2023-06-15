package convert

import (
	"fmt"
	refl "reflect"
	"strconv"
	"strings"
	"time"
)

// StringConverter converts arbitrary values into strings using extended conversion rules:
// - Numbers: are converted with '.' as decimal point
// - DateTime: using ISO format
// - Boolean: "true" for true and "false" for false
// - Arrays: as comma-separated list
// - Other objects: using toString() method
//
// Example:
//
//  value1, ok1 = convert.StringConverter.ToString(123.456)
//  value2, ok2 = convert.StringConverter.ToString(true)
//  value3, ok3 = convert.StringConverter.ToString(time.Now())
//  value4, ok4 = convert.StringConverter.ToString([...]int{1, 2, 3})
//  fmt.Println(value1, ok1) // 123.456, true
//  fmt.Println(value2, ok2) // true, true
//  fmt.Println(value3, ok3) // 2019-08-20T23:54:47+03:00, true
//  fmt.Println(value4, ok4) // 1,2,3, true
var StringConverter = &_TStringConverter{}

type _TStringConverter struct{}

// ToNullableString converts value into string or returns null when value is null.
// Parameters: "value" - the value to convert
// Returns: string value and true or "" and false when value is null.
func (c *_TStringConverter) ToNullableString(value any) (string, bool) {
	return toNullableString(value)
}

// ToString converts value into string or returns "" when value is null.
// Parameters: "value" - the value to convert
// Returns: string value or "" when value is null.
func (c *_TStringConverter) ToString(value any) string {
	return toString(value)
}

// ToStringWithDefault converts value into string or returns default when value is null.
// Parameters:
//  "value" - the value to convert.
//  "defaultValue" - the default value.
// Returns: string value or default when value is null.
func (c *_TStringConverter) ToStringWithDefault(value any, defaultValue string) string {
	return toStringWithDefault(value, defaultValue)
}

// ToNullableString converts value into string or returns null when value is null.
// Parameters: "value" - the value to convert
// Returns: string value and true or "" and false when value is null.
func toNullableString(value any) (string, bool) {
	if value == nil {
		return "", false
	}

	switch value.(type) {
	case string:
		r, ok := value.(string)
		return r, ok

	case uint, uint32, uint64:
		r := strconv.FormatUint(LongConverter.ToULong(value), 10)
		return r, true

	case byte, int, int32, int64:
		r := strconv.FormatInt(LongConverter.ToLong(value), 10)
		return r, true

	case float32, float64:
		r := strconv.FormatFloat(DoubleConverter.ToDouble(value), 'f', -1, 64)
		return r, true

	case bool:
		if b, ok := value.(bool); ok {
			if b {
				return "true", true
			}
			return "false", true
		}
		break

	case time.Time:
		if r, ok := value.(time.Time); ok {
			return r.Format(time.RFC3339), true
		}
		break

	case time.Duration:
		if r, ok := value.(time.Duration); ok {
			return strconv.FormatInt(r.Nanoseconds()/1000000, 10), true
		}
		break

	default:
		val := refl.ValueOf(value)
		if val.Kind() == refl.Slice || val.Kind() == refl.Array {
			builder := strings.Builder{}
			for index := 0; index < val.Len(); index++ {
				if builder.Len() > 0 {
					builder.WriteString(",")
				}
				builder.WriteString(fmt.Sprint(val.Index(index).Interface()))
			}
			return builder.String(), true
		}

		r := fmt.Sprint(value)
		return r, true
	}

	return "", false
}

// ToString converts value into string or returns "" when value is null.
// Parameters: "value" - the value to convert
// Returns: string value or "" when value is null.
func toString(value any) string {
	return toStringWithDefault(value, "")
}

// ToStringWithDefault converts value into string or returns default when value is null.
// Parameters:
//  "value" - the value to convert.
//  "defaultValue" - the default value.
// Returns: string value or default when value is null.
func toStringWithDefault(value any, defaultValue string) string {
	if r, ok := toNullableString(value); ok {
		return r
	}
	return defaultValue
}
