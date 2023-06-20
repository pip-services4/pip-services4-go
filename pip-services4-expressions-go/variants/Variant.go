package variants

import (
	"time"

	cconv "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
)

// Variant defines container for variant values.
type Variant struct {
	typ   VariantType
	value any
}

// Empty variant constant
var Empty *Variant = EmptyVariant()

// EmptyVariant constructs an empty variant object
func EmptyVariant() *Variant {
	return &Variant{
		typ:   Null,
		value: nil,
	}
}

// NewVariant constructs this class and assignes a value.
//
//	Parameters:
//		- value: another variant value.
func NewVariant(value any) *Variant {
	c := &Variant{}
	c.SetAsObject(value)
	return c
}

// VariantFromInteger creates a new variant from Integer value.
//
//	Parameters:
//		- value: a variant value.
//	Returns: a created variant object
func VariantFromInteger(value int) *Variant {
	c := &Variant{}
	c.SetAsInteger(value)
	return c
}

// VariantFromLong сreates a new variant from Long value.
// Parameters:
//   - value: a variant value.
//
// Returns: A created variant object
func VariantFromLong(value int64) *Variant {
	c := &Variant{}
	c.SetAsLong(value)
	return c
}

// VariantFromBoolean creates a new variant from Boolean value.
//
//	Parameters:
//		- value: a variant value.
//	Returns: A created variant object
func VariantFromBoolean(value bool) *Variant {
	c := &Variant{}
	c.SetAsBoolean(value)
	return c
}

// VariantFromFloat creates a new variant from Float value.
//
//	Parameters:
//		- value: a variant value.
//	Returns: A created variant object
func VariantFromFloat(value float32) *Variant {
	c := &Variant{}
	c.SetAsFloat(value)
	return c
}

// VariantFromDouble creates a new variant from Double value.
//
//	Parameters:
//		- value: a variant value.
//	Returns: A created variant object
func VariantFromDouble(value float64) *Variant {
	c := &Variant{}
	c.SetAsDouble(value)
	return c
}

// VariantFromString creates a new variant from String value.
//
//	Parameters:
//		- value: a variant value.
//	Returns: A created variant object
func VariantFromString(value string) *Variant {
	c := &Variant{}
	c.SetAsString(value)
	return c
}

// VariantFromDateTime creates a new variant from DateTime value.
//
//	Parameters:
//		- value: a variant value.
//	Returns: A created variant object
func VariantFromDateTime(value time.Time) *Variant {
	c := &Variant{}
	c.SetAsDateTime(value)
	return c
}

// VariantFromTimeSpan creates a new variant from TimeSpan value.
//
//	Parameters:
//		- value: a variant value.
//	Returns: A created variant object
func VariantFromTimeSpan(value time.Duration) *Variant {
	c := &Variant{}
	c.SetAsTimeSpan(value)
	return c
}

// VariantFromObject creates a new variant from Object value.
//
//	Parameters:
//		- value: a variant value.
//	Returns: A created variant object
func VariantFromObject(value any) *Variant {
	c := &Variant{}
	c.SetAsObject(value)
	return c
}

// VariantFromArray creates a new variant from Array value.
//
//	Parameters:
//		- value: a variant value.
//	Returns: A created variant object
func VariantFromArray(value []*Variant) *Variant {
	c := &Variant{}
	c.SetAsArray(value)
	return c
}

// Type gets a variant type
func (c *Variant) Type() VariantType {
	return c.typ
}

// AsInteger gets variant value as integer
func (c *Variant) AsInteger() int {
	return c.value.(int)
}

// SetAsInteger sets variant value as integer
//
//	Parameters:
//		- value a value to be set
func (c *Variant) SetAsInteger(value int) {
	c.typ = Integer
	c.value = value
}

// AsLong gets variant value as int64
func (c *Variant) AsLong() int64 {
	return c.value.(int64)
}

// SetAsLong sets variant value as int64
//
//	Parameters:
//		- value a value to be set
func (c *Variant) SetAsLong(value int64) {
	c.typ = Long
	c.value = value
}

// AsBoolean gets variant value as boolean
func (c *Variant) AsBoolean() bool {
	return c.value.(bool)
}

// SetAsBoolean sets variant value as boolean
//
//	Parameters:
//		- value a value to be set
func (c *Variant) SetAsBoolean(value bool) {
	c.typ = Boolean
	c.value = value
}

// AsFloat gets variant value as float
func (c *Variant) AsFloat() float32 {
	return c.value.(float32)
}

// SetAsFloat sets variant value as float
//
//	Parameters:
//		- value a value to be set
func (c *Variant) SetAsFloat(value float32) {
	c.typ = Float
	c.value = value
}

// AsDouble gets variant value as double
func (c *Variant) AsDouble() float64 {
	return c.value.(float64)
}

// SetAsDouble sets variant value as double
//
//	Parameters:
//		- value a value to be set
func (c *Variant) SetAsDouble(value float64) {
	c.typ = Double
	c.value = value
}

// AsString gets variant value as string
func (c *Variant) AsString() string {
	return c.value.(string)
}

// SetAsString sets variant value as string
//
//	Parameters:
//		- value a value to be set
func (c *Variant) SetAsString(value string) {
	c.typ = String
	c.value = value
}

// AsDateTime gets variant value as DateTime
func (c *Variant) AsDateTime() time.Time {
	return c.value.(time.Time)
}

