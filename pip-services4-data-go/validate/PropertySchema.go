package validate

// PropertySchema to validate object properties
//	see ObjectSchema
//	Example
//		var schema = NewObjectSchema()
//			.WithProperty(NewPropertySchema("id", TypeCode.String));
//		schema.Validate({ id: "1", name: "ABC" });       // Result: no errors
//		schema.Validate({ name: "ABC" });                // Result: no errors
//		schema.Validate({ id: 1, name: "ABC" });         // Result: id type mismatch
type PropertySchema struct {
	*Schema
	name string
	typ  any
}

// Creates a new validation schema and sets its values.
// Returns *PropertySchema
func NewPropertySchema() *PropertySchema {
	c := &PropertySchema{}
	c.Schema = InheritSchema(c)
	return c
}

// NewPropertySchemaWithRules creates a new validation schema and sets its values.
//	see IValidationRule
//	see TypeCode
//	Parameters:
//		- name string a property name
//		- type any a property type
//		- required bool true to always require non-null values.
//		- rules []IValidationRule a list with validation rules.
//	Returns: *PropertySchema
func NewPropertySchemaWithRules(name string, typ any, required bool, rules []IValidationRule) *PropertySchema {
	c := &PropertySchema{
		name: name,
		typ:  typ,
	}
	c.Schema = InheritSchemaWithRules(c, required, rules)
	return c
}

// Name gets the property name.
//	Returns: string the property name.
func (c *PropertySchema) Name() string {
	return c.name
}

// SetName sets the property name.
//	Parameters: value string a new property name.
func (c *PropertySchema) SetName(value string) {
	c.name = value
}

// Type gets the property type.
//	Returns: any the property type.
func (c *PropertySchema) Type() any {
	return c.typ
}

// SetType sets a new property type. The type can be defined as type, type name or TypeCode
//	Parameters: value any a new property type.
func (c *PropertySchema) SetType(value any) {
	c.typ = value
}

// PerformValidation validates a given value against the schema and configured validation rules.
//	Parameters:
//		- path string a dot notation path to the value.
//		- value any a value to be validated.
//	Returns:  []*ValidationResult a list with validation results to add new results.
func (c *PropertySchema) PerformValidation(path string, value any) []*ValidationResult {
	if path != "" {
		path = path + "." + c.name
	} else {
		path = c.name
	}

	results := make([]*ValidationResult, 0)

	innerResults := c.Schema.PerformValidation(path, value)
	if innerResults != nil {
		results = append(results, innerResults...)
	}

	innerResults = c.Schema.PerformTypeValidation(path, c.Type(), value)
	if innerResults != nil {
		results = append(results, innerResults...)
	}

	return results
}
