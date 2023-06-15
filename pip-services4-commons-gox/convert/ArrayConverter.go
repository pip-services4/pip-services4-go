package convert

import (
	"reflect"
	"strings"
)

// ArrayConverter converts arbitrary values into array objects.
//
//	Example:
//		value1 := convert.ArrayConverter.ToArray([...]int{1, 2})
//		value2 := convert.ArrayConverter.ToArray(1)
//		value3 := convert.ArrayConverter.ListToArray("1,2,3")
//
//		fmt.Println(value1) // [1 2]
//		fmt.Println(value2) // [1]
//		fmt.Println(value3) // [1 2 3]
var ArrayConverter = &_TArrayConverter{}

type _TArrayConverter struct{}

// ToNullableArray converts value into array object. Single values are converted into arrays with a single element.
// 	Parameters: "value" - the value to convert.
// 	Returns: array object and true or null and false when value is null.
func (c *_TArrayConverter) ToNullableArray(value any) ([]any, bool) {
	return toNullableArray(value)
}

// ToArray converts value into array object with empty array as default. Single values
// are converted into arrays with single element.
// 	Parameters: "value" - the value to convert.
// 	Returns: array object or empty array when value is null.
func (c *_TArrayConverter) ToArray(value any) []any {
	return toArray(value)
}

// ToArrayWithDefault converts value into array object with empty array as default. Single values
// are converted into arrays with single element.
// 	Parameters:
//  	"value" - the value to convert.
//  	"defaultValue" - default array object.
// 	Returns: array object or empty array when value is null.
func (c *_TArrayConverter) ToArrayWithDefault(value any, defaultValue []any) []any {
	return toArrayWithDefault(value, defaultValue)
}

// ListToArray converts value into array object with empty array as default.
// Strings with comma-delimited values are split into array of strings.
// 	Parameters:
//  	"value" - the list to convert.
// 	Returns: array object or empty array when value is null
func (c *_TArrayConverter) ListToArray(value any) []any {
	return listToArray(value)
}

// ToNullableArray converts value into array object. Single values are converted into arrays with a single element.
// 	Parameters:
//  	"value" - the value to convert.
// 	Returns: array object and true or null and false when value is null.
func toNullableArray(value any) ([]any, bool) {
	if value == nil {
		return nil, false
	}

	v := reflect.ValueOf(value)

	switch v.Kind() {

	case reflect.Map:
		r := mapToArray(v)
		return r, true

	case reflect.Array, reflect.Slice:
		r := arrayToArray(v)
		return r, true

	default:
		value = valueToInterface(v)
		r := []any{value}
		return r, true
	}
}

// ToArray converts value into array object with empty array as default. Single values
// are converted into arrays with single element.
// 	Parameters:
//  	"value" - the value to convert.
// 	Returns: array object or empty array when value is null.
func toArray(value any) []any {
	return toArrayWithDefault(value, []any{})
}

// ToArrayWithDefault converts value into array object with empty array as default. Single values
// are converted into arrays with single element.
// 	Parameters:
//  	"value" - the value to convert.
//  	"defaultValue" - default array object.
// 	Returns: array object or empty array when value is null.
func toArrayWithDefault(value any, defaultValue []any) []any {
	if m, ok := toNullableArray(value); ok {
		return m
	}
	return defaultValue
}

// ListToArray converts value into array object with empty array as default.
// Strings with comma-delimited values are split into array of strings.
// 	Parameters:
//  	"value" - the list to convert.
// 	Returns: array object or empty array when value is null
func listToArray(value any) []any {
	if value == nil {
		return []any{}
	}

	v := reflect.ValueOf(value)

	if v.Kind() == reflect.String {
		value = strings.Split(value.(string), ",")
	}

	return toArray(value)
}
