package validate

import (
	refl "reflect"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"
	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/reflect"
)

// Schema basic schema that validates values against a set of validation rules.
// This schema is used as a basis for specific schemas to validate
// objects, project properties, arrays and maps.
type Schema struct {
	base     ISchemaBase
	required bool
	rules    []IValidationRule
}

// NewSchema creates a new instance of validation schema and sets its values.
//
//	Returns: *Schema
func NewSchema() *Schema {
	c := &Schema{
		required: false,
		rules:    make([]IValidationRule, 0),
	}
	c.base = c
	return c
}

// NewSchemaWithRules creates a new instance of validation schema and sets its values.
//
//	see IValidationRule
//	Parameters:
//		- required bool true to always require non-null values.
//		- rules []IValidationRule a list with validation rules.
//	Returns: *Schema
func NewSchemaWithRules(required bool, rules []IValidationRule) *Schema {
	c := &Schema{
		required: required,
		rules:    rules,
	}
	c.base = c
	return c
}

// InheritSchema inherit schema
//
//	Parameters: base ISchemaBase base foe create new schema
//	Returns: *Schema
func InheritSchema(base ISchemaBase) *Schema {
	c := &Schema{
		required: false,
		rules:    []IValidationRule{},
	}
	c.base = base
	return c
}

// InheritSchemaWithRules inherit schema with rules
//
//	Parameters:
//		- base ISchemaBase base foe create new schema
//		- required bool true to always require non-null values.
//		- rules []IValidationRule a list with validation rules.
//	Returns: *Schema
func InheritSchemaWithRules(base ISchemaBase, required bool, rules []IValidationRule) *Schema {
	c := &Schema{
		required: required,
		rules:    rules,
	}
	c.base = base
	return c
}

// Required gets a flag that always requires non-null values.
// For null values it raises a validation error.
//
//	Returns: bool true to always require non-null values and false to allow null values.
func (c *Schema) Required() bool {
	return c.required
}

// SetRequired sets a flag that always requires non-null values.
//
//	Parameters: value bool true to always require non-null values and false to allow null values.
func (c *Schema) SetRequired(value bool) {
	c.required = value
}

// Rules gets validation rules to check values against.
//
//	Returns: []IValidationRule a list with validation rules.
func (c *Schema) Rules() []IValidationRule {
	return c.rules
}

// SetRules sets validation rules to check values against.
//
//	Parameters: value []IValidationRule a list with validation rules.
func (c *Schema) SetRules(value []IValidationRule) {
	c.rules = value
}

// MakeRequired makes validated values always required (non-null).
// For null values the schema will raise errors.
// This method returns reference to this exception to implement
// Builder pattern to chain additional calls.
//
//	see MakeOptional
//	Returns: *Schema this validation schema
func (c *Schema) MakeRequired() *Schema {
	c.required = true
	return c
}

// MakeOptional makes validated values optional. Validation for null values will be skipped.
// This method returns reference to this exception to
// implement Builder pattern to chain additional calls.
//
//	see MakeRequired
//	Returns: *Schema this validation schema
func (c *Schema) MakeOptional() *Schema {
	c.required = false
	return c
}

// WithRule adds validation rule to this schema.
// This method returns reference to this exception to implement Builder pattern to chain additional calls.
//
//	Parameters: rule IValidationRule a validation rule to be added.
//	Returns: Schema this validation schema.
func (c *Schema) WithRule(rule IValidationRule) *Schema {
	if c.rules == nil {
		c.rules = []IValidationRule{}
	}
	c.rules = append(c.rules, rule)
	return c
}

// PerformValidation validates a given value against the schema and configured validation rules.
//
//	Parameters:
//		- path string a dot notation path to the value.
//		- value any a value to be validated.
//	Returns: []*ValidationResult a list with validation results to add new results.
func (c *Schema) PerformValidation(path string, value any) []*ValidationResult {
	results := make([]*ValidationResult, 0)

	name := path
	if name == "" {
		name = "value"
	}

	if value == nil {
		if c.Required() {
			results = append(results, NewValidationResult(
				path,
				Error,
				"VALUE_IS_NULL",
				name+" must not be null",
				"NOT NULL",
				nil,
			))
		}
	} else {
		value = reflect.ObjectReader.GetValue(value)

		// Check validation rules
		if c.rules != nil {
			for _, rule := range c.rules {
				ruleResults := rule.Validate(path, c, value)
				if ruleResults != nil {
					results = append(results, ruleResults...)
				}
			}
		}
	}

	return results
}

func (c *Schema) typeToString(typ any) string {
	if typ == nil {
		return "unknown"
	}
	if typeCode, ok := convert.IntegerConverter.ToNullableInteger(typ); ok {
		return convert.TypeConverter.ToString(convert.TypeCode(typeCode))
	}
	return convert.StringConverter.ToString(typ)
}

// PerformTypeValidation validates a given value to match specified type.
// The type can be defined as a Schema, type, a type name or
// TypeCode When type is a Schema, it executes validation recursively against that Schema.
//
//	see PerformValidation
//	Parameters:
//		- path string a dot notation path to the value.
//		- type any a type to match the value type
//		- value any a value to be validated.
//	Returns: []*ValidationResult a list with validation results to add new results.
func (c *Schema) PerformTypeValidation(path string, typ any, value any) []*ValidationResult {
	results := make([]*ValidationResult, 0)

	// If type it not defined then skip
	if typ == nil {
		return results
	}

	// Perform validation against the schema
	if schema, ok := typ.(ISchemaBase); ok {
		schemaResults := schema.PerformValidation(path, value)
		if schemaResults != nil {
			results = append(results, schemaResults...)
		}
		return results
	}

	// If value is null then skip
	if value = reflect.ObjectReader.GetValue(value); value == nil {
		return results
	}

	name := path
	if name == "" {
		name = "value"
	}
	valueType := refl.TypeOf(value)
	valueTypeCode := convert.TypeConverter.ToTypeCode(value)

	// Match types
	if reflect.TypeMatcher.MatchType(typ, valueType) {
		return results
	}

	results = append(results,
		NewValidationResult(
			path,
			Error,
			"TYPE_MISMATCH",
			name+" type must be "+c.typeToString(typ)+" but found "+c.typeToString(valueType),
			typ,
			convert.TypeConverter.ToString(valueTypeCode),
		),
	)
	return results
}

// Validate validates the given value and results validation results.
//
//	see ValidationResult
//	Parameters: value any a value to be validated.
//	Returns: []*ValidationResult a list with validation results.
func (c *Schema) Validate(value any) []*ValidationResult {
	return c.base.PerformValidation("", value)
}

// ValidateAndReturnError validates the given value and returns a *errors.ApplicationError if errors were found.
//
//	Parameters:
//		- traceId string transaction id to trace execution through call chain.
//		- value any a value to be validated.
//		- strict bool true to treat warnings as errors.
//	Returns: *errors.ApplicationError
func (c *Schema) ValidateAndReturnError(traceId string, value any, strict bool) *errors.ApplicationError {
	return NewValidationErrorFromResults(traceId, c.Validate(value), strict)
}

// ValidateAndThrowError validates the given value and throws a *errors.ApplicationError if errors were found.
//
//	see errors.ApplicationError.ThrowExceptionIfNeeded
//	Parameters:
//		- traceId string transaction id to trace execution through call chain.
//		- value any a value to be validated.
//		- strict: bool true to treat warnings as errors.
func (c *Schema) ValidateAndThrowError(traceId string, value any, strict bool) {
	results := c.Validate(value)
	ThrowValidationErrorIfNeeded(traceId, results, strict)
}
