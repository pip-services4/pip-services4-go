package validate

import (
	"strings"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/reflect"
)

// AtLeastOneExistsRule validation rule that check that at least one of the object properties is not null.
//
//	see IValidationRule
//	Example:
//		schema := NewSchema()
//			.WithRule(NewAtLeastOneExistsRule("field1", "field2"))
//
//		schema.Validate(struct {
//				field1 int
//				field2 string
//			}{field1: 1, field2: "A"})     // Result: no errors
//		schema.Validate(struct {
//				field1 int
//				field2 string
//			}{field1: 1})                  // Result: no errors
//		schema.Validate(struct {}{})       // Result: at least one of properties field1, field2 must exist
type AtLeastOneExistsRule struct {
	properties []string
}

// NewAtLeastOneExistsRule creates a new validation rule and sets its values
//
//	Parameters: properties ...string a list of property names where at least one property must exist
//	Returns: *AtLeastOneExistsRule
func NewAtLeastOneExistsRule(properties ...string) *AtLeastOneExistsRule {
	return &AtLeastOneExistsRule{
		properties: properties,
	}
}

// Validate validates a given value against this rule.
//
//	Parameters:
//		- path string a dot notation path to the value.
//		- schema ISchema a schema this rule is called from
//		- value any a value to be validated.
//	Returns: []*ValidationResult a list with validation results to add new results.
func (c *AtLeastOneExistsRule) Validate(path string, schema ISchema, value any) []*ValidationResult {
	name := path
	if name == "" {
		name = "value"
	}
	found := 0

	for _, property := range c.properties {
		propertyValue := reflect.ObjectReader.GetProperty(value, property)
		if propertyValue != nil {
			found++
		}
	}

	if found == 0 {
		return []*ValidationResult{
			NewValidationResult(
				path,
				Error,
				"VALUE_NULL",
				name+" must have at least one property from "+strings.Join(c.properties, ","),
				c.properties,
				nil,
			),
		}
	}

	return nil
}
