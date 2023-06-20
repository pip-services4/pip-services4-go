package validate

import (
	refl "reflect"
	"strconv"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/reflect"
)

// ArraySchema to validate arrays.
//
//	Example:
//		schema := NewArraySchema(TypeCode.String);
//		schema.Validate(["A", "B", "C"]);    // Result: no errors
//		schema.Validate([1, 2, 3]);          // Result: element type mismatch
//		schema.Validate("A");                // Result: type mismatch
type ArraySchema struct {
	*Schema
	valueType any
}

// NewArraySchema creates a new instance of validation schema and sets its values.
//
//	see TypeCode
//	Parameters: valueType any a type of array elements.
//		Null means that elements may have any type.
//	Returns: *ArraySchema
func NewArraySchema(valueType any) *ArraySchema {
	c := &ArraySchema{
		valueType: valueType,
	}
	c.Schema = InheritSchema(c)
	return c
}

// ValueType gets the type of array elements. Null means that elements may have any type.
//
//	Returns: any the type of array elements.
func (c *ArraySchema) ValueType() any {
	return c.valueType
}

// SetValueType sets the type of array elements. Null means that elements may have any type.
//
//	Parameters: value any a type of array elements.
func (c *ArraySchema) SetValueType(value any) {
	c.valueType = value
}

// PerformValidation validates a given value against the schema and configured validation rules.
//
//	Parameters:
//		- path string a dot notation path to the value.
//		- value any a value to be validated.
//	Returns: []*ValidationResult  a list with validation results to add new results.
func (c *ArraySchema) PerformValidation(path string, value any) []*ValidationResult {
	name := path
	if name == "" {
		name = "value"
	}
	value = reflect.ObjectReader.GetValue(value)

	results := c.Schema.PerformValidation(path, value)
	if results == nil {
		results = make([]*ValidationResult, 0)
	}

	if value == nil {
		return results
	}

	val := refl.ValueOf(value)
	if val.Kind() == refl.Ptr {
		val = val.Elem()
	}

	if val.Kind() == refl.Slice || val.Kind() == refl.Array {
		for index := 0; index < val.Len(); index++ {
			elementPath := strconv.Itoa(index)
			if path != "" {
				elementPath = path + "." + elementPath
			}
			elemResults := c.PerformTypeValidation(elementPath, c.valueType, val.Index(index).Interface())
			if elemResults != nil {
				results = append(results, elemResults...)
			}
		}
	} else {
		results = append(results,
			NewValidationResult(
				path,
				Error,
				"VALUE_ISNOT_ARRAY",
				name+" type must to be List or Array",
				convert.Array,
				convert.TypeConverter.ToTypeCode(value),
			),
		)
	}

	return results
}
