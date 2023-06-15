package data

import (
	"strings"
	"time"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
)

// AnyValueArray Cross-language implementation of dynamic object array what can hold values of any type. The stored values can be converted to different types using variety of accessor methods.
//
//	Example:
//		value1 := NewAnyValueArray([1, "123.456", "2018-01-01"]);
//		value1.GetAsBoolean(0);   // Result: true
//		value1.GetAsInteger(1);   // Result: 123
//		value1.GetAsFloat(1);     // Result: 123.456
//		value1.GetAsDateTime(2);  // Result: new Date(2018,0,1)
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
type AnyValueArray struct {
	_value []any
}

// NewEmptyAnyValueArray creates a new instance of the empty array.
//
//	Returns *AnyValueArray
func NewEmptyAnyValueArray() *AnyValueArray {
	return &AnyValueArray{
		_value: make([]any, 0, 10),
	}
}

// NewAnyValueArray creates a new instance of the array and assigns its value.
//
//	Parameters:
//		values []any
//	Returns: *AnyValueArray
func NewAnyValueArray(values []any) *AnyValueArray {
	c := &AnyValueArray{
		_value: make([]any, len(values)),
	}
	copy(c._value, values)
	return c
}

// NewAnyValueArrayFromValues creates a new AnyValueArray from a list of values
//
//	Parameters: values ...values any a list of values to initialize the created AnyValueArray
//	Returns: *AnyValueArray a newly created AnyValueArray.
func NewAnyValueArrayFromValues(values ...any) *AnyValueArray {
	return NewAnyValueArray(values)
}

// NewAnyValueArrayFromValue converts specified value into AnyValueArray.
//
//	see convertor.ArrayConverter.ToArray
//	Parameters: value any value to be converted
//	Returns: *AnyValueArray a newly created AnyValueArray.
func NewAnyValueArrayFromValue(value any) *AnyValueArray {
	return NewAnyValueArray(convert.ArrayConverter.ToArray(value))
}

// NewAnyValueArrayFromString splits specified string into elements using a separator and assigns the elements to a newly created AnyValueArray.
//
//	Parameters:
//		- values string a string value to be split and assigned to AnyValueArray
//		separator string a separator to split the string
//		- removeDuplicates bool	true to remove duplicated elements
//	Returns *AnyValueArray a newly created AnyValueArray.
func NewAnyValueArrayFromString(values string, separator string, removeDuplicates bool) *AnyValueArray {
	result := NewEmptyAnyValueArray()

	if values == "" {
		return result
	}

	items := strings.Split(values, separator)
	ln := len(items)

	if removeDuplicates {
		buffer := make(map[string]uint8, ln)
		for index := 0; index < ln; index++ {
			item := items[index]
			_, contains := buffer[item]
			if item != "" && !contains {
				result.Push(item)
				buffer[item] = 1
			}
		}
	} else {
		for index := 0; index < ln; index++ {
			item := items[index]
			if item != "" {
				result.Push(item)
			}
		}
	}

	return result
}

// InnerValue return inner value of array as any
func (c *AnyValueArray) InnerValue() any {
	return c._value
}

// Value returns array of elements []any
func (c *AnyValueArray) Value() []any {
	return c._value
}

// Len returns length of array
func (c *AnyValueArray) Len() int {
	return len(c._value)
}

// IsValidIndex checks that 0 <= index < len.
//
//	Parameters:
//		index int an index of the element to get.
//
// Returns: bool
func (c *AnyValueArray) IsValidIndex(index int) bool {
	return index >= 0 && index < c.Len()
}

// Get an array element specified by its index.
//
//	Parameters:
//		index int an index of the element to get.
//	Returns: value and true or nil and false if index is not valid
func (c *AnyValueArray) Get(index int) (any, bool) {
	if c.IsValidIndex(index) {
		return c._value[index], true
	}
	return nil, false
}