// SetAsDateTime sets variant value as DateTime
//
//	Parameters:
//		- value a value to be set
func (c *Variant) SetAsDateTime(value time.Time) {
	c.typ = DateTime
	c.value = value
}

// AsTimeSpan gets variant value as TimeSpan
func (c *Variant) AsTimeSpan() time.Duration {
	return c.value.(time.Duration)
}

// SetAsTimeSpan sets variant value as TimeSpan
//
//	Parameters:
//		- value a value to be set
func (c *Variant) SetAsTimeSpan(value time.Duration) {
	c.typ = TimeSpan
	c.value = value
}

// AsObject gets variant value as object
func (c *Variant) AsObject() any {
	return c.value
}

// SetAsObject sets variant value as object
//
//	Parameters:
//		- value a value to be set
func (c *Variant) SetAsObject(value any) {
	c.value = value

	if value == nil {
		c.typ = Null
		return
	}

	switch value.(type) {
	case int:
		c.typ = Integer
	case int32:
		_val, _ := c.value.(int32)
		c.value = int(_val)
		c.typ = Integer
	case uint:
		_val, _ := c.value.(uint)
		c.value = int64(_val)
		c.typ = Long
	case uint32:
		_val, _ := c.value.(uint32)
		c.value = int64(_val)
		c.typ = Long
	case int64:
		c.typ = Long
	case float32:
		c.typ = Float
	case float64:
		c.typ = Double
	case bool:
		c.typ = Boolean
	case time.Time:
		c.typ = DateTime
	case time.Duration:
		c.typ = TimeSpan
	case string:
		c.typ = String
	case []*Variant:
		c.typ = Array
		v1, _ := c.value.([]*Variant)
		v2 := make([]*Variant, len(v1))
		copy(v2, v1)
		c.value = v2
	case *Variant:
		v, _ := c.value.(*Variant)
		c.typ = v.typ
		c.value = v.value
	default:
		c.typ = Object
	}
}

// AsArray gets variant value as variant array
func (c *Variant) AsArray() []*Variant {
	if val, ok := c.value.([]*Variant); ok {
		return val
	}
	return nil
}

// SetAsArray sets variant value as variant array
//
//	Parameters:
//		- value a value to be set
func (c *Variant) SetAsArray(value []*Variant) {
	c.typ = Array
	a := make([]*Variant, len(value))
	copy(a, value)
	c.value = a
}

// Length gets length of the array
//
//	Returns: the length of the array or 0
func (c *Variant) Length() int {
	if c.typ == Array {
		return len(c.value.([]*Variant))
	}
	return 0
}

// SetLength sets a new array length
//
//	Parameters:
//		- value a new array length
func (c *Variant) SetLength(value int) {
	if c.typ == Array {
		a := c.value.([]*Variant)
		for len(a) < value {
			a = append(a, &Variant{typ: Null, value: nil})
		}
		c.value = a
	} else {
		panic("Cannot set array length for non-array data type.")
	}
}

// GetByIndex пets an array element by its index.
//
//	Parameters:
//		- index an element index
//	Returns: a requested array element
func (c *Variant) GetByIndex(index int) *Variant {
	if c.typ == Array {
		a := c.value.([]*Variant)
		if len(a) > index {
			return a[index]
		} else {
			panic("Requested element of array is not accessible.")
		}
	} else {
		panic("Cannot access array element for none-array data type.")
	}
}

// SetByIndex sets an array element by its index.
//
//	Parameters:
//		- index an element index
//		- element an element value
func (c *Variant) SetByIndex(index int, element *Variant) {
	if c.typ == Array {
		a := c.value.([]*Variant)
		for len(a) <= index {
			a = append(a, &Variant{typ: Null, value: nil})
		}
		a[index] = element
		c.value = a
	} else {
		panic("Cannot access array element for none-array data type.")
	}
}

// IsNull checks is this variant value Null.
//
//	Returns: <code>true</code> if this variant value is Null.
func (c *Variant) IsNull() bool {
	return c.typ == Null
}

// IsEmpty checks is this variant value empty.
//
//	Returns <code>true</code< is this variant value is empty.
func (c *Variant) IsEmpty() bool {
	return c.value == nil
}

// Assign a new value to this object.
//
//	Parameters:
//		- value A new value to be assigned.
func (c *Variant) Assign(value *Variant) {
	if value != nil {
		c.typ = value.typ
		c.value = value.value
	} else {
		c.typ = Null
		c.value = nil
	}
}

// Clear this object and assigns a VariantType.Null type.
func (c *Variant) Clear() {
	c.typ = Null
	c.value = nil
}

// String gets a string value for this object.
//
//	Returns: a string value for this object.
func (c *Variant) String() string {
	if c.value == nil {
		return "null"
	}
	return cconv.StringConverter.ToString(c.value)
}

// Equals compares this object to the specified one.
//
//	Parameters:
//		- obj An object to be compared.
//	Returns: <code>true</code> if objects are equal.
func (c *Variant) Equals(obj *Variant) bool {
	if obj == nil {
		return false
	}
	value1 := c.value
	value2 := obj.value
	if value1 == nil || value2 == nil {
		return value1 == value2
	}
	return c.typ == obj.typ && value1 == value2
}

// Clone the variant value
//
//	Returns: The cloned value of this variant
func (c *Variant) Clone() *Variant {
	return NewVariant(c)
}
