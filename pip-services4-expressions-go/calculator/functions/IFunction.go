package functions

import "github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/variants"

// IFunction вefines an interface for expression function.
// </summary>
type IFunction interface {
	// Name The function name.
	Name() string

	// Calculate the function calculation method.
	// Parameters:
	//		- parameters: A list with function parameters<
	//		- variantOperations: Variants operations manager.
	Calculate(parameters []*variants.Variant,
		variantOperations variants.IVariantOperations) (*variants.Variant, error)
}
