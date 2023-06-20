package variables

import "github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/variants"

// Variable implements a variable holder object.
type Variable struct {
	name  string
	value *variants.Variant
}

// EmptyVariable constructs a new empty variable.
//	Parameters:
//		- name: The name of this variable.
func EmptyVariable(name string) *Variable {
	return NewVariable(name, nil)
}

// NewVariable constructs this variable with name and value.
//	Parameters:
//		- name: The name of this variable.
//		- value: The variable value.
func NewVariable(name string, value *variants.Variant) *Variable {
	if name == "" {
		panic("Name parameter cannot be empty")
	}
	if value == nil {
		value = variants.EmptyVariant()
	}
	c := &Variable{
		name:  name,
		value: value,
	}
	return c
}

// Name variable name.
func (c *Variable) Name() string {
	return c.name
}

// Value the variable value.
func (c *Variable) Value() *variants.Variant {
	return c.value
}

// SetValue the variable value.
func (c *Variable) SetValue(value *variants.Variant) {
	c.value = value
}
