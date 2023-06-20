package validate

import "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"

// ValueComparisonRule validation rule that compares value to a constant.
//	see IValidationRule
//	Example:
//		var schema = NewSchema()
//			.WithRule(NewValueComparisonRule("EQ", 1));
//		schema.Validate(1);          // Result: no errors
//		schema.Validate(2);          // Result: 2 is not equal to 1
type ValueComparisonRule struct {
	value     any
	operation string
}

// NewValueComparisonRule creates a new validation rule and sets its values.
//	Parameters:
//		- operation string a comparison operation:
//			"==" ("=", "EQ"), "!= " ("<>", "NE"); "<"/">" ("LT"/"GT"), "<="/">=" ("LE"/"GE"); "LIKE".
//		- value any a constant value to compare to
//	Returns: *ValueComparisonRule
func NewValueComparisonRule(operation string, value any) *ValueComparisonRule {
	return &ValueComparisonRule{
		value:     value,
		operation: operation,
	}
}

func (c *ValueComparisonRule) Validate(path string, schema ISchema, value any) []*ValidationResult {
	name := path
	if name == "" {
		name = "value"
	}

	if !ObjectComparator.Compare(value, c.operation, c.value) {
		expectedValue := convert.StringConverter.ToString(c.value)
		actualValue := convert.StringConverter.ToString(value)

		return []*ValidationResult{
			NewValidationResult(
				path,
				Error,
				"BAD_VALUE",
				name+" must "+c.operation+" "+expectedValue+" but found "+actualValue,
				c.operation+" "+expectedValue,
				value,
			),
		}
	}

	return nil
}
