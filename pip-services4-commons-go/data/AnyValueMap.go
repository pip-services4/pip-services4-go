package data

import (
	"fmt"
	"strings"
	"time"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
)

// AnyValueMap is a Cross-language implementation of dynamic object map (dictionary)
// what can hold values of any type. The stored values can be converted to different
// types using variety of accessor methods.
//
//	Example
//		value1 := NewAnyValueMap(map[string]any{"key1": 1, "key2": "123.456", "key3": "2018-01-01"})
//
//		value1.GetAsBoolean("key1")  // Result: true
//		value1.GetAsInteger("key2")  // Result: 123
//		value1.GetAsFloat("key2")    // Result: 123.456
//		value1.GetAsDateTime("key3") // Result: new Date(2018,0,1)
//
//	see convert.StringConverter
//	see convert.TypeConverter
//	see convert.BooleanConverter
//	see convert.IntegerConverter
//	see convert.LongConverter
//	see convert.DoubleConverter
//	see convert.FloatConverter
//	see convert.DateTimeConverter
//	see ICloneable
type AnyValueMap struct {
	_value map[string]any
	_base  IMap
}

// NewEmptyAnyValueMap creates a new empty instance of the map.
//
//	Returns: *AnyValueMap
func NewEmptyAnyValueMap() *AnyValueMap {
	c := &AnyValueMap{
		_value: map[string]any{},
	}
	c._base = c
	return c
}

// InheritAnyValueMap creates a new instance of the map and assigns base methods from interface.
//
//	Parameters: base IMap
//	Returns: *AnyValueMap
func InheritAnyValueMap(base IMap) *AnyValueMap {
	c := &AnyValueMap{
		_value: make(map[string]any),
	}
	c._base = base
	return c
}

// NewAnyValueMap creates a new instance of the map and assigns its value.
// Parameters:
//   - values map[string]any
//
// Returns *AnyValueMap
func NewAnyValueMap(value map[string]any) *AnyValueMap {
	c := &AnyValueMap{
		_value: make(map[string]any),
	}
	c._base = c
	c.Append(value)
	return c
}

// NewAnyValueMapFromValue Converts specified value into AnyValueMap.
//
//	see SetAsSingleObject
//	Parameters: value any value to be converted
//	Returns: *AnyValueMap a newly created AnyValueMap.
func NewAnyValueMapFromValue(value any) *AnyValueMap {
	result := NewEmptyAnyValueMap()
	result.SetAsSingleObject(value)
	return result
}

// NewAnyValueMapFromTuples Creates a new AnyValueMap from a list of key-value pairs called tuples.
//
//	see NewAnyValueMapFromTuplesArray
//	Parameters: tuples ...tuples: any a list of values where odd elements
//		are keys and the following even elements are values
//	Returns: *AnyValueMap a newly created AnyValueArray.
func NewAnyValueMapFromTuples(tuples ...any) *AnyValueMap {
	return NewAnyValueMapFromTuplesArray(tuples)
}

// NewAnyValueMapFromTuplesArray creates a new AnyValueMap from a list of key-value pairs called tuples.
// The method is similar to fromTuples but tuples are passed as array instead of parameters.
//
//	Parameters: tuples: []any a list of values where odd elements are keys and
//	the following even elements are values
//	Returns: *AnyValueMap a newly created AnyValueArray.
func NewAnyValueMapFromTuplesArray(tuples []any) *AnyValueMap {
	result := NewEmptyAnyValueMap()
	if len(tuples) == 0 {
		return result
	}

	for index := 0; index < len(tuples); index = index + 2 {
		if index+1 >= len(tuples) {
			break
		}

		name := convert.StringConverter.ToString(tuples[index])
		value := tuples[index+1]

		result.SetAsObject(name, value)
	}

	return result
}

