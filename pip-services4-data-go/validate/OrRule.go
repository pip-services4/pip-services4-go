package validate

// OrRule validation rule to combine rules with OR logical operation.
// When one of rules returns no errors, then this rule also returns no errors.
// When all rules return errors, then the rule returns all errors.
//	see IValidationRule
//	Example:
//		var schema = NewObjectSchema().WithProperty(NewPropertySchema("id", TypeCode.String))
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
type OrRule struct {
	rules []IValidationRule
}

// NewOrRule creates a new validation rule and ses its values
//	Parameters: rule IValidationRule a rule to be negated.
//	Returns: *OrRule
func NewOrRule(rules ...IValidationRule) *OrRule {
	return &OrRule{
		rules: rules,
	}
}

// Validate validates a given value against this rule.
//	Parameters:
//		- path string a dot notation path to th value.
//		- schema  ISchema a schema this rule is called from
//		- value any a value to be validated.
//	Returns: []*ValidationResult a list with validation results to add new results.
func (c *OrRule) Validate(path string, schema ISchema, value any) []*ValidationResult {
	if len(c.rules) == 0 {
		return nil
	}

	results := make([]*ValidationResult, 0)

	for _, rule := range c.rules {
		ruleResults := rule.Validate(path, schema, value)

		if len(ruleResults) == 0 {
			return nil
		}

		results = append(results, ruleResults...)
	}

	return results
}
