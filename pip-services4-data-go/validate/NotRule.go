package validate

// NotRule validation rule negate another rule.
// When embedded rule returns no errors, than this rule return an error.
// When embedded rule return errors, than the rule returns no errors.
//	see IValidationRule
//	Example:
//		var schema = NewSchema()
//			.WithRule(NewNotRule(
//				NewValueComparisonRule("EQ", 1),
//			));
//		schema.Validate(1);  // Result: error
//		schema.Validate(5);  // Result: no error
type NotRule struct {
	rule IValidationRule
}

// NewNotRule creates a new validation rule and sets its values
//	Parameters: rule IValidationRule a rule to be negated.
//	Returns: *NotRule
func NewNotRule(rule IValidationRule) *NotRule {
	return &NotRule{
		rule: rule,
	}
}

// Validate validates a given value against this rule.
//	Parameters:
//		- path string a dot notation path to the value.
//		- schema ISchema a schema this rule is called from
//		- value any a value to be validated.
//	Returns: []*ValidationResult a list with validation results to add new results.
func (c *NotRule) Validate(path string, schema ISchema, value any) []*ValidationResult {
	if c.rule == nil {
		return nil
	}

	name := path
	if name == "" {
		name = "value"
	}

	results := c.rule.Validate(path, schema, value)

	if len(results) > 0 {
		return nil
	}

	return []*ValidationResult{
		NewValidationResult(
			path,
			Error,
			"NOT_FAILED",
			"Negative check for "+name+" failed",
			nil,
			nil,
		),
	}
}