// NewAnyValueMapFromMaps creates a new AnyValueMap by merging two or more maps.
// Maps defined later in the list override values from previously defined maps.
//
//	Parameters: maps ...maps: any[] an array of maps to be merged
//	Returns: *AnyValueMap a newly created AnyValueMap.
func NewAnyValueMapFromMaps(maps ...map[string]any) *AnyValueMap {
	result := NewEmptyAnyValueMap()
	if len(maps) > 0 {
		for index := 0; index < len(maps); index++ {
			result.Append(maps[index])
		}
	}
	return result
}

// InnerValue return inner values of map as any
func (c *AnyValueMap) InnerValue() any {
	return c._value
}

// Value returns map of elements as map[string]any
func (c *AnyValueMap) Value() map[string]any {
	return c._value
}

// Get a map element specified by its key.
//
//	Parameters: key string a key of the element to get.
//	Returns: any the value of the map element.
func (c *AnyValueMap) Get(key string) (any, bool) {
	value, ok := c._value[key]
	return value, ok
}

// Keys gets keys of all elements stored in this map.
//
//	Returns: []string a list with all map keys.
func (c *AnyValueMap) Keys() []string {
	keys := make([]string, 0, len(c._value))
	for key := range c._value {
		keys = append(keys, key)
	}
	return keys
}

// Put a new value into map element specified by its key.
//
//	Parameters:
//		- key string a key of the element to put.
//		- value any a new value for map element.
//	Returns: any
func (c *AnyValueMap) Put(key string, value any) {
	c._value[key] = value
}

// Remove a map element specified by its key
//
//	Parameters: key string a key of the element to remove.
func (c *AnyValueMap) Remove(key string) {
	delete(c._value, key)
}

// Contains checks if this map contains a key. The check uses direct comparison between key and the specified key value.
//
//	Parameters: key string a value to be checked
//	Returns: bool true if this map contains the key or false otherwise.
func (c *AnyValueMap) Contains(key string) bool {
	_, ok := c._value[key]
	return ok
}

// Append new elements to this map.
//
//	Parameters: value: map[string]any a map of elements to be added.
func (c *AnyValueMap) Append(value map[string]any) {
	if value == nil {
		return
	}

	for key, val := range value {
		c._value[key] = val
	}
}

// Clear this map by removing all its elements.
func (c *AnyValueMap) Clear() {
	c._value = map[string]any{}
}

// Len gets a number of elements stored in this map.
//
//	Returns: int the number of elements in this map.
func (c *AnyValueMap) Len() int {
	return len(c._value)
}

// GetAsSingleObject gets the value stored in map element without any conversions.
// When element index is not defined it returns the entire array value.
//
//	Returns: any the element value or value of the array when index is not defined.
func (c *AnyValueMap) GetAsSingleObject() any {
	return c._value
}

// SetAsSingleObject sets a new value to map.
// see convert.MapConverter.ToMap
//
//	Parameters: value any a new element or array value.
func (c *AnyValueMap) SetAsSingleObject(value any) {
	c._value = convert.MapConverter.ToMap(value)
}

// GetAsObject gets the value stored in map element without any conversions.
// When element key is not defined it returns the entire map value.
//
//	Parameters: key string a key of the element to get
//	Returns: any the element value or value of the map when index is not defined.
func (c *AnyValueMap) GetAsObject(key string) (any, bool) {
	return c._base.Get(key)
}

// SetAsObject sets a new value to map element specified by its index. When the index is not defined,
// it resets the entire map value.
//
//	see convert.MapConverter.ToMap
//	Parameters:
//		- key string a key of the element to set
//		- value any a new element or map value.
func (c *AnyValueMap) SetAsObject(key string, value any) {
	c._base.Put(key, value)
}

// GetAsNullableString converts map element into a string or returns null if conversion is not possible.
// see convert.StringConverter.ToNullableString
//
//	Parameters: key string a key of element to get.
//	Returns: string value of the element or null if conversion is not supported.
func (c *AnyValueMap) GetAsNullableString(key string) (string, bool) {
	if value, ok := c._base.Get(key); ok && value != "" {
		return convert.StringConverter.ToNullableString(value)
	}
	return "", false
}

