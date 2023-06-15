package convert

import "reflect"

// MapConverter converts arbitrary values into map objects using extended conversion rules:
// - Objects: property names as keys, property values as values
// - Arrays: element indexes as keys, elements as values
//
// Example:
//
//  value1, ok1 := convert.MapConverter.ToNullableMap("ABC")
//  value2, ok2 := convert.MapConverter.ToNullableMap(map[string]int{"key": 123})
//  value3, ok3 := convert.MapConverter.ToNullableMap([...]int{1, 2, 3})
//  fmt.Println(value1, ok1) // <nil>, false
//  fmt.Println(value2, ok2) // map[key:123], true
//  fmt.Println(value3, ok3) // map[0:1 1:2 2:3], true
var MapConverter = &_TMapConverter{}

type _TMapConverter struct{}

// ToNullableMap converts value into map object or returns null when conversion is not possible.
// Parameters: "value" - the value to convert
// Returns: map object and true or null and false when conversion is not supported.
func (c *_TMapConverter) ToNullableMap(value any) (map[string]any, bool) {
	return toNullableMap(value)
}

// ToMap converts value into map object or returns empty map when conversion is not possible.
// Parameters: "value" - the value to convert
// Returns: map object or empty map when conversion is not supported.
func (c *_TMapConverter) ToMap(value any) map[string]any {
	return toMap(value)
}

// ToMapWithDefault converts value into map object or returns default map when conversion is not possible.
// Parameters:
//  "value" - the value to convert
//  "defaultValue" - the default value.
// Returns: map object or default map when conversion is not supported.
func (c *_TMapConverter) ToMapWithDefault(value any, defaultValue map[string]any) map[string]any {
	return toMapWithDefault(value, defaultValue)
}

// ToNullableMap converts value into map object or returns null when conversion is not possible.
// Parameters: "value" - the value to convert
// Returns: map object and true or null and false when conversion is not supported.
func toNullableMap(value any) (map[string]any, bool) {
	if value == nil {
		return nil, false
	}

	v := reflect.ValueOf(value)

	switch v.Kind() {

	case reflect.Map:
		r := mapToMap(v)
		return r, true

	case reflect.Array, reflect.Slice:
		r := arrayToMap(v)
		return r, true

	case reflect.Struct:
		r := structToMap(v)
		return r, true

	case reflect.Interface, reflect.Ptr:
		if v.IsNil() {
			break
		}
		value = valueToInterface(v.Elem())
		return toNullableMap(value)
	}

	return nil, false
}

// ToMap converts value into map object or returns empty map when conversion is not possible.
// Parameters: "value" - the value to convert
// Returns: map object or empty map when conversion is not supported.
func toMap(value any) map[string]any {
	return toMapWithDefault(value, map[string]any{})
}

// ToMapWithDefault converts value into map object or returns default map when conversion is not possible.
// Parameters:
//  "value" - the value to convert
//  "defaultValue" - the default value.
// Returns: map object or default map when conversion is not supported.
func toMapWithDefault(value any, defaultValue map[string]any) map[string]any {
	if m, ok := toNullableMap(value); ok {
		return m
	}
	return defaultValue
}
