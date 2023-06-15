package data

import (
	"time"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
)

// AnyValue Cross-language implementation of dynamic object what can hold value of any type.
// The stored value can be converted to different types using variety of accessor methods.
//
//	Example:
//		value := data.NewAnyValue("123.456")
//
//		value.GetAsInteger() // Result: 123
//		value.GetAsString()  // Result: "123.456"
//		value.GetAsFloat()   // Result: 123.456
type AnyValue struct {
	_value any
}

// NewEmptyAnyValue creates a new empty instance of the object
//
//	Returns: new empty object
func NewEmptyAnyValue() *AnyValue {
	return &AnyValue{_value: nil}
}

// NewAnyValue creates a new instance of the object and assigns its value.
//
//	Parameters: "value" - value to initialize this object.
//	Returns: new object.
func NewAnyValue(value any) *AnyValue {
	v, ok := value.(*AnyValue)
	if ok {
		return v
	} else {
		return &AnyValue{_value: value}
	}
}

// InnerValue gets the value stored in this object without any conversions.
//
//	Returns: the object value.
func (c *AnyValue) InnerValue() any {
	return c._value
}

// Value gets the value stored in this object without any conversions.
//
//	Returns: the object value.
func (c *AnyValue) Value() any {
	return c._value
}

// TypeCode gets type code for the value stored in this object.
//
//	Returns: type code of the object value.
func (c *AnyValue) TypeCode() convert.TypeCode {
	return convert.TypeConverter.ToTypeCode(c._value)
}

// GetAsObject gets the value stored in this object without any conversions.
//
//	Returns: the object value.
func (c *AnyValue) GetAsObject() any {
	return c._value
}

// SetAsObject sets a new value for this object.
//
//	Parameters: "value" - the new object value.
func (c *AnyValue) SetAsObject(value any) {
	c._value = value
}

// GetAsNullableString converts object value into a string or returns null if conversion is not possible.
//
//	Returns: string and true value or "" and false if conversion is not supported.
func (c *AnyValue) GetAsNullableString() (string, bool) {
	return convert.StringConverter.ToNullableString(c._value)
}

// GetAsString converts object value into a string or returns "" if conversion is not possible.
// Returns: string value or "" if conversion is not supported.
func (c *AnyValue) GetAsString() string {
	return c.GetAsStringWithDefault("")
}

// GetAsStringWithDefault converts object value into a string or returns default value if conversion is not possible.
//
//	Parameters: "defaultValue" - the default value.
//	Returns: string value or default if conversion is not supported.
func (c *AnyValue) GetAsStringWithDefault(defaultValue string) string {
	return convert.StringConverter.ToStringWithDefault(c._value, defaultValue)
}

// GetAsNullableBoolean converts object value into a boolean or returns null if conversion is not possible.
//
//	Returns: boolean value and true or false and false if conversion is not supported.
func (c *AnyValue) GetAsNullableBoolean() (bool, bool) {
	return convert.BooleanConverter.ToNullableBoolean(c._value)
}

// GetAsBoolean converts object value into a boolean or returns false if conversion is not possible.
//
//	Returns: string value or false if conversion is not supported.
func (c *AnyValue) GetAsBoolean() bool {
	return c.GetAsBooleanWithDefault(false)
}

// GetAsBooleanWithDefault converts object value into a boolean or returns default value if conversion is not possible.
//
//	Parameters: "defaultValue" - the default value.
//	Returns: boolean value or default if conversion is not supported.
func (c *AnyValue) GetAsBooleanWithDefault(defaultValue bool) bool {
	return convert.BooleanConverter.ToBooleanWithDefault(c._value, defaultValue)
}

// GetAsNullableInteger converts object value into an integer or returns null if conversion is not possible.
//
//	Returns: integer value and true or 0 and false if conversion is not supported.
func (c *AnyValue) GetAsNullableInteger() (int, bool) {
	return convert.IntegerConverter.ToNullableInteger(c._value)
}

// GetAsInteger converts object value into an integer or returns 0 if conversion is not possible.
//
//	Returns: integer value or 0 if conversion is not supported.
func (c *AnyValue) GetAsInteger() int {
	return c.GetAsIntegerWithDefault(0)
}

// GetAsIntegerWithDefault converts object value into a integer or returns default value if conversion is not possible.
//
//	Parameters: "defaultValue" - the default value
//	Returns: integer value or default if conversion is not supported.
func (c *AnyValue) GetAsIntegerWithDefault(defaultValue int) int {
	return convert.IntegerConverter.ToIntegerWithDefault(c._value, defaultValue)
}

// GetAsNullableLong converts object value into a long or returns null if conversion is not possible.
//
//	Returns: long value and true or 0 and false if conversion is not supported.
func (c *AnyValue) GetAsNullableLong() (int64, bool) {
	return convert.LongConverter.ToNullableLong(c._value)
}

// GetAsLong converts object value into a long or returns 0 if conversion is not possible.
//
//	Returns: string value or 0 if conversion is not supported.
func (c *AnyValue) GetAsLong() int64 {
	return c.GetAsLongWithDefault(0)
}