// GetAsString converts map element into a string or returns "" if conversion is not possible.
//
//	see GetAsStringWithDefault
//	Parameters: key string a key of element to get.
//	Returns: string value of the element or "" if conversion is not supported.
func (c *AnyValueMap) GetAsString(key string) string {
	return c.GetAsStringWithDefault(key, "")
}

// GetAsStringWithDefault converts map element into a string or returns default value if conversion is not possible.
// see convert.StringConverter.ToStringWithDefault
//
//	Parameters:
//		- key string a key of element to get.
//		- defaultValue string the default value
//	Returns: string value of the element or default value if conversion is not supported.
func (c *AnyValueMap) GetAsStringWithDefault(key string, defaultValue string) string {
	if value, ok := c._base.Get(key); ok && value != "" {
		return convert.StringConverter.ToStringWithDefault(value, defaultValue)
	}
	return defaultValue
}

// GetAsNullableBoolean converts map element into a boolean or returns null if conversion is not possible.
// see convert.BooleanConverter.toNullableBoolean
//
//	Parameters: key string a key of element to get.
//	Returns: bool value of the element or null if conversion is not supported.
func (c *AnyValueMap) GetAsNullableBoolean(key string) (bool, bool) {
	if value, ok := c._base.Get(key); ok {
		return convert.BooleanConverter.ToNullableBoolean(value)
	}
	return false, false
}

// GetAsBoolean converts map element into a boolean or returns false if conversion is not possible.
//
//	see GetAsBooleanWithDefault
//	Parameters: key: string a key of element to get.
//	Returns: bool value of the element or false if conversion is not supported.
func (c *AnyValueMap) GetAsBoolean(key string) bool {
	return c.GetAsBooleanWithDefault(key, false)
}

// GetAsBooleanWithDefault converts map element into a boolean or returns default value if conversion is not possible.
//
//	see convert.BooleanConverter.toBooleanWithDefault
//	Parameters:
//		- key string a key of element to get.
//		- defaultValue bool the default value
//	Returns: bool value of the element or default value if conversion is not supported.
func (c *AnyValueMap) GetAsBooleanWithDefault(key string, defaultValue bool) bool {
	if value, ok := c._base.Get(key); ok {
		return convert.BooleanConverter.ToBooleanWithDefault(value, defaultValue)
	}
	return defaultValue
}

// GetAsNullableInteger converts map element into an integer or returns null if conversion is not possible.
//
//	see convert.IntegerConverter.toNullableInteger
//	Parameters: key string a key of element to get.
//	Returns: value of the element or null if conversion is not supported.
func (c *AnyValueMap) GetAsNullableInteger(key string) (int, bool) {
	if value, ok := c._base.Get(key); ok {
		return convert.IntegerConverter.ToNullableInteger(value)
	}
	return 0, false
}

// GetAsInteger converts map element into an integer or returns 0 if conversion is not possible.
//
//	see GetAsIntegerWithDefault
//	Parameters: key string a key of element to get.
//	Returns: integer value of the element or 0 if conversion is not supported.
func (c *AnyValueMap) GetAsInteger(key string) int {
	return c.GetAsIntegerWithDefault(key, 0)
}

// GetAsIntegerWithDefault converts map element into an integer or returns default value if conversion is not possible.
//
//	see convert.IntegerConverter.toIntegerWithDefault
//	Parameters:
//		- key string a key of element to get.
//		-defaultValue int the default value
//	Returns: integer value of the element or default value if conversion is not supported.
func (c *AnyValueMap) GetAsIntegerWithDefault(key string, defaultValue int) int {
	if value, ok := c._base.Get(key); ok {
		return convert.IntegerConverter.ToIntegerWithDefault(value, defaultValue)
	}
	return defaultValue
}

// GetAsNullableUInteger converts map element into an unsigned integer or returns false if conversion is not possible.
//
//	see convert.IntegerConverter.toNullableUInteger
//	Parameters: key string a key of element to get.
//	Returns: value of the element or false if conversion is not supported.
func (c *AnyValueMap) GetAsNullableUInteger(key string) (uint, bool) {
	if value, ok := c._base.Get(key); ok {
		return convert.IntegerConverter.ToNullableUInteger(value)
	}
	return 0, false
}

