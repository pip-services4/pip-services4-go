package validate

import (
	refl "reflect"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/reflect"
)

// MapSchema to validate maps.
//
//	Example
//		var schema = NewMapSchema(TypeCode.String, TypeCode.Integer);
//		schema.Validate({ "key1": "A", "key2": "B" });       // Result: no errors
//		schema.Validate({ "key1": 1, "key2": 2 });           // Result: element type mismatch
//		schema.Validate([ 1, 2, 3 ]);                        // Result: type mismatch
type MapSchema struct {
	*Schema
	keyType   any
	valueType any
}

// NewMapSchema creates a new instance of validation schema and sets its values.
//
//	see IValidationRule
//	see TypeCode
//	Parameters:
//		- keyType any a type of map keys. Null means that keys may have any type.
//		- valueType any a type of map values. Null means that values may have any type.
//	Returns: *MapSchema
func NewMapSchema(keyType any, valueType any) *MapSchema {
	c := &MapSchema{
		keyType:   keyType,
		valueType: valueType,
	}
	c.Schema = InheritSchema(c)
	return c
}

// NewMapSchemaWithRules creates a new instance of validation schema and sets its values.
//
//	see IValidationRule
//	see TypeCode
//	Parameters:
//		- keyType any a type of map keys. Null means that keys may have any type.
//		- valueType any a type of map values. Null means that values may have any type.
//		- required: boolean true to always require non-null values.
//		- rules: []IValidationRule a list with validation rules.
//	Returns: *MapSchema
func NewMapSchemaWithRules(keyType any, valueType any, required bool, rules []IValidationRule) *MapSchema {
	c := &MapSchema{
		keyType:   keyType,
		valueType: valueType,
	}
	c.Schema = InheritSchemaWithRules(c, required, rules)
	return c
}

// KeyType gets the type of map keys. Null means that keys may have any type.
//
//	Returns: any the type of map keys.
func (c *MapSchema) KeyType() any {
	return c.keyType
}

// SetKeyType sets the type of map keys. Null means that keys may have any type.
//
//	Parameters: value any a type of map keys.
func (c *MapSchema) SetKeyType(value any) {
	c.keyType = value
}

// ValueType gets the type of map values. Null means that values may have any type.
//
//	Returns: any the type of map values.
func (c *MapSchema) ValueType() any {
	return c.valueType
}

// SetValueType sets the type of map values. Null means that values may have any type.
//
//	Parameters: value any a type of map values.
func (c *MapSchema) SetValueType(value any) {
	c.valueType = value
}

// PerformValidation validates a given value against the schema and configured validation rules.
//
//	Parameters:
//		- path string a dot notation path to the value.
//		- value any a value to be validated.
//	Returns: []*ValidationResult[] a list with validation results to add new results.
func (c *MapSchema) PerformValidation(path string, value any) []*ValidationResult {
	value = reflect.ObjectReader.GetValue(value)

	results := c.Schema.PerformValidation(path, value)
	if results == nil {
		results = []*ValidationResult{}
	}

	if value == nil {
		return results
	}

	name := path
	if name == "" {
		name = "value"
	}

	val := refl.ValueOf(value)

	if val.Kind() == refl.Map {
		for _, keyVal := range val.MapKeys() {
			elementPath := convert.StringConverter.ToString(keyVal.Interface())
			if path != "" {
				elementPath = path + "." + elementPath
			}

			keyResults := c.PerformTypeValidation(elementPath, c.keyType, keyVal.Interface())
			if keyResults != nil {
				results = append(results, keyResults...)
			}

			elemResults := c.PerformTypeValidation(elementPath, c.valueType, val.MapIndex(keyVal).Interface())
			if elemResults != nil {
				results = append(results, elemResults...)
			}
		}
	} else {
		if c.Required() {
			results = append(results,
				NewValidationResult(
					path,
					Error,
					"VALUE_ISNOT_MAP",
					name+" type must be Map",
					convert.Map,
					convert.TypeConverter.ToTypeCode(value),
				),
			)
		}
	}

	return results
}