// Put a new value into array element specified by its index.
//
//	Parameters:
//		index int an index of the element to put.
//		value any a new value for array element.
//	Returns: true or false if index is invalid
func (c *AnyValueArray) Put(index int, value any) bool {
	if index <= 0 && index <= c.Len() {
		after := c._value[index:]
		before := c._value[:index]
		c._value = append(make([]any, 0, len(c._value)+1), before...)
		c._value = append(c._value, value)
		c._value = append(c._value, after...)
		return true
	}
	return false
}

// Remove an array element specified by its index
//
//	Parameters:
//		index int an index of the element to remove.
func (c *AnyValueArray) Remove(index int) bool {
	if c.IsValidIndex(index) {
		c._value = append(c._value[:index], c._value[index+1:]...)
		return true
	}
	return false
}

// Push element in the end of array
//
//	Parameters:
//		value any an value what need to insert
func (c *AnyValueArray) Push(value any) {
	c._value = append(c._value, value)
}

// Append new elements to this array.
//
//	Parameters:
//		elements []any a list of elements to be added.
func (c *AnyValueArray) Append(elements []any) {
	if elements != nil {
		c._value = append(c._value, elements...)
	}
}

// Clear this array by removing all its elements.
func (c *AnyValueArray) Clear() {
	c._value = make([]any, 0, 10)
}

// GetAsSingleObject inflate AnyValueArray as single object
//
//	Returns: any
func (c *AnyValueArray) GetAsSingleObject() any {
	return *c
}

// SetAsSingleObject sets AnyValueArray from input object
//
//	Parameters:
//		value any input object
func (c *AnyValueArray) SetAsSingleObject(value any) {
	c._value = convert.ArrayConverter.ToArray(value)
}

// GetAsObject gets the value stored in array element without any conversions.
// When element index is not defined it returns the entire array value.
//
//	Parameters
//		indexint an index of the element to get
//	Returns: any the element value or value of the array when index is not defined.
func (c *AnyValueArray) GetAsObject(index int) (any, bool) {
	return c.Get(index)
}

// SetAsObject sets a new value to array element specified by its index.
// When the index is not defined, it resets the entire array value.
//
//	see convert.ArrayConverter.ToArray
//	Parameters:
//		index int an index of the element to set
//		value any a new element or array value.
func (c *AnyValueArray) SetAsObject(index int, value any) bool {
	return c.Put(index, value)
}

// GetAsNullableString converts array element into a string or returns nil if conversion is not possible.
//
//	see convert.StringConverter.ToNullableString
//	Parameters:
//		index int an index of element to get.
//	Returns: string value of the element and true or "" and false if conversion is not supported or index is invalid.
func (c *AnyValueArray) GetAsNullableString(index int) (string, bool) {
	if value, ok := c.Get(index); ok && value != "" {
		return convert.StringConverter.ToNullableString(value)
	}
	return "", false
}

// GetAsString converts array element into a string or returns "" if conversion is not possible.
//
//	see GetAsStringWithDefault
//	Parameters:
//		index int an index of element to get.
//	Returns: string value ot the element or "" if conversion is not supported.
func (c *AnyValueArray) GetAsString(index int) string {
	return c.GetAsStringWithDefault(index, "")
}

// GetAsStringWithDefault converts array element into a string or returns default value if conversion is not possible.
// see convert.StringConverter.ToStringWithDefault
//
//	Parameters:
//		index int an index of element to get.
//		defaultValue: string the default value
//	Returns: string value ot the element or default value if conversion is not supported.
func (c *AnyValueArray) GetAsStringWithDefault(index int, defaultValue string) string {
	if value, ok := c.Get(index); ok && value != "" {
		return convert.StringConverter.ToStringWithDefault(value, defaultValue)
	}
	return defaultValue
}

// GetAsNullableBoolean converts array element into a boolean or returns nil if conversion is not possible.
// see convert.BooleanConverter.toNullableBoolean
//
//	Parameters:
//		index int an index of element to get.
//	Returns: boolean value of the element and true or false and false if conversion is not supported or index is invalid.
func (c *AnyValueArray) GetAsNullableBoolean(index int) (bool, bool) {
	if value, ok := c.Get(index); ok {
		return convert.BooleanConverter.ToNullableBoolean(value)
	}
	return false, false
}