// GetAsUInteger converts map element into an integer or returns 0 if conversion is not possible.
//
//	see GetAsIntegerWithDefault
//	Parameters: key string a key of element to get.
//	Returns: unsigned integer value of the element or 0 if conversion is not supported.
func (c *AnyValueMap) GetAsUInteger(key string) uint {
	return c.GetAsUIntegerWithDefault(key, 0)
}

// GetAsUIntegerWithDefault converts map element into an unsigned integer or
// returns default value if conversion is not possible.
//
//	see IntegerConverter.ToUIntegerWithDefault
//	Parameters:
//		- key string a key of element to get.
//		-defaultValue int the default value
//	Returns: unsigned integer value of the element or default value if conversion is not supported.
func (c *AnyValueMap) GetAsUIntegerWithDefault(key string, defaultValue uint) uint {
	if value, ok := c._base.Get(key); ok {
		return convert.IntegerConverter.ToUIntegerWithDefault(value, defaultValue)
	}
	return defaultValue
}

// GetAsNullableLong converts map element into a long or returns null if conversion is not possible.
//
//	see convert.LongConverter.toNullableLong
//	Parameters: key string a key of element to get.
//	Returns: int64 value of the element or null if conversion is not supported.
func (c *AnyValueMap) GetAsNullableLong(key string) (int64, bool) {
	if value, ok := c._base.Get(key); ok {
		return convert.LongConverter.ToNullableLong(value)
	}
	return 0, false
}

// GetAsLong converts map element into a long or returns 0 if conversion is not possible.
//
//	see GetAsLongWithDefault
//	Parameters: key string a key of element to get.
//	Returns: int64 value of the element or 0 if conversion is not supported.
func (c *AnyValueMap) GetAsLong(key string) int64 {
	return c.GetAsLongWithDefault(key, 0)
}

// GetAsLongWithDefault converts map element into a long or returns default value if conversion is not possible.
//
//	see convert.LongConverter.toLongWithDefault
//	Parameters:
//		-key string a key of element to get.
//		- defaultValue int64 the default value
//	Returns: int64 value of the element or default value if conversion is not supported.
func (c *AnyValueMap) GetAsLongWithDefault(key string, defaultValue int64) int64 {
	if value, ok := c._base.Get(key); ok {
		return convert.LongConverter.ToLongWithDefault(value, defaultValue)
	}
	return defaultValue
}

// GetAsNullableULong converts map element into a long or returns null if conversion is not possible.
//
//	see convert.LongConverter.toNullableLong
//	Parameters: key string a key of element to get.
//	Returns: int64 value of the element or null if conversion is not supported.
func (c *AnyValueMap) GetAsNullableULong(key string) (uint64, bool) {
	if value, ok := c._base.Get(key); ok {
		return convert.LongConverter.ToNullableULong(value)
	}
	return 0, false
}

// GetAsULong converts map element into a unsigned long or returns 0 if conversion is not possible.
//
//	see GetAsLongWithDefault
//	Parameters: key string a key of element to get.
//	Returns: uint64 value of the element or 0 if conversion is not supported.
func (c *AnyValueMap) GetAsULong(key string) uint64 {
	return c.GetAsULongWithDefault(key, 0)
}

// GetAsULongWithDefault converts map element into a unsigned long or
// returns default value if conversion is not possible.
//
//	see convert.LongConverter.ToLongWithDefault
//
// Parameters:
//   - key string a key of element to get.
//   - defaultValue int64 the default value
//     Returns: uint64 value of the element or default value if conversion is not supported.
func (c *AnyValueMap) GetAsULongWithDefault(key string, defaultValue uint64) uint64 {
	if value, ok := c._base.Get(key); ok {
		return convert.LongConverter.ToULongWithDefault(value, defaultValue)
	}
	return defaultValue
}