// GetAsLongWithDefault converts object value into a long or returns default value if conversion is not possible.
//
//	Parameters: "defaultValue" - the default value
//	Returns: long value or default if conversion is not supported.
func (c *AnyValue) GetAsLongWithDefault(defaultValue int64) int64 {
	return convert.LongConverter.ToLongWithDefault(c._value, defaultValue)
}

// GetAsNullableUInteger converts object value into an unsigned integer
// or returns default value if conversion is not possible.
//
//	Parameters: "defaultValue" - the default value
//	Returns: integer value or default if conversion is not supported.
func (c *AnyValue) GetAsNullableUInteger() (uint, bool) {
	return convert.IntegerConverter.ToNullableUInteger(c._value)
}

// GetAsUInteger converts object value into an unsigned integer or returns 0 if conversion is not possible.
//
//	Returns: unsigned integer value or 0 if conversion is not supported.
func (c *AnyValue) GetAsUInteger() uint {
	return c.GetAsUIntegerWithDefault(0)
}

// GetAsUIntegerWithDefault converts object value into a unsigned integer or returns
// default value if conversion is not possible.
//
//	Parameters: "defaultValue" - the default value
//	Returns: unsigned integer value or default if conversion is not supported.
func (c *AnyValue) GetAsUIntegerWithDefault(defaultValue uint) uint {
	return convert.IntegerConverter.ToUIntegerWithDefault(c._value, defaultValue)
}

// GetAsNullableULong converts object value into an unsigned long
// or returns default value if conversion is not possible.
//
//	Parameters: "defaultValue" - the default value
//	Returns: long value or default if conversion is not supported.
func (c *AnyValue) GetAsNullableULong() (uint64, bool) {
	return convert.LongConverter.ToNullableULong(c._value)
}

// GetAsULong converts object value into an unsigned long or returns 0 if conversion is not possible.
//
//	Returns: unsigned long value or 0 if conversion is not supported.
func (c *AnyValue) GetAsULong() uint64 {
	return c.GetAsULongWithDefault(0)
}

// GetAsULongWithDefault converts object value into a unsiged long or returns default value
// if conversion is not possible.
//
//	Parameters: "defaultValue" - the default value
//	Returns: unsigned long value or default if conversion is not supported.
func (c *AnyValue) GetAsULongWithDefault(defaultValue uint64) uint64 {
	return convert.LongConverter.ToULongWithDefault(c._value, defaultValue)
}

// GetAsNullableFloat converts object value into a float or returns null if conversion is not possible.
//
//	Returns: float value and true or 0.0 and false if conversion is not supported.
func (c *AnyValue) GetAsNullableFloat() (float32, bool) {
	return convert.FloatConverter.ToNullableFloat(c._value)
}

// GetAsFloat converts object value into a float or returns 0 if conversion is not possible.
//
//	Returns: float value or 0 if conversion is not supported.
func (c *AnyValue) GetAsFloat() float32 {
	return c.GetAsFloatWithDefault(0)
}

// GetAsFloatWithDefault converts object value into a float or returns default value if conversion is not possible.
//
//	Parameters: "defaultValue" - the default value
//	Returns: float value or default if conversion is not supported.
func (c *AnyValue) GetAsFloatWithDefault(defaultValue float32) float32 {
	return convert.FloatConverter.ToFloatWithDefault(c._value, defaultValue)
}

// GetAsNullableDouble converts object value into a double or returns null if conversion is not possible.
//
//	Returns: double value and true or 0.0 and false if conversion is not supported.
func (c *AnyValue) GetAsNullableDouble() (float64, bool) {
	return convert.DoubleConverter.ToNullableDouble(c._value)
}

// GetAsDouble converts object value into a double or returns 0 if conversion is not possible.
//
//	Returns: double value or 0 if conversion is not supported.
func (c *AnyValue) GetAsDouble() float64 {
	return c.GetAsDoubleWithDefault(0)
}

// GetAsDoubleWithDefault converts object value into a double or returns default value if conversion is not possible.
//
//	Parameters: "defaultValue" - the default value
//	Returns: double value or default if conversion is not supported.
func (c *AnyValue) GetAsDoubleWithDefault(defaultValue float64) float64 {
	return convert.DoubleConverter.ToDoubleWithDefault(c._value, defaultValue)
}

// GetAsNullableDateTime converts object value into a Date or returns null if conversion is not possible.
//
//	Returns: DateTime value and true or zero time and false if conversion is not supported.
func (c *AnyValue) GetAsNullableDateTime() (time.Time, bool) {
	return convert.DateTimeConverter.ToNullableDateTime(c._value)
}

// GetAsDateTime converts object value into a Date or returns current date if conversion is not possible.
//
//	Returns: DateTime value or current date if conversion is not supported.
func (c *AnyValue) GetAsDateTime() time.Time {
	return c.GetAsDateTimeWithDefault(time.Time{})
}