// GetAsBoolean converts array element into a boolean or returns false if conversion is not possible.
//
//	see GetAsBooleanWithDefault
//	Parameters:
//		index int an index of element to get.
//	Returns: boolean value ot the element or false if conversion is not supported.
func (c *AnyValueArray) GetAsBoolean(index int) bool {
	return c.GetAsBooleanWithDefault(index, false)
}

// GetAsBooleanWithDefault converts array element into a boolean or returns default value if conversion is not possible.
//
//	see convert.BooleanConverter.toBooleanWithDefault
//	Parameters:
//		index int an index of element to get.
//		defaultValue: boolean the default value
//	Returns: boolean value ot the element or default value if conversion is not supported.
func (c *AnyValueArray) GetAsBooleanWithDefault(index int, defaultValue bool) bool {
	if value, ok := c.Get(index); ok {
		return convert.BooleanConverter.ToBooleanWithDefault(value, defaultValue)
	}
	return defaultValue
}

// GetAsNullableInteger converts array element into an integer or returns nil if conversion is not possible.
//
//	see convert.IntegerConverter.toNullableInteger
//	Parameters:
//		index int an index of element to get.
//	Returns: integer value of the element and true or 0 and false if conversion is not supported or index is invalid.
func (c *AnyValueArray) GetAsNullableInteger(index int) (int, bool) {
	if value, ok := c.Get(index); ok {
		return convert.IntegerConverter.ToNullableInteger(value)
	}
	return 0, false
}

// GetAsInteger converts array element into an integer or returns 0 if conversion is not possible.
//
//	see GetAsIntegerWithDefault
//	Parameters:
//		index int an index of element to get.
//	Returns: integer value ot the element or 0 if conversion is not supported.
func (c *AnyValueArray) GetAsInteger(index int) int {
	return c.GetAsIntegerWithDefault(index, 0)
}

// GetAsIntegerWithDefault converts array element into an integer or returns default value if conversion is not possible.
//
//	see convert.IntegerConverter.toIntegerWithDefault
//	Parameters:
//		index int an index of element to get.
//		defaultValue int the default value
//	Returns: integer value ot the element or default value if conversion is not supported.
func (c *AnyValueArray) GetAsIntegerWithDefault(index int, defaultValue int) int {
	if value, ok := c.Get(index); ok {
		return convert.IntegerConverter.ToIntegerWithDefault(value, defaultValue)
	}
	return defaultValue
}

// GetAsNullableUInteger converts array element into an unsigned integer or returns nil if conversion is not possible.
//
//	see convert.IntegerConverter.ToNullableUInteger
//	Parameters:
//		index int an index of element to get.
//	Returns: unsigned integer value of the element and true or
//		0 and false if conversion is not supported or index is invalid.
func (c *AnyValueArray) GetAsNullableUInteger(index int) (uint, bool) {
	if value, ok := c.Get(index); ok {
		return convert.IntegerConverter.ToNullableUInteger(value)
	}
	return 0, false
}

// GetAsUInteger converts array element into an unsigned integer or returns 0 if conversion is not possible.
//
//	see GetAsIntegerWithDefault
//	Parameters: index int an index of element to get.
//	Returns: uint unsigned integer value ot the element or 0 if conversion is not supported.
func (c *AnyValueArray) GetAsUInteger(index int) uint {
	return c.GetAsUIntegerWithDefault(index, 0)
}

// GetAsUIntegerWithDefault converts array element into an integer or
// returns default value if conversion is not possible.
// see IntegerConverter.ToIntegerWithDefault
//
//	Parameters:
//		- index int an index of element to get.
//		- defaultValue uint the default value
//	Returns: uint unsigned integer value ot the element or default value if conversion is not supported.
func (c *AnyValueArray) GetAsUIntegerWithDefault(index int, defaultValue uint) uint {
	if value, ok := c.Get(index); ok {
		return convert.IntegerConverter.ToUIntegerWithDefault(value, defaultValue)
	}
	return defaultValue
}

