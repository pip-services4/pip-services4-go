package validate

// IValidationRule interface for validation rules.
// Validation rule can validate one or multiple values against complex
// rules like: value is in range, one property is less than another property,
// enforce enumerated values and more.
// This interface allows implementing custom rules.
type IValidationRule interface {
	// Validate a given value against this rule.
	//	Parameters:
	//		- path string a dot notation path to the value.
	//		- schema ISchema a schema this rule is called from
	//		- value any a value to be validated.
	//	Returns: []*ValidationResult a list with validation results to add new results.
	Validate(path string, schema ISchema, value any) []*ValidationResult
}
