package convert

import (
	"reflect"
	"time"
)

// TypeConverter converts arbitrary values into objects specific by TypeCodes.
// For each TypeCode this class calls corresponding converter
// which applies extended conversion rules to convert the values.
//
// Example:
//
//  value1 := convert.TypeConverter.ToType(convert.Integer, "123.456")
//  value2 := convert.TypeConverter.ToType(convert.DateTime, 123)
//  value3 := convert.TypeConverter.ToType(convert.Boolean, "F")
//  fmt.Println(value1) // 123
//  fmt.Println(value2) // 1970-01-01 02:02:03 +0200 EET
//  fmt.Println(value3) // false
var TypeConverter = &_TTypeConverter{}

type _TTypeConverter struct{}

// ToTypeCode gets TypeCode for specific value.
// Parameters: "value" - value whose TypeCode is to be resolved.
// Returns: the TypeCode that corresponds to the passed object's type.
func (c *_TTypeConverter) ToTypeCode(value any) TypeCode {
	return toTypeCode(value)
}

// ToNullableType converts value into an object type specified by Type Code or returns null
// when conversion is not possible.
// Parameters:
//  "typ" - the TypeCode for the data type.
//  "value" - the value to convert.
// Returns: object value of type corresponding to TypeCode, or null when
// conversion is not supported.
func (c *_TTypeConverter) ToNullableType(typ TypeCode, value any) (any, bool) {
	return toNullableType(typ, value)
}

// ToType converts value into an object type specified by Type Code
// or returns default value when conversion is not possible.
// Parameters:
//  "typ" - the TypeCode for the data type into which 'value' is to be converted.
//  "value" - the value to convert.
// Returns: object value of type corresponding to TypeCode, or default value when
// conversion is not supported
func (c *_TTypeConverter) ToType(typ TypeCode, value any) any {
	return toType(typ, value)
}

// ToTypeWithDefault converts value into an object type specified by Type Code
// or returns default value when conversion is not possible.
// Parameters:
//  "typ" - the TypeCode for the data type into which 'value' is to be converted.
//  "value" - the value to convert.
//  "defaultValue" - the default value to return if conversion is not possible
//  (returns null).
// Returns: object value of type corresponding to TypeCode, or default value when
// conversion is not supported
func (c *_TTypeConverter) ToTypeWithDefault(typ TypeCode, value any, defaultValue any) any {
	return toTypeWithDefault(typ, value, defaultValue)
}

// ToString converts a TypeCode into its string name.
// Parameters: "typ" - the TypeCode to convert into a string.
// Returns: the name of the TypeCode passed as a string value.
func (c *_TTypeConverter) ToString(typ TypeCode) string {
	return typeCodeToString(typ)
}

// ToTypeCode gets TypeCode for specific value.
// Parameters: "value" - value whose TypeCode is to be resolved.
// Returns: the TypeCode that corresponds to the passed object's type.
func toTypeCode(value any) TypeCode {
	if value == nil {
		return Unknown
	}

	switch value.(type) {
	case string:
		return String

	case bool:
		return Boolean

	case byte, uint, int, int32:
		return Integer

	case uint32, uint64, int64:
		return Long

	case float32:
		return Float

	case float64:
		return Double

	case time.Time:
		return DateTime

	case time.Duration:
		return Duration
	}

	rt, ok := value.(reflect.Type)
	if !ok {
		rt = reflect.TypeOf(value)
	}
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	if rt == reflect.TypeOf((*time.Time)(nil)).Elem() {
		return DateTime
	}

	if rt == reflect.TypeOf((*time.Duration)(nil)).Elem() {
		return Duration
	}

	switch rt.Kind() {
	case reflect.String:
		return String

	case reflect.Bool:
		return Boolean

	case reflect.Int8, reflect.Uint8, reflect.Int16, reflect.Uint16,
		reflect.Int32, reflect.Uint32, reflect.Int, reflect.Uint:
		return Integer

	case reflect.Int64, reflect.Uint64:
		return Long

	case reflect.Float32:
		return Float

	case reflect.Float64:
		return Double

	case reflect.Struct:
		return Object

	case reflect.Map:
		return Map

	case reflect.Array, reflect.Slice:
		return Array

	default:
		return Unknown
	}
}