// GetAsNullableLong converts array element into a long or returns nil if conversion is not possible.
//
//	see convert.LongConverter.ToNullableLong
//	Parameters:
//		index int an index of element to get.
//	Returns: int64 value of the element and true or 0 and false if conversion is not supported or index is invalid.
func (c *AnyValueArray) GetAsNullableLong(index int) (int64, bool) {
	if value, ok := c.Get(index); ok {
		return convert.LongConverter.ToNullableLong(value)
	}
	return 0, false
}

// GetAsLong converts array element into a long or returns 0 if conversion is not possible.
//
//	see GetAsLongWithDefault
//	Parameters: index int an index of element to get.
//	Returns: int64 value ot the element or 0 if conversion is not supported.
func (c *AnyValueArray) GetAsLong(index int) int64 {
	return c.GetAsLongWithDefault(index, 0)
}

// GetAsLongWithDefault converts array element into a long or returns default value if conversion is not possible.
//
//	see convert.LongConverter.ToLongWithDefault
//	Parameters:
//		- index int an index of element to get.
//		- defaultValue int64 the default value.
//	Returns: int64 value ot the element or default value if conversion is not supported.
func (c *AnyValueArray) GetAsLongWithDefault(index int, defaultValue int64) int64 {
	if value, ok := c.Get(index); ok {
		return convert.LongConverter.ToLongWithDefault(value, defaultValue)
	}
	return defaultValue
}

// GetAsNullableULong converts array element into a unsigned long or returns nil if conversion is not possible.
//
//	see convert.LongConverter.ToNullableULong
//	Parameters:
//		index int an index of element to get.
//	Returns: uint64 value of the element and true or 0 and false if conversion is not supported or index is invalid.
func (c *AnyValueArray) GetAsNullableULong(index int) (uint64, bool) {
	if value, ok := c.Get(index); ok {
		return convert.LongConverter.ToNullableULong(value)
	}
	return 0, false
}

// GetAsULong converts array element into a unsigned long or returns 0 if conversion is not possible.
// see GetAsLongWithDefault
//
//	Parameters: index int an index of element to get.
//	Returns: uint64 value ot the element or 0 if conversion is not supported.
func (c *AnyValueArray) GetAsULong(index int) uint64 {
	return c.GetAsULongWithDefault(index, 0)
}

// GetAsULongWithDefault converts array element into a unsigned long or
// returns default value if conversion is not possible.
//
//	see convert.LongConverter.ToLongWithDefault
//	Parameters:
//		- index int an index of element to get.
//		- defaultValue int64 the default value
//	Returns: uint64 value ot the element or default value if conversion is not supported.
func (c *AnyValueArray) GetAsULongWithDefault(index int, defaultValue uint64) uint64 {
	if value, ok := c.Get(index); ok {
		return convert.LongConverter.ToULongWithDefault(value, defaultValue)
	}
	return defaultValue
}

// GetAsNullableFloat converts array element into a float or returns nil if conversion is not possible.
//
//	see convert.FloatConverter.ToNullableFloat
//	Parameters: index int an index of element to get.
//	Returns: float64 value of the element and true or 0 and false if conversion is not supported or index is invalid.
func (c *AnyValueArray) GetAsNullableFloat(index int) (float32, bool) {
	if value, ok := c.Get(index); ok {
		return convert.FloatConverter.ToNullableFloat(value)
	}
	return 0, false
}

// GetAsFloat converts array element into a float or returns 0 if conversion is not possible.
//
//	see GetAsFloatWithDefault
//	Parameters: index int an index of element to get.
//	Returns: float value ot the element or 0 if conversion is not supported.
func (c *AnyValueArray) GetAsFloat(index int) float32 {
	return c.GetAsFloatWithDefault(index, 0)
}

// GetAsFloatWithDefault converts array element into a float or returns default value if conversion is not possible.
//
//	see convert.FloatConverter.toFloatWithDefault
//	Parameters:
//		- index int an index of element to get.
//		- defaultValue: number the default value
//	Returns: number float value ot the element or default value if conversion is not supported.
func (c *AnyValueArray) GetAsFloatWithDefault(index int, defaultValue float32) float32 {
	if value, ok := c.Get(index); ok {
		return convert.FloatConverter.ToFloatWithDefault(value, defaultValue)
	}
	return defaultValue
}

