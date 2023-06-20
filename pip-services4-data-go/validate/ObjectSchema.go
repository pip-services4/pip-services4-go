package validate

import (
	"strings"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/reflect"
)

// ObjectSchema to validate user defined objects.
//
//	Example:
//		schema.Validate(struct {
//			id   string
//			name string
//		}{id: "1", name: "ABC"}) // Result: no errors
//
//		schema.Validate(struct {
//			id   string
//			name string
//		}{name: "ABC"}) // Result: no errors
//
//		schema.Validate(struct {
//			id   int
//			name string
//		}{id: 1, name: "ABC"}) // Result: id type mismatch
//
//		schema.Validate(struct {
//			id    int
//			_name string
//		}{id: 1, _name: "ABC"}) // Result: name is missing, unexpected_name
//
//		schema.Validate("ABC")
type ObjectSchema struct {
	*Schema
	properties     []*PropertySchema
	allowUndefined bool
}

// NewObjectSchema creates a new validation schema and sets its values.
//
//	Returns: *ObjectSchema
func NewObjectSchema() *ObjectSchema {
	c := &ObjectSchema{
		allowUndefined: false,
	}
	c.Schema = InheritSchema(c)
	return c
}

// NewObjectSchemaWithRules creates a new validation schema and sets its values.
//
//	see IValidationRule
//	Parameters:
//		- allowUndefined bool true to allow properties undefines in the schema
//		- required bool true to always require non-null values.
//		- rules []IValidationRule a list with validation rules.
//	Returns: *ObjectSchema
func NewObjectSchemaWithRules(allowUndefined bool, required bool, rules []IValidationRule) *ObjectSchema {
	c := &ObjectSchema{
		allowUndefined: allowUndefined,
	}
	c.Schema = InheritSchemaWithRules(c, required, rules)
	return c
}

// Properties gets validation schemas for object properties.
//
//	see PropertySchema
//	Returns: []*PropertySchema the list of property validation schemas.
func (c *ObjectSchema) Properties() []*PropertySchema {
	return c.properties
}

// SetProperties sets validation schemas for object properties.
//
//	see PropertySchema
//	Parameters: value []*PropertySchema a list of property validation schemas.
func (c *ObjectSchema) SetProperties(value []*PropertySchema) {
	c.properties = value
}

// UndefinedAllowed gets flag to allow undefined properties
//
//	Returns: bool true to allow undefined properties and false to disallow.
func (c *ObjectSchema) UndefinedAllowed() bool {
	return c.allowUndefined
}

// SetUndefinedAllowed gets flag to allow undefined properties
//
//	Parameters: value bool true to allow undefined properties and false to disallow.
func (c *ObjectSchema) SetUndefinedAllowed(value bool) {
	c.allowUndefined = value
}

// AllowUndefined sets flag to allow undefined properties
// This method returns reference to this exception to implement Builder pattern to chain additional calls.
//
//	Parameters: value bool true to allow undefined properties and false to disallow.
//	Returns: *ObjectSchema this validation schema.
func (c *ObjectSchema) AllowUndefined(value bool) *ObjectSchema {
	c.allowUndefined = value
	return c
}

// WithProperty adds a validation schema for an object property.
// This method returns reference to this exception to implement
// Builder pattern to chain additional calls.
//
//	see PropertySchema
//	Parameters: schema *PropertySchema a property validation schema to be added.
//	Returns: *ObjectSchema this validation schema.
func (c *ObjectSchema) WithProperty(schema *PropertySchema) *ObjectSchema {
	if c.properties == nil {
		c.properties = make([]*PropertySchema, 0)
	}
	c.properties = append(c.properties, schema)
	return c
}

// WithRequiredProperty adds a validation schema for a required object property.
//
//	Parameters:
//		- name string a property name.
//		- type any a property schema or type.
//		- rules ...IValidationRule a list of property validation rules.
//	Returns: *ObjectSchema
func (c *ObjectSchema) WithRequiredProperty(name string, typ any, rules ...IValidationRule) *ObjectSchema {
	return c.WithProperty(NewPropertySchemaWithRules(name, typ, true, rules))
}

// WithOptionalProperty adds a validation schema for an optional object property.
//
//	Parameters:
//		- name string a property name.
//		- type any a property schema or type.
//		- rules ...IValidationRule a list of property validation rules.
//	Returns: *ObjectSchema
func (c *ObjectSchema) WithOptionalProperty(name string, typ any, rules ...IValidationRule) *ObjectSchema {
	return c.WithProperty(NewPropertySchemaWithRules(name, typ, false, rules))
}

// PerformValidation validates a given value against the schema and configured validation rules.
//
//	Parameters:
//		- path string a dot notation path to the value.
//		- value any a value to be validated.
//	Returns: []*ValidationResult a list with validation results to add new results.
func (c *ObjectSchema) PerformValidation(path string, value any) []*ValidationResult {
	results := c.Schema.PerformValidation(path, value)
	if results == nil {
		results = make([]*ValidationResult, 0)
	}

	if value == nil {
		return results
	}

	name := path
	if name == "" {
		name = "value"
	}
	properties := reflect.ObjectReader.GetProperties(value)

	if properties != nil {
		for _, propertySchema := range c.properties {
			processedName := ""

			for propertyName, propertyValue := range properties {
				if strings.EqualFold(propertySchema.Name(), propertyName) {
					propResults := propertySchema.PerformValidation(path, propertyValue)
					if propResults != nil {
						results = append(results, propResults...)
					}
					processedName = propertyName
					break
				}
			}

			if processedName != "" {
				delete(properties, processedName)
			} else {
				propResults := propertySchema.PerformValidation(path, nil)
				if propResults != nil {
					results = append(results, propResults...)
				}
			}
		}
	}

	if !c.allowUndefined {
		for propertyName := range properties {
			propertyPath := propertyName
			if path != "" {
				propertyPath = path + "." + propertyName
			}

			results = append(results, NewValidationResult(
				propertyPath,
				Warning,
				"UNEXPECTED_PROPERTY",
				name+" contains unexpected property "+propertyName,
				nil,
				propertyName,
			))
		}
	}

	return results
}
