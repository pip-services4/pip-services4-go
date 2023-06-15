package data

import (
	"fmt"
	"strings"
	"time"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
)

// StringValueMap Cross-language implementation of a map (dictionary) where all keys and values are strings.
// The stored values can be converted to different types using variety of accessor methods.
// The string map is highly versatile. It can be converted into many formats, stored and sent over the wire.
//
// This class is widely used in Pip.Services as a basis for variety of classes, such as ConfigParams, ConnectionParams,
// CredentialParams and others.
//
//	Example:
//		value1 := NewStringValueMapFromString("key1=1;key2=123.456;key3=2018-01-01")
//
//		value1.GetAsBoolean("key1")   // Result: true
//		value1.GetAsInteger("key2")   // Result: 123
//		value1.GetAsFloat("key2")     // Result: 123.456
//		value1.GetAsDateTime("key3")  // Result: new Date(2018,0,1)
//		see StringConverter
//		see TypeConverter
//		see BooleanConverter
//		see IntegerConverter
//		see LongConverter
//		see DoubleConverter
//		see FloatConverter
//		see DateTimeConverter
type StringValueMap struct {
	_value map[string]string
}

// NewEmptyStringValueMap creates a new instance of the map.
func NewEmptyStringValueMap() *StringValueMap {
	return &StringValueMap{
		_value: make(map[string]string),
	}
}

// NewStringValueMap creates a new instance of the map and assigns its value.
//
//	Parameters: value map[string]string
//	Returns: *StringValueMap
func NewStringValueMap(value map[string]string) *StringValueMap {
	c := NewEmptyStringValueMap()
	c.Append(value)
	return c
}

// NewStringValueMapFromValue converts specified value into StringValueMap.
//
//	see SetAsSingleObject
//	Parameters: value any value to be converted
//	Returns: *StringValueMap a newly created StringValueMap.
func NewStringValueMapFromValue(value any) *StringValueMap {
	result := NewEmptyStringValueMap()
	result.SetAsSingleObject(value)
	return result
}

// NewStringValueMapFromTuples creates a new StringValueMap from a list of key-value pairs called tuples.
//
//	see NewStringValueMapFromTuplesArray
//	Parameters: tuples ...any a list of values where odd elements
//		are keys and the following even elements are values
//	Returns: *StringValueMap a newly created StringValueMap.
func NewStringValueMapFromTuples(tuples ...any) *StringValueMap {
	return NewStringValueMapFromTuplesArray(tuples)
}

