package validate

// AndRule validation rule to combine rules with AND logical operation.
// When all rules returns no errors, than this rule also returns no errors.
// When one of the rules return errors, than the rules returns all errors.
//	Example:
//		schema = NewSchema()
//			.WithRule(NewAndRule(
//				NewValueComparisonRule("GTE", 1),
//				NewValueComparisonRule("LTE", 10),
//			));
//		schema.Validate(0);  // Result: 0 must be greater or equal to 1
//		schema.Validate(5);  // Result: no error
//		schema.Validate(20); // Result: 20 must be letter or equal 10
type AndRule struct {
	rules []IValidationRule
}

// NewAndRule creates a new validation rule and sets its values.
//	Parameters: rules ...IValidationRule a list of rules to join with AND operator
//	Returns: *AndRule
func NewAndRule(rules ...IValidationRule) *AndRule {
	return &AndRule{
		rules: rules,
	}
}

// Validate validates a given value against this rule.
//	Parameters:
//		- path string a dot notation path to the value.
//		- schema ISchema a schema this rule is called from
//		- value any a value to be validated.
//	Returns: []*ValidationResult a list with validation results to add new results.
func (c *AndRule) Validate(path string, schema ISchema, value any) []*ValidationResult {
	if len(c.rules) == 0 {
		return nil
	}

	results := make([]*ValidationResult, 0)

	for _, rule := range c.rules {
		ruleResults := rule.Validate(path, schema, value)
		if ruleResults != nil {
			results = append(results, ruleResults...)
		}
	}

	return results
}