// GetAsNullableFloat converts map element into a float or returns null if conversion is not possible.
//
//	see convert.FloatConverter.toNullableFloat
//	Parameters: key string a key of element to get.
//	Returns: float32 value of the element or null if conversion is not supported.
func (c *AnyValueMap) GetAsNullableFloat(key string) (float32, bool) {
	if value, ok := c._base.Get(key); ok {
		return convert.FloatConverter.ToNullableFloat(value)
	}
	return 0, false
}

// GetAsFloat converts map element into a float or returns 0 if conversion is not possible.
//
//	see GetAsFloatWithDefault
//	Parameters: key string a key of element to get.
//	Returns: float32 value of the element or 0 if conversion is not supported.
func (c *AnyValueMap) GetAsFloat(key string) float32 {
	return c.GetAsFloatWithDefault(key, 0)
}

// GetAsFloatWithDefault converts map element into a flot or returns default value if conversion is not possible.
//
//	see convert.FloatConverter.toFloatWithDefault
//	Parameters:
//		- key string a key of element to get.
//		- defaultValue float32 the default value
//	Returns: float32 value of the element or default value if conversion is not supported.
func (c *AnyValueMap) GetAsFloatWithDefault(key string, defaultValue float32) float32 {
	if value, ok := c._base.Get(key); ok {
		return convert.FloatConverter.ToFloatWithDefault(value, defaultValue)
	}
	return defaultValue
}

// GetAsNullableDouble converts map element into a double or returns null if conversion is not possible.
//
//	see convert.DoubleConverter.toNullableDouble
//	Parameters: key string a key of element to get.
//	Returns: float64 value of the element or null if conversion is not supported.
func (c *AnyValueMap) GetAsNullableDouble(key string) (float64, bool) {
	if value, ok := c._base.Get(key); ok {
		return convert.DoubleConverter.ToNullableDouble(value)
	}
	return 0, false
}

// GetAsDouble converts map element into a double or returns 0 if conversion is not possible.
//
//	see GetAsDoubleWithDefault
//	Parameters: key string a key of element to get.
//	Returns: float64 value of the element or 0 if conversion is not supported.
func (c *AnyValueMap) GetAsDouble(key string) float64 {
	return c.GetAsDoubleWithDefault(key, 0)
}

// GetAsDoubleWithDefault converts map element into a double or returns default value if conversion is not possible.
// see convert.DoubleConverter.ToDoubleWithDefault
//
//	Parameters:
//		- key string a key of element to get.
//		- defaultValue float64 the default value
//	Returns: float64 value of the element or default value if conversion is not supported.
func (c *AnyValueMap) GetAsDoubleWithDefault(key string, defaultValue float64) float64 {
	if value, ok := c._base.Get(key); ok {
		return convert.DoubleConverter.ToDoubleWithDefault(value, defaultValue)
	}
	return defaultValue
}

// GetAsNullableDateTime converts map element into a time.Time or returns null if conversion is not possible.
//
//	see convert.DateTimeConverter.ToNullableDateTime
//	Parameters: key string a key of element to get.
//	Returns: time.Time value of the element or null if conversion is not supported.
func (c *AnyValueMap) GetAsNullableDateTime(key string) (time.Time, bool) {
	if value, ok := c._base.Get(key); ok {
		return convert.DateTimeConverter.ToNullableDateTime(value)
	}
	return time.Time{}, false
}

// GetAsDateTime converts map element into a time.Time or returns the current date if conversion is not possible.
//
//	see GetAsDateTimeWithDefault
//	Parameters: key string a key of element to get.
//	Returns: time.Time value of the element or the current date if conversion is not supported.
func (c *AnyValueMap) GetAsDateTime(key string) time.Time {
	return c.GetAsDateTimeWithDefault(key, time.Time{})
}

// GetAsDateTimeWithDefault converts map element into a time.Time or returns default value if conversion is not possible.
//
//	see convert.DateTimeConverter.toDateTimeWithDefault
//	Parameters:
//		- key: string a key of element to get.
//		- defaultValue: Date the default value
//	Returns: time.Time value of the element or default value if conversion is not supported.
func (c *AnyValueMap) GetAsDateTimeWithDefault(key string, defaultValue time.Time) time.Time {
	if value, ok := c._base.Get(key); ok {
		return convert.DateTimeConverter.ToDateTimeWithDefault(value, defaultValue)
	}
	return defaultValue
}

