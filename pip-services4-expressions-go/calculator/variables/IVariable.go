package variables

import "github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/variants"

// IVariable defines a variable interface.
type IVariable interface {
	// Name the variable name.
	Name() string

	// Value gets the variable value.
	Value() *variants.Variant

	// SetValue sets the variable value.
	SetValue(value *variants.Variant)
}