// GetAsNullableDouble converts array element into a double or returns nil if conversion is not possible.
//
//	see convert.DoubleConverter.toNullableDouble
//	Parameters: index int an index of element to get.
//	Returns: float64 value of the element and true or 0 and false if conversion is not supported or index is invalid.
func (c *AnyValueArray) GetAsNullableDouble(index int) (float64, bool) {
	if value, ok := c.Get(index); ok {
		return convert.DoubleConverter.ToNullableDouble(value)
	}
	return 0, false
}

// GetAsDouble converts array element into a double or returns 0 if conversion is not possible.
//
//	see GetAsDoubleWithDefault
//	Parameters: index int an index of element to get.
//	Returns: double value ot the element or 0 if conversion is not supported.
func (c *AnyValueArray) GetAsDouble(index int) float64 {
	return c.GetAsDoubleWithDefault(index, 0)
}

// GetAsDoubleWithDefault converts array element into a double or returns default value if conversion is not possible.
//
//	see convert.DoubleConverter.ToDoubleWithDefault
//	Parameters:
//		- index int an index of element to get.
//		- defaultValue: float64 the default value.
//	Returns: double value ot the element or default value if conversion is not supported.
func (c *AnyValueArray) GetAsDoubleWithDefault(index int, defaultValue float64) float64 {
	if value, ok := c.Get(index); ok {
		return convert.DoubleConverter.ToDoubleWithDefault(value, defaultValue)
	}
	return defaultValue
}

// GetAsNullableDateTime converts array element into a time.Time or returns nil if conversion is not possible.
//
//	see convert.DateTimeConverter.ToNullableDateTime
//	Parameters: index int an index of element to get.
//	Returns: time.Time value of the element and true or zero time and false if conversion is not supported or index is invalid.
func (c *AnyValueArray) GetAsNullableDateTime(index int) (time.Time, bool) {
	if value, ok := c.Get(index); ok {
		return convert.DateTimeConverter.ToNullableDateTime(value)
	}
	return time.Time{}, false
}

// GetAsDateTime converts array element into a time.Time or returns the current date if conversion is not possible.
//
//	see GetAsDateTimeWithDefault
//	Parameters: index int an index of element to get.
//	Returns: time.Time value ot the element or the current date if conversion is not supported.
func (c *AnyValueArray) GetAsDateTime(index int) time.Time {
	return c.GetAsDateTimeWithDefault(index, time.Time{})
}

// GetAsDateTimeWithDefault converts array element into a time.Time or returns default value if conversion is not possible.
//
//	see covert.DateTimeConverter.toDateTimeWithDefault
//	Parameters:
//		- index int an index of element to get.
//		- defaultValue: time.Time the default value.
//	Returns: time.time value ot the element or default value if conversion is not supported.
func (c *AnyValueArray) GetAsDateTimeWithDefault(index int, defaultValue time.Time) time.Time {
	if value, ok := c.Get(index); ok {
		return convert.DateTimeConverter.ToDateTimeWithDefault(value, defaultValue)
	}
	return defaultValue
}

// GetAsNullableDuration converts array element into a time.Duration or returns nil if conversion is not possible.
//
//	see convert.DateTimeConverter.toNullableDateTime
//	Parameters: index int an index of element to get.
//	Returns: time.Duration value of the element and true or 0 and false if conversion is not supported or index is ivalid.
func (c *AnyValueArray) GetAsNullableDuration(index int) (time.Duration, bool) {
	if value, ok := c.Get(index); ok {
		return convert.DurationConverter.ToNullableDuration(value)
	}
	return 0, false
}

// GetAsDuration converts array element into a time.Duration or returns the current date if conversion is not possible.
//
//	see GetAsDurationWithDefault
//	Parameters: index int an index of element to get.
//	Returns: time.Duration value ot the element or the current date if conversion is not supported.
func (c *AnyValueArray) GetAsDuration(index int) time.Duration {
	return c.GetAsDurationWithDefault(index, 0*time.Millisecond)
}