// GetAsNullableDuration converts map element into a time.Duration or returns null if conversion is not possible.
//
//	see convert.DurationConverter.toNullableDateTime
//	Parameters: key string a key of element to get.
//	Returns: time.Duration value of the element or null if conversion is not supported.
func (c *AnyValueMap) GetAsNullableDuration(key string) (time.Duration, bool) {
	if value, ok := c._base.Get(key); ok {
		return convert.DurationConverter.ToNullableDuration(value)
	}
	return 0, false
}

// GetAsDuration converts map element into a time.Duration or returns the current date if conversion is not possible.
//
//	see GetAsDurationWithDefault
//	Parameters: key string a key of element to get.
//	Returns: time.Duration value of the element or the current date if conversion is not supported.
func (c *AnyValueMap) GetAsDuration(key string) time.Duration {
	return c.GetAsDurationWithDefault(key, 0*time.Millisecond)
}

// GetAsDurationWithDefault converts map element into a time.Duration or returns default value if conversion is not possible.
//
//	see convert.DurationConverter.toDateTimeWithDefault
//	Parameters:
//		- key: string a key of element to get.
//		- defaultValue: Date the default value
//	Returns: time.Duration value of the element or default value if conversion is not supported.
func (c *AnyValueMap) GetAsDurationWithDefault(key string, defaultValue time.Duration) time.Duration {
	if value, ok := c._base.Get(key); ok {
		return convert.DurationConverter.ToDurationWithDefault(value, defaultValue)
	}
	return defaultValue
}

// GetAsNullableType converts map element into a value defined by specied typecode.
// If conversion is not possible it returns null.
//
//	see TypeConverter.ToNullableType
//	Parameters:
//		- type TypeCode the TypeCode that defined the type of the result
//		- key string a key of element to get.
//	Returns: any element value defined by the typecode or null if conversion is not supported.
func (c *AnyValueMap) GetAsNullableType(typ convert.TypeCode, key string) (any, bool) {
	if value, ok := c._base.Get(key); ok {
		return convert.TypeConverter.ToNullableType(typ, value)
	}
	return nil, false
}

// GetAsType converts map element into a value defined by specied typecode.
// If conversion is not possible it returns default value for the specified type.
//
//	see GetAsTypeWithDefault
//	Parameters:
//		- type TypeCode the TypeCode that defined the type of the result
//		- key string a key of element to get.
//	Returns: any element value defined by the typecode or default if conversion is not supported.
func (c *AnyValueMap) GetAsType(typ convert.TypeCode, key string) any {
	return c.GetAsTypeWithDefault(typ, key, nil)
}

// GetAsTypeWithDefault converts map element into a value defined by specied typecode.
// If conversion is not possible it returns default value.
//
//	see convert.TypeConverter.toTypeWithDefault
//	Parameters:
//		- type TypeCode the TypeCode that defined the type of the result
//		- key string a key of element to get.
//		- defaultValue any the default value
//	Returns: any element value defined by the typecode or default value if conversion is not supported.
func (c *AnyValueMap) GetAsTypeWithDefault(typ convert.TypeCode, key string, defaultValue any) any {
	if value, ok := c._base.Get(key); ok {
		return convert.TypeConverter.ToTypeWithDefault(typ, value, defaultValue)
	}
	return defaultValue
}

// GetAsValue converts map element into an AnyValue or returns an empty AnyValue if conversion is not possible.
//
//	see AnyValue
//	see NewAnyValue
//	Parameters: key string a key of element to get.
//	Returns: *AnyValue value of the element or empty AnyValue if conversion is not supported.
func (c *AnyValueMap) GetAsValue(key string) *AnyValue {
	if value, ok := c._base.Get(key); ok {
		return NewAnyValue(value)
	}
	return NewEmptyAnyValue()
}