// GetAsDateTimeWithDefault converts object value into a Date or returns default value if conversion is not possible.
//
//	Parameters: "defaultValue" - the default value
//	Returns: DateTime value or default if conversion is not supported.
func (c *AnyValue) GetAsDateTimeWithDefault(defaultValue time.Time) time.Time {
	return convert.DateTimeConverter.ToDateTimeWithDefault(c._value, defaultValue)
}

// GetAsNullableDuration converts object value into a Duration or returns null if conversion is not possible.
//
//	Returns: Duration value and true or 0 and false if conversion is not supported.
func (c *AnyValue) GetAsNullableDuration() (time.Duration, bool) {
	return convert.DurationConverter.ToNullableDuration(c._value)
}

// GetAsDuration converts object value into a Duration or returns current date if conversion is not possible.
//
//	Returns: Duration value or current date if conversion is not supported.
func (c *AnyValue) GetAsDuration() time.Duration {
	return c.GetAsDurationWithDefault(0 * time.Millisecond)
}

// GetAsDurationWithDefault converts object value into a Duration or returns default value if conversion is not possible.
// Parameters: "defaultValue" - the default value
// Returns: Duration value or default if conversion is not supported.
func (c *AnyValue) GetAsDurationWithDefault(defaultValue time.Duration) time.Duration {
	return convert.DurationConverter.ToDurationWithDefault(c._value, defaultValue)
}

// GetAsNullableType converts object value into a value defined by specied typecode. If conversion is not possible it returns null.
//
//	Parameters: "typ" - the TypeCode that defined the type of the result.
//	Returns: value defined by the typecode and true or null and false if conversion is not supported.
func (c *AnyValue) GetAsNullableType(typ convert.TypeCode) (any, bool) {
	return convert.TypeConverter.ToNullableType(typ, c._value)
}

// GetAsType converts object value into a value defined by specied typecode. If conversion
// is not possible it returns default value for the specified type.
//
//	Parameters: "typ" - the TypeCode that defined the type of the result.
//	Returns: value defined by the typecode or type default value if conversion is not supported.
func (c *AnyValue) GetAsType(typ convert.TypeCode) any {
	return c.GetAsTypeWithDefault(typ, nil)
}

// GetAsTypeWithDefault converts object value into a value defined by specied typecode. If conversion
// is not possible it returns default value.
//
//		Parameters:
//	 	"typ" - the TypeCode that defined the type of the result;
//	 	"defaultValue" - the default value.
//		Returns: value defined by the typecode or type default value if conversion is not supported.
func (c *AnyValue) GetAsTypeWithDefault(typ convert.TypeCode, defaultValue any) any {
	return convert.TypeConverter.ToTypeWithDefault(typ, c._value, defaultValue)
}

// GetAsArray converts object value into an AnyArray or returns empty AnyArray if conversion is not possible.
//
//	Returns: AnyArray value or empty AnyArray if conversion is not supported.
func (c *AnyValue) GetAsArray() *AnyValueArray {
	return NewAnyValueArrayFromValue(c._value)
}

// GetAsMap converts object value into AnyMap or returns empty AnyMap if conversion is not possible.
//
//	Returns: AnyMap value or empty AnyMap if conversion is not supported.
func (c *AnyValue) GetAsMap() *AnyValueMap {
	return NewAnyValueMapFromValue(c._value)
}

// Equals compares this object value to specified specified value. When direct
// comparison gives negative results it tries to compare values as strings.
//
//	Parameters: "obj" - the value to be compared with.
//	Returns: true when objects are equal and false otherwise.
func (c *AnyValue) Equals(obj any) bool {
	if obj == nil && c._value == nil {
		return true
	}
	if obj == nil || c._value == nil {
		return false
	}

	v, ok := obj.(*AnyValue)
	if ok {
		obj = v._value
	}

	strThisValue, strThisValueOk := convert.StringConverter.ToNullableString(c._value)
	strValue, strValueOk := convert.StringConverter.ToNullableString(obj)

	if !strThisValueOk && !strValueOk {
		return true
	}
	if !strThisValueOk || !strValueOk {
		return false
	}
	return strThisValue == strValue
}

// EqualsAsType compares this object value to specified specified value. When direct
// comparison gives negative results it converts values to type specified by
// type code and compare them again.
//
//		Parameters:
//	 	"typ" - the TypeCode that defined the type of the result.
//	 	"obj" - the value to be compared with.
//		Returns: true when objects are equal and false otherwise.
func (c *AnyValue) EqualsAsType(typ convert.TypeCode, obj any) bool {
	if obj == nil && c._value == nil {
		return true
	}
	if obj == nil || c._value == nil {
		return false
	}

	v, ok := obj.(*AnyValue)
	if ok {
		obj = v._value
	}

	typedThisValue := convert.TypeConverter.ToType(typ, c._value)
	typedValue := convert.TypeConverter.ToType(typ, obj)

	return typedThisValue == typedValue
}

// Clone creates a binary clone of this object.
//
//	Returns: a clone of this object.
func (c *AnyValue) Clone() *AnyValue {
	return NewAnyValue(c._value)
}

// String gets a string representation of the object.
//
//	Returns: a string representation of the object.
func (c *AnyValue) String() string {
	return convert.StringConverter.ToString(c._value)
}