// GetAsDurationWithDefault converts array element into a time.Duration or returns default value if conversion is not possible.
//
//	see convert.DateTimeConverter.toDateTimeWithDefault
//	Parameters:
//		- index int an index of element to get.
//		- defaultValue: time.Duration the default value
//	Returns: time.Duration value ot the element or default value if conversion is not supported.
func (c *AnyValueArray) GetAsDurationWithDefault(index int, defaultValue time.Duration) time.Duration {
	if value, ok := c.Get(index); ok {
		return convert.DurationConverter.ToDurationWithDefault(value, defaultValue)
	}
	return defaultValue
}

// GetAsNullableType converts array element into a value defined by specied typecode. If conversion is not possible it returns nil.
//
//	see convert.TypeConverter.ToNullableType
//	Parameters
//		- type: TypeCode the TypeCode that defined the type of the result
//		- index int an index of element to get.
//	Returns: any element value defined by the typecode and true or nil and false if conversion is not supported or index is invalid.
func (c *AnyValueArray) GetAsNullableType(typ convert.TypeCode, index int) (any, bool) {
	if value, ok := c.Get(index); ok {
		return convert.TypeConverter.ToNullableType(typ, value)
	}
	return nil, false
}

// GetAsType converts array element into a value defined by specied typecode.
// If conversion is not possible it returns default value for the specified type.
//
//	see GetAsTypeWithDefault
//	Parameters:
//		- type TypeCode the TypeCode that defined the type of the result
//		- index int an index of element to get.
//	Returns: any element value defined by the typecode or default if conversion is not supported.
func (c *AnyValueArray) GetAsType(typ convert.TypeCode, index int) any {
	return c.GetAsTypeWithDefault(typ, index, nil)
}

// GetAsTypeWithDefault converts array element into a value defined by specied typecode.
// If conversion is not possible it returns default value.
//
//	see convert.TypeConverter.ToTypeWithDefault
//	Parameters:
//		- type TypeCode the TypeCode that defined the type of the result
//		- index int an index of element to get.
//		- defaultValue any the default value
//	Returns: any element value defined by the typecode or default value if conversion is not supported.
func (c *AnyValueArray) GetAsTypeWithDefault(typ convert.TypeCode, index int, defaultValue any) any {
	if value, ok := c.Get(index); ok {
		return convert.TypeConverter.ToTypeWithDefault(typ, value, defaultValue)
	}
	return defaultValue
}

// GetAsValue converts array element into an AnyValue or returns an empty AnyValue if conversion is not possible.
//
//	see AnyValue
//	see AnyValue.constructor
//	Parameters: index int an index of element to get.
//	Returns: AnyValue value of the element and true or empty AnyValue and false if conversion is not supported or index is invalid.
func (c *AnyValueArray) GetAsValue(index int) (*AnyValue, bool) {
	if value, ok := c.Get(index); ok {
		return NewAnyValue(value), true
	}
	return NewEmptyAnyValue(), false
}

// GetAsNullableArray converts array element into an AnyValueArray or returns nil if conversion is not possible.
//
//	see NewAnyValueArrayFromValue
//	Parameters: index int an index of element to get.
//	Returns: *AnyValueArray value of the element and true or nil and false if conversion is not supported or index is invalid.
func (c *AnyValueArray) GetAsNullableArray(index int) (*AnyValueArray, bool) {
	if value, ok := c.Get(index); ok {
		return NewAnyValueArrayFromValue(value), true
	}
	return nil, false
}

// GetAsArray converts array element into an AnyValueArray or returns empty AnyValueArray if conversion is not possible.
//
//	see NewAnyValueArrayFromValue
//	Parameters: index int an index of element to get.
//	Returns *AnyValueArray value of the element and true or empty AnyValueArray and false if conversion is not supported or index is ivalid.
func (c *AnyValueArray) GetAsArray(index int) *AnyValueArray {
	if value, ok := c.Get(index); ok {
		return NewAnyValueArrayFromValue(value)
	}
	return NewEmptyAnyValueArray()
}