// GetAsNullableArray converts map element into an AnyValueArray or returns null if conversion is not possible.
//
//	see AnyValueArray
//	see NewAnyValueArrayFromValue
//	Parameters:  key string a key of element to get.
//	Returns: *AnyValueArray value of the element or null if conversion is not supported.
func (c *AnyValueMap) GetAsNullableArray(key string) (*AnyValueArray, bool) {
	if value, ok := c._base.Get(key); ok {
		return NewAnyValueArrayFromValue(value), true
	}
	return nil, false
}

// GetAsArray converts map element into an AnyValueArray or returns empty AnyValueArray if conversion is not possible.
//
//	see AnyValueArray
//	see NewAnyValueArrayFromValue
//	Parameters: key string a key of element to get.
//	Returns: *AnyValueArray value of the element or empty AnyValueArray if conversion is not supported.
func (c *AnyValueMap) GetAsArray(key string) *AnyValueArray {
	if value, ok := c._base.Get(key); ok {
		return NewAnyValueArrayFromValue(value)
	}
	return NewEmptyAnyValueArray()
}

// GetAsArrayWithDefault converts map element into an AnyValueArray or returns default value if conversion is not possible.
//
//	see AnyValueArray
//	see GetAsNullableArray
//	Parameters:
//		- key string a key of element to get.
//		- defaultValue: *AnyValueArray the default value
//	Returns: *AnyValueArray value of the element or default value if conversion is not supported.
func (c *AnyValueMap) GetAsArrayWithDefault(key string, defaultValue *AnyValueArray) *AnyValueArray {
	if result, ok := c.GetAsNullableArray(key); ok {
		return result
	}
	return defaultValue
}

// GetAsNullableMap converts map element into an AnyValueMap or returns null if conversion is not possible.
//
//	see NewAnyValueMapFromValue
//	Parameters: key string  a key of element to get.
//	Returns: *AnyValueMap value of the element or null if conversion is not supported.
func (c *AnyValueMap) GetAsNullableMap(key string) (*AnyValueMap, bool) {
	if value, ok := c._base.Get(key); ok {
		return NewAnyValueMapFromValue(value), true
	}
	return nil, false
}

// GetAsMap converts map element into an AnyValueMap or returns empty AnyValueMap if conversion is not possible.
//
//	see NewAnyValueMapFromValue
//	Parameters: key string a key of element to get.
//	Returns: *AnyValueMap value of the element or empty AnyValueMap if conversion is not supported.
func (c *AnyValueMap) GetAsMap(key string) *AnyValueMap {
	if value, ok := c._base.Get(key); ok {
		return NewAnyValueMapFromValue(value)
	}
	return NewEmptyAnyValueMap()
}

// GetAsMapWithDefault converts map element into an AnyValueMap or returns default value if conversion is not possible.
//
//	see GetAsNullableMap
//	Parameters:
//		- key string a key of element to get.
//		- defaultValue *AnyValueMap the default value
//	Returns: *AnyValueMap value of the element or default value if conversion is not supported.
func (c *AnyValueMap) GetAsMapWithDefault(key string, defaultValue *AnyValueMap) *AnyValueMap {
	if result, ok := c.GetAsNullableMap(key); ok {
		return result
	}
	return defaultValue
}

// String gets a string representation of the object.
// The result is a semicolon-separated list of key-value pairs as "key1=value1;key2=value2;key=value3"
// Returns: string a string representation of the object.
func (c *AnyValueMap) String() string {
	builder := strings.Builder{}

	// Todo:: User encoder
	for key := range c.Value() {
		value, _ := c._base.Get(key)

		if builder.Len() > 0 {
			builder.WriteString(";")
		}

		if value != nil {
			builder.WriteString(fmt.Sprintf("%s=%v", key, value))
		} else {
			builder.WriteString(key)
		}
	}

	return builder.String()
}

// Clone creates a binary clone of this object.
//
//	Returns: any a clone of this object.
func (c *AnyValueMap) Clone() *AnyValueMap {
	return NewAnyValueMap(c._value)
}