// ToNullableType converts value into an object type specified by Type Code or returns null
// when conversion is not possible.
// Parameters:
//  "typ" - the TypeCode for the data type.
//  "value" - the value to convert.
// Returns: object value of type corresponding to TypeCode and true, or null and false when
// conversion is not supported.
func toNullableType(typ TypeCode, value any) (any, bool) {
	if value == nil {
		return nil, false
	}
	// Convert to known types
	switch typ {
	case String:
		return StringConverter.ToNullableString(value)
	case Boolean:
		return BooleanConverter.ToNullableBoolean(value)
	case Integer:
		return IntegerConverter.ToNullableInteger(value)
	case Long:
		return LongConverter.ToNullableLong(value)
	case Float:
		return FloatConverter.ToNullableFloat(value)
	case Double:
		return DoubleConverter.ToNullableDouble(value)
	case DateTime:
		return DateTimeConverter.ToNullableDateTime(value)
	case Duration:
		return DurationConverter.ToNullableDuration(value)
	case Array:
		return ArrayConverter.ToNullableArray(value)
	case Map:
		return MapConverter.ToNullableMap(value)
	default:
		return nil, false
	}
}

// ToType converts value into an object type specified by Type Code
// or returns default value when conversion is not possible.
// Parameters:
//  "typ" - the TypeCode for the data type into which 'value' is to be converted.
//  "value" - the value to convert.
// Returns: object value of type corresponding to TypeCode, or default value when
// conversion is not supported
func toType(typ TypeCode, value any) any {
	if value == nil {
		return nil
	}

	// Convert to known types
	switch typ {
	case String:
		return StringConverter.ToString(value)
	case Boolean:
		return BooleanConverter.ToBoolean(value)
	case Integer:
		return IntegerConverter.ToInteger(value)
	case Long:
		return LongConverter.ToLong(value)
	case Float:
		return FloatConverter.ToFloat(value)
	case Double:
		return DoubleConverter.ToDouble(value)
	case DateTime:
		return DateTimeConverter.ToDateTime(value)
	case Duration:
		return DurationConverter.ToDuration(value)
	case Array:
		return ArrayConverter.ToArray(value)
	case Map:
		return MapConverter.ToMap(value)
	default:
		return value
	}
}

// ToTypeWithDefault converts value into an object type specified by Type Code
// or returns default value when conversion is not possible.
// Parameters:
//  "typ" - the TypeCode for the data type into which 'value' is to be converted.
//  "value" - the value to convert.
//  "defaultValue" - the default value to return if conversion is not possible
//  (returns null).
// Returns: object value of type corresponding to TypeCode, or default value when
// conversion is not supported
func toTypeWithDefault(typ TypeCode, value any, defaultValue any) any {
	if value == nil {
		return defaultValue
	}

	// Convert to known types
	switch typ {
	case String:
		defVal, _ := defaultValue.(string)
		return StringConverter.ToStringWithDefault(value, defVal)
	case Boolean:
		defVal, _ := defaultValue.(bool)
		return BooleanConverter.ToBooleanWithDefault(value, defVal)
	case Integer:
		defVal, _ := defaultValue.(int)
		return IntegerConverter.ToIntegerWithDefault(value, defVal)
	case Long:
		defVal, _ := defaultValue.(int64)
		return LongConverter.ToLongWithDefault(value, defVal)
	case Float:
		defVal, _ := defaultValue.(float32)
		return FloatConverter.ToFloatWithDefault(value, defVal)
	case Double:
		defVal, _ := defaultValue.(float64)
		return DoubleConverter.ToDoubleWithDefault(value, defVal)
	case DateTime:
		defVal, _ := defaultValue.(time.Time)
		return DateTimeConverter.ToDateTimeWithDefault(value, defVal)
	case Duration:
		defVal, _ := defaultValue.(time.Duration)
		return DurationConverter.ToDurationWithDefault(value, defVal)
	case Array:
		defVal, _ := defaultValue.([]any)
		return ArrayConverter.ToArrayWithDefault(value, defVal)
	case Map:
		defVal, _ := defaultValue.(map[string]any)
		return MapConverter.ToMapWithDefault(value, defVal)
	default:
		return defaultValue
	}
}

// TypeCodeToString converts a TypeCode into its string name.
// Parameters: "typ" - the TypeCode to convert into a string.
// Returns: the name of the TypeCode passed as a string value.
func typeCodeToString(typ TypeCode) string {
	switch typ {
	case Unknown:
		return "unknown"
	case String:
		return "string"
	case Boolean:
		return "boolean"
	case Integer:
		return "integer"
	case Long:
		return "long"
	case Float:
		return "float"
	case Double:
		return "double"
	case DateTime:
		return "datetime"
	case Duration:
		return "duration"
	case Object:
		return "object"
	case Enum:
		return "enum"
	case Array:
		return "array"
	case Map:
		return "map"
	default:
		return "unknown"
	}
}
