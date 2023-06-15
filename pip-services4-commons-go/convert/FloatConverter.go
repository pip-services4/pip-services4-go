package convert

// FloatConverter Converts arbitrary values into float using extended conversion rules:
// - Strings are converted to float values
// - DateTime: total number of milliseconds since unix epoсh
// - Boolean: 1 for true and 0 for false
//
// Example:
//
//  value1, ok1 := convert.FloatConverter.ToNullableFloat("ABC")
//  value2, ok2 := convert.FloatConverter.ToNullableFloat("123.456")
//  value3, ok3 := convert.FloatConverter.ToNullableFloat(true)
//  value4, ok4 := convert.FloatConverter.ToNullableFloat(time.Now())
//  fmt.Println(value1, ok1) // 0, false
//  fmt.Println(value2, ok2) // 123.456, true
//  fmt.Println(value3, ok3) // 1, true
//  fmt.Println(value4, ok4) // current milliseconds (e.g. 1.566333114e+09), true
var FloatConverter = &_TFloatConverter{}

type _TFloatConverter struct{}

// ToNullableFloat converts value into float or returns null when conversion is not possible.
// Parameters: "value" - the value to convert.
// Returns: float value and true or 0 and false when conversion is not supported.
func (c *_TFloatConverter) ToNullableFloat(value any) (float32, bool) {
	return toNullableFloat(value)
}

// ToFloat converts value into float or returns 0 when conversion is not possible.
// Parameters: "value" - the value to convert.
// Returns: float value or 0 when conversion is not supported.
func (c *_TFloatConverter) ToFloat(value any) float32 {
	return toFloat(value)
}

// ToFloatWithDefault converts value into float or returns default when conversion is not possible.
// Parameters:
//  "value" - the value to convert.
//  "defaultValue" - the default value.
// Returns: float value or default when conversion is not supported.
func (c *_TFloatConverter) ToFloatWithDefault(value any, defaultValue float32) float32 {
	return toFloatWithDefault(value, defaultValue)
}

// ToNullableFloat converts value into float or returns null when conversion is not possible.
// Parameters: "value" - the value to convert.
// Returns: float value and true or 0 and false when conversion is not supported.
func toNullableFloat(value any) (float32, bool) {
	if r, ok := DoubleConverter.ToNullableDouble(value); ok {
		return float32(r), ok
	}
	return 0, false
}

// ToFloat converts value into float or returns 0 when conversion is not possible.
// Parameters: "value" - the value to convert.
// Returns: float value or 0 when conversion is not supported.
func toFloat(value any) float32 {
	return toFloatWithDefault(value, 0)
}

// ToFloatWithDefault converts value into float or returns default when conversion is not possible.
// Parameters:
//  "value" - the value to convert.
//  "defaultValue" - the default value.
// Returns: float value or default when conversion is not supported.
func toFloatWithDefault(value any, defaultValue float32) float32 {
	if r, ok := toNullableFloat(value); ok {
		return r
	}
	return defaultValue
}