// GetAsArrayWithDefault converts array element into an AnyValueArray or returns default value if conversion is not possible.
//
//	see GetAsNullableArray
//	Parameters:
//		- index int an index of element to get.
//		- defaultValue *AnyValueArray the default value
//	Returns: *AnyValueArray value of the element or default value if conversion is not supported.
func (c *AnyValueArray) GetAsArrayWithDefault(index int, defaultValue *AnyValueArray) *AnyValueArray {
	if result, ok := c.GetAsNullableArray(index); ok {
		return result
	}
	return defaultValue
}

// GetAsNullableMap converts array element into an AnyValueMap or returns nil if conversion is not possible.
//
//	see AnyValueMap
//	see NewAnyValueMapFromValue
//	Parameters: index int an index of element to get.
//	Returns: *AnyValueMap value of the element and true or nil and false if conversion is not supported or index is invalid.
func (c *AnyValueArray) GetAsNullableMap(index int) (*AnyValueMap, bool) {
	if value, ok := c.Get(index); ok {
		return NewAnyValueMapFromValue(value), true
	}
	return nil, false
}

// GetAsMap converts array element into an AnyValueMap or returns empty AnyValueMap if conversion is not possible.
//
//	see AnyValueMap
//	see NewAnyValueMapFromValue
//	Parameters: index int an index of element to get.
//	Returns: *AnyValueMap
func (c *AnyValueArray) GetAsMap(index int) *AnyValueMap {
	if value, ok := c.Get(index); ok {
		return NewAnyValueMapFromValue(value)
	}
	return NewEmptyAnyValueMap()
}

// GetAsMapWithDefault converts array element into an AnyValueMap or returns default value if conversion is not possible.
//
//	see GetAsNullableMap
//	Parameters
//		- index int an index of element to get.
//		- defaultValue *AnyValueMap the default value
//	Returns *AnyValueMap value of the element or default value if conversion is not supported.
func (c *AnyValueArray) GetAsMapWithDefault(index int, defaultValue *AnyValueMap) *AnyValueMap {
	if result, ok := c.GetAsNullableMap(index); ok {
		return NewAnyValueMapFromValue(result)
	}
	return defaultValue
}

// Contains checks if this array contains a value. The check uses direct comparison between elements and the specified value.
//
//	Parameters: value any a value to be checked
//	Returns: bool true if this array contains the value or false otherwise.
func (c *AnyValueArray) Contains(value any) bool {
	for index := 0; index < c.Len(); index++ {
		if value == c._value[index] {
			return true
		}
	}
	return false
}

// ContainsAsType checks if this array contains a value.
//
//	The check before comparison converts elements and the value to type specified by type code.
//	see convert.TypeConverter.ToType
//	see convert.TypeConverter.ToNullableType
//	Parameters:
//		- typeCode TypeCode a type code that defines a type to convert values before comparison
//		- value any a value to be checked
//	Returns: bool true if this array contains the value or false otherwise.
func (c *AnyValueArray) ContainsAsType(typ convert.TypeCode, value any) bool {
	typedValue := convert.TypeConverter.ToType(typ, value)

	for index := 0; index < c.Len(); index++ {
		thisTypedValue := convert.TypeConverter.ToType(typ, c._value[index])
		if typedValue == thisTypedValue {
			return true
		}
	}
	return false
}

// Clone creates a binary clone of this object.
//
//	Returns: any a clone of this object.
func (c *AnyValueArray) Clone() *AnyValueArray {
	return NewAnyValueArray(c._value)
}

// String converts array to string with coma separator
func (c *AnyValueArray) String() string {
	if c.Len() == 0 {
		return ""
	}
	builder := strings.Builder{}
	builder.WriteString(convert.StringConverter.ToStringWithDefault(c._value[0], ""))
	for index := 1; index < c.Len(); index++ {
		builder.WriteString(",")
		builder.WriteString(convert.StringConverter.ToStringWithDefault(c._value[index], ""))
	}
	return builder.String()
}