// NewStringValueMapFromTuplesArray creates a new StringValueMap from a list of key-value pairs called tuples.
// The method is similar to fromTuples but tuples are passed as array instead of parameters.
//
//	Parameters: tuples: []any a list of values where odd elements
//		are keys and the following even elements are values
//	Returns: *StringValueMap a newly created StringValueMap.
func NewStringValueMapFromTuplesArray(tuples []any) *StringValueMap {
	result := NewEmptyStringValueMap()
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

// NewStringValueMapFromString parses semicolon-separated key-value pairs and returns them as a StringValueMap.
//
//	Parameters: line string semicolon-separated key-value list to initialize StringValueMap.
//	Returns: *StringValueMap a newly created StringValueMap.
func NewStringValueMapFromString(line string) *StringValueMap {
	result := NewEmptyStringValueMap()
	if line == "" {
		return result
	}

	// Todo: User tokenizer / decoder
	tokens := strings.Split(line, ";")

	for index := 0; index < len(tokens); index++ {
		token := tokens[index]
		if len(token) == 0 {
			continue
		}

		pos := strings.Index(token, "=")

		var key string
		if pos > 0 {
			key = token[0:pos]
			key = strings.TrimSpace(key)
		} else {
			key = strings.TrimSpace(token)
		}

		var value string
		if pos > 0 {
			value = token[pos+1:]
			value = strings.TrimSpace(value)
		} else {
			value = ""
		}

		result.Put(key, value)
	}

	return result
}

// NewStringValueMapFromMaps creates a new AnyValueMap by merging two or more maps.
// Maps defined later in the list override values from previously defined maps.
//
//	Parameters: maps...map[string]string an array of maps to be merged
//	Returns: StringValueMap a newly created AnyValueMap.
func NewStringValueMapFromMaps(maps ...map[string]string) *StringValueMap {
	result := NewEmptyStringValueMap()
	if len(maps) > 0 {
		for index := 0; index < len(maps); index++ {
			result.Append(maps[index])
		}
	}
	return result
}

// InnerValue return inner values of map as any
func (c *StringValueMap) InnerValue() any {
	return c._value
}

// Value returns map of elements as map[string]any
func (c *StringValueMap) Value() map[string]string {
	return c._value
}

// Get a map element specified by its key.
//
//	Parameters: key string a key of the element to get.
//	Returns string the value of the map element.
func (c *StringValueMap) Get(key string) (any, bool) {
	value, ok := c._value[key]
	return value, ok
}

// Keys gets keys of all elements stored in this map.
//
//	Returns: []string a list with all map keys.
func (c *StringValueMap) Keys() []string {
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
func (c *StringValueMap) Put(key string, value any) {
	c._value[key] = convert.StringConverter.ToString(value)
}

// Remove a map element specified by its key
//
//	Parameters: key string a key of the element to remove.
func (c *StringValueMap) Remove(key string) {
	delete(c._value, key)
}

// Contains checks if this map contains a key. The check uses direct comparison between key and the specified key value.
// Parameters
//   - key string
//     a value to be checked
//
// Returns bool
// true if this map contains the key or false otherwise.
func (c *StringValueMap) Contains(key string) bool {
	_, ok := c._value[key]
	return ok
}

// Append new elements to this map.
//
//	Parameters: values map[string]string a map with elements to be added.
func (c *StringValueMap) Append(values map[string]string) {
	if values == nil {
		return
	}

	for key, val := range values {
		c._value[key] = val
	}
}

// AppendAny new elements to this map.
//
//	Parameters: values map[string]any a map with elements to be added.
func (c *StringValueMap) AppendAny(values map[string]any) {
	if values == nil {
		return
	}

	for key, val := range values {
		c._value[key] = convert.StringConverter.ToString(val)
	}
}

// Clear this map by removing all its elements.
func (c *StringValueMap) Clear() {
	c._value = make(map[string]string)
}

// Len gets a number of elements stored in this map.
//
//	Returns: int the number of elements in this map.
func (c *StringValueMap) Len() int {
	return len(c._value)
}

// GetAsSingleObject Gets the value stored in map element without any conversions.
// When element index is not defined it returns the entire array value.
//
//	Returns: any the element value or value of the array when index is not defined.
func (c *StringValueMap) GetAsSingleObject() any {
	return *c
}

// SetAsSingleObject sets a new value to map.
//
//	Parameters: value any a new element or array value.
func (c *StringValueMap) SetAsSingleObject(value any) {
	a := convert.MapConverter.ToMap(value)
	c.Clear()
	c.AppendAny(a)
}

// GetAsObject gets the value stored in map element without any conversions.
// When element key is not defined it returns the entire map value.
//
//	Parameters: key string a key of the element to get
//	Returns: any the element value or value of the map when index is not defined.
func (c *StringValueMap) GetAsObject(key string) (any, bool) {
	return c.Get(key)
}

// SetAsObject sets a new value to map element specified by its index.
// When the index is not defined, it resets the entire map value.
//
//	Parameters:
//		- key any a key of the element to set
//		- value any a new element or map value.
func (c *StringValueMap) SetAsObject(key string, value any) {
	c.Put(key, value)
}

// GetAsNullableString converts map element into a string or returns null if conversion is not possible.
//
//	see StringConverter.ToNullableString
//	Parameters: key string a key of element to get.
//	Returns: string value of the element or null if conversion is not supported.
func (c *StringValueMap) GetAsNullableString(key string) (string, bool) {
	if value, ok := c.Get(key); ok && value != "" {
		return convert.StringConverter.ToNullableString(value)
	}
	return "", false
}

// GetAsString converts map element into a string or returns "" if conversion is not possible.
//
//	see GetAsStringWithDefault
//	Parameters: key string a key of element to get.
//	Returns: string value of the element or "" if conversion is not supported.
func (c *StringValueMap) GetAsString(key string) string {
	return c.GetAsStringWithDefault(key, "")
}

// GetAsStringWithDefault converts map element into a string or returns default value if conversion is not possible.
//
//	see convert.StringConverter.toStringWithDefault
//	Parameters:
//		- key string a key of element to get.
//		- defaultValue string the default value
//	Returns: string value of the element or default value if conversion is not supported.
func (c *StringValueMap) GetAsStringWithDefault(key string, defaultValue string) string {
	if value, ok := c.Get(key); ok && value != "" {
		return convert.StringConverter.ToStringWithDefault(value, defaultValue)
	}
	return defaultValue
}

// GetAsNullableBoolean converts map element into a boolean or returns null if conversion is not possible.
//
//	see convert.BooleanConverter.toNullableBoolean
//	Parameters:  key string a key of element to get.
//	Returns: boolean value of the element or null if conversion is not supported.
func (c *StringValueMap) GetAsNullableBoolean(key string) (bool, bool) {
	if value, ok := c.Get(key); ok {
		return convert.BooleanConverter.ToNullableBoolean(value)
	}
	return false, false
}

// GetAsBoolean converts map element into a boolean or returns false if conversion is not possible.
//
//	see GetAsBooleanWithDefault
//	Parameters: key string a key of element to get.
//	Returns: boolean value of the element or false if conversion is not supported.
func (c *StringValueMap) GetAsBoolean(key string) bool {
	return c.GetAsBooleanWithDefault(key, false)
}

// GetAsBooleanWithDefault converts map element into a boolean or returns default value if conversion is not possible.
// see
// BooleanConverter.toBooleanWithDefault
// Parameters
//   - key string
//     a key of element to get.
//   - defaultValue bool
//     the default value
//
// Returns bool
// boolean value of the element or default value if conversion is not supported.
func (c *StringValueMap) GetAsBooleanWithDefault(key string, defaultValue bool) bool {
	if value, ok := c.Get(key); ok {
		return convert.BooleanConverter.ToBooleanWithDefault(value, defaultValue)
	}
	return defaultValue
}

// GetAsNullableInteger converts map element into an integer or returns null if conversion is not possible.
//
//	see convert.IntegerConverter.toNullableInteger
//	Parameters: key string a key of element to get.
//	Returns: integer value of the element or null if conversion is not supported.
func (c *StringValueMap) GetAsNullableInteger(key string) (int, bool) {
	if value, ok := c.Get(key); ok {
		return convert.IntegerConverter.ToNullableInteger(value)
	}
	return 0, false
}

// GetAsInteger converts map element into an integer or returns 0 if conversion is not possible.
//
//	see GetAsIntegerWithDefault
//	Parameters: key string a key of element to get.
//	Returns: int
func (c *StringValueMap) GetAsInteger(key string) int {
	return c.GetAsIntegerWithDefault(key, 0)
}

// GetAsIntegerWithDefault converts map element into an integer or returns default value if conversion is not possible.
//
//	see convert.IntegerConverter.toIntegerWithDefault
//	Parameters:
//		- key string a key of element to get.
//		- defaultValue int the default value
//	Returns: integer value of the element or default value if conversion is not supported.
func (c *StringValueMap) GetAsIntegerWithDefault(key string, defaultValue int) int {
	if value, ok := c.Get(key); ok {
		return convert.IntegerConverter.ToIntegerWithDefault(value, defaultValue)
	}
	return defaultValue
}

// GetAsNullableUInteger converts map element into an unsigned integer or returns null if conversion is not possible.
//
//	see IntegerConverter.ToNullableInteger
//	Parameters: key string a key of element to get.
//	Returns: unsigned integer value of the element or o and false if conversion is not supported.
func (c *StringValueMap) GetAsNullableUInteger(key string) (uint, bool) {
	if value, ok := c.Get(key); ok {
		return convert.IntegerConverter.ToNullableUInteger(value)
	}
	return 0, false
}

// GetAsUInteger converts map element into an unsigned integer or returns 0 if conversion is not possible.
//
//	see GetAsIntegerWithDefault
//	Parameters: key string a key of element to get.
//	Returns uint
func (c *StringValueMap) GetAsUInteger(key string) uint {
	return c.GetAsUIntegerWithDefault(key, 0)
}

// GetAsUIntegerWithDefault converts map element into an unsigned integer or
// returns default value if conversion is not possible.
//
//	see IntegerConverter.ToIntegerWithDefault
//	Parameters:
//		- key string a key of element to get.
//		- defaultValue uint the default value
//	Returns uint integer value of the element or default value if conversion is not supported.
func (c *StringValueMap) GetAsUIntegerWithDefault(key string, defaultValue uint) uint {
	if value, ok := c.Get(key); ok {
		return convert.IntegerConverter.ToUIntegerWithDefault(value, defaultValue)
	}
	return defaultValue
}

// GetAsNullableLong converts map element into a int64 or returns null if conversion is not possible.
//
//	see convert.LongConverter.toNullableLong
//	Parameters: key string a key of element to get.
//	Returns: int64 value of the element or null if conversion is not supported.
func (c *StringValueMap) GetAsNullableLong(key string) (int64, bool) {
	if value, ok := c.Get(key); ok {
		return convert.LongConverter.ToNullableLong(value)
	}
	return 0, false
}

// GetAsLong converts map element into a int64 or returns 0 if conversion is not possible.
//
//	see GetAsLongWithDefault
//	Parameters: key string a key of element to get.
//	Returns: int64 value of the element or 0 if conversion is not supported.
func (c *StringValueMap) GetAsLong(key string) int64 {
	return c.GetAsLongWithDefault(key, 0)
}

// GetAsLongWithDefault converts map element into a int64 or returns default value if conversion is not possible.
//
//	see convert.LongConverter.toLongWithDefault
//	Parameters:
//		- key string a key of element to get.
//		- defaultValue int64 the default value
//	Returns: value of the element or default value if conversion is not supported.
func (c *StringValueMap) GetAsLongWithDefault(key string, defaultValue int64) int64 {
	if value, ok := c.Get(key); ok {
		return convert.LongConverter.ToLongWithDefault(value, defaultValue)
	}
	return defaultValue
}

// GetAsNullableULong converts map element into a uint64 or returns null if conversion is not possible.
// see convert.LongConverter.ToNullableLong
//
//	Parameters: key string a key of element to get.
//	Returns: int64 value of the element or false if conversion is not supported.
func (c *StringValueMap) GetAsNullableULong(key string) (uint64, bool) {
	if value, ok := c.Get(key); ok {
		return convert.LongConverter.ToNullableULong(value)
	}
	return 0, false
}

// GetAsULong converts map element into a uint64 or returns 0 if conversion is not possible.
//
//	see GetAsULongWithDefault
//	Parameters: key string a key of element to get.
//	Returns: uint64 value of the element or 0 if conversion is not supported.
func (c *StringValueMap) GetAsULong(key string) uint64 {
	return c.GetAsULongWithDefault(key, 0)
}

// GetAsULongWithDefault converts map element into a uint64 or returns default value if conversion is not possible.
//
//	see convert.LongConverter.ToULongWithDefault
//	Parameters:
//		- key string a key of element to get.
//		- defaultValue uint64 the default value
//	Returns: uint64 value of the element or default value if conversion is not supported.
func (c *StringValueMap) GetAsULongWithDefault(key string, defaultValue uint64) uint64 {
	if value, ok := c.Get(key); ok {
		return convert.LongConverter.ToULongWithDefault(value, defaultValue)
	}
	return defaultValue
}

// GetAsNullableFloat converts map element into a float32 or returns null if conversion is not possible.
//
//	see convert.FloatConverter.toNullableFloat
//	Parameters: key string a key of element to get.
//	Returns: float32 value of the element or null if conversion is not supported.
func (c *StringValueMap) GetAsNullableFloat(key string) (float32, bool) {
	if value, ok := c.Get(key); ok {
		return convert.FloatConverter.ToNullableFloat(value)
	}
	return 0, false
}

// GetAsFloat converts map element into a float32 or returns 0 if conversion is not possible.
//
//	see GetAsFloatWithDefault
//	Parameters: key string a key of element to get.
//	Returns: float32 value of the element or 0 if conversion is not supported.
func (c *StringValueMap) GetAsFloat(key string) float32 {
	return c.GetAsFloatWithDefault(key, 0)
}

// GetAsFloatWithDefault converts map element into a float32 or returns default value if conversion is not possible.
//
//	see convert.FloatConverter.toFloatWithDefault
//	Parameters:
//		- key string a key of element to get.
//		- defaultValue: float32 the default value
//	Returns: float32 value of the element or default value if conversion is not supported.
func (c *StringValueMap) GetAsFloatWithDefault(key string, defaultValue float32) float32 {
	if value, ok := c.Get(key); ok {
		return convert.FloatConverter.ToFloatWithDefault(value, defaultValue)
	}
	return defaultValue
}

// GetAsNullableDouble converts map element into a float64 or returns null if conversion is not possible.
//
//	see convert.DoubleConverter.toNullableDouble
//	Parameters: key string a key of element to get.
//	Returns: float64 value of the element or null if conversion is not supported.
func (c *StringValueMap) GetAsNullableDouble(key string) (float64, bool) {
	if value, ok := c.Get(key); ok {
		return convert.DoubleConverter.ToNullableDouble(value)
	}
	return 0, false
}

// GetAsDouble converts map element into a float64 or returns 0 if conversion is not possible.
//
//	see GetAsDoubleWithDefault
//	Parameters: key string a key of element to get.
//	Returns: value of the element or 0 if conversion is not supported.
func (c *StringValueMap) GetAsDouble(key string) float64 {
	return c.GetAsDoubleWithDefault(key, 0)
}

// GetAsDoubleWithDefault converts map element into a float64 or returns default value if conversion is not possible.
//
//	see convert.DoubleConverter.toDoubleWithDefault
//	Parameters:
//		- key string a key of element to get.
//		- defaultValue float64 the default value
//	Returns: float64 value of the element or default value if conversion is not supported.
func (c *StringValueMap) GetAsDoubleWithDefault(key string, defaultValue float64) float64 {
	if value, ok := c.Get(key); ok {
		return convert.DoubleConverter.ToDoubleWithDefault(value, defaultValue)
	}
	return defaultValue
}

// GetAsNullableDateTime converts map element into a time.Time or returns null if conversion is not possible.
//
//	see convert.DateTimeConverter.toNullableDateTime
//	Parameters: key string a key of element to get.
//	Returns: time.Time value of the element or null if conversion is not supported.
func (c *StringValueMap) GetAsNullableDateTime(key string) (time.Time, bool) {
	if value, ok := c.Get(key); ok {
		return convert.DateTimeConverter.ToNullableDateTime(value)
	}
	return time.Time{}, false
}

// GetAsDateTime converts map element into a time.Time or returns the current date if conversion is not possible.
//
//	see GetAsDateTimeWithDefault
//	Parameters: key string a key of element to get.
//	Returns: time.Time value of the element or the current date if conversion is not supported.
func (c *StringValueMap) GetAsDateTime(key string) time.Time {
	return c.GetAsDateTimeWithDefault(key, time.Time{})
}

// GetAsDateTimeWithDefault converts map element into a time.Time or returns default value if conversion is not possible.
//
//	see convert.DateTimeConverter.toDateTimeWithDefault
//	Parameters:
//		- key string a key of element to get.
//		- defaultValue time.Time the default value
//	Returns: time.Time value of the element or default value if conversion is not supported.
func (c *StringValueMap) GetAsDateTimeWithDefault(key string, defaultValue time.Time) time.Time {
	if value, ok := c.Get(key); ok {
		return convert.DateTimeConverter.ToDateTimeWithDefault(value, defaultValue)
	}
	return defaultValue
}

// GetAsValue converts map element into an AnyValue or returns an empty AnyValue if conversion is not possible.
//
//	see AnyValue
//	Parameters: key string a key of element to get.
//	Returns: *AnyValue value of the element or empty AnyValue if conversion is not supported.
func (c *StringValueMap) GetAsValue(key string) *AnyValue {
	if value, ok := c.Get(key); ok {
		return NewAnyValue(value)
	}
	return NewEmptyAnyValue()
}

// GetAsNullableArray converts map element into an AnyValueArray or returns null if conversion is not possible.
//
//	see AnyValueArray
//	see NewAnyValueArrayFromValue
//	Parameters: key string a key of element to get.
//	Returns: *AnyValueArray value of the element or nil if conversion is not supported.
func (c *StringValueMap) GetAsNullableArray(key string) (*AnyValueArray, bool) {
	if value, ok := c.Get(key); ok {
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
func (c *StringValueMap) GetAsArray(key string) *AnyValueArray {
	if value, ok := c.Get(key); ok {
		return NewAnyValueArrayFromValue(value)
	}
	return NewEmptyAnyValueArray()
}

// GetAsArrayWithDefault converts map element into an AnyValueArray or returns default value if conversion is not possible.
//
//	see AnyValueArray
//	see GetAsNullableArray
//	Parameters
//		- key string a key of element to get.
//		- defaultValue *AnyValueArray the default value
//	Returns: *AnyValueArray value of the element or default value if conversion is not supported.
func (c *StringValueMap) GetAsArrayWithDefault(key string, defaultValue *AnyValueArray) *AnyValueArray {
	if result, ok := c.GetAsNullableArray(key); ok {
		return result
	}
	return defaultValue
}

// GetAsNullableMap converts map element into an AnyValueMap or returns null if conversion is not possible.
//
//	see NewAnyValueMapFromValue
//	Parameters: key string a key of element to get.
//	Returns: *AnyValueMap value of the element or nil if conversion is not supported.
func (c *StringValueMap) GetAsNullableMap(key string) (*AnyValueMap, bool) {
	if value, ok := c.Get(key); ok {
		return NewAnyValueMapFromValue(value), true
	}
	return nil, false
}

// GetAsMap converts map element into an AnyValueMap or returns empty AnyValueMap if conversion is not possible.
//
//	see NewAnyValueMapFromValue
//	Parameters: key string a key of element to get.
//	Returns: *AnyValueMap value of the element or empty AnyValueMap if conversion is not supported.
func (c *StringValueMap) GetAsMap(key string) *AnyValueMap {
	if value, ok := c.Get(key); ok {
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
func (c *StringValueMap) GetAsMapWithDefault(key string, defaultValue *AnyValueMap) *AnyValueMap {
	if result, ok := c.GetAsNullableMap(key); ok {
		return result
	}
	return defaultValue
}

// String gets a string representation of the object. The result is a semicolon-separated
// list of key-value pairs as "key1=value1;key2=value2;key=value3"
//
//	Returns: a string representation of the object.
func (c *StringValueMap) String() string {
	builder := strings.Builder{}

	// Todo:: User encoder
	for key := range c.Value() {
		value, _ := c.Get(key)

		if builder.Len() > 0 {
			builder.WriteString(";")
		}

		if value != nil {
			builder.WriteString(fmt.Sprintf("%s=%s", key, value))
		} else {
			builder.WriteString(key)
		}
	}

	return builder.String()
}

// Clone creates a binary clone of this object.
//
//	Returns any a clone of this object.
func (c *StringValueMap) Clone() *StringValueMap {
	return NewStringValueMap(c._value)
}

func (c *StringValueMap) MarshalJSON() ([]byte, error) {
	bts, err := convert.JsonConverter.ToJson(c._value)
	return []byte(bts), err
}

func (c *StringValueMap) UnmarshalJSON(data []byte) error {
	values, err := convert.JsonConverter.FromJson(string(data))
	if err != nil {
		return err
	}
	c.Clear()
	val, _ := values.(map[string]any)
	c.AppendAny(val)
	return nil
}
