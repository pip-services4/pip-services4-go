package variants

// IVariantOperations defines an interface for variant operations manager.
type IVariantOperations interface {
	// Convert variant to specified type
	//	Parameters:
	//		- value: A variant value to be converted.
	//		- newType: A type of object to be returned.
	//	Returns: A converted Variant value.
	Convert(value *Variant, newType VariantType) (*Variant, error)

	// Add performs '+' operation for two variants.
	//	Parameters:
	//		- value1: The first operand for this operation.
	//		- value2: The second operand for this operation.
	//	Returns: A result variant object.
	Add(value1 *Variant, value2 *Variant) (*Variant, error)

	// Sub performs '-' operation for two variants.
	//	Parameters:
	//		- value1: The first operand for this operation.
	//		- value2: The second operand for this operation.
	//	Returns: A result variant object.
	Sub(value1 *Variant, value2 *Variant) (*Variant, error)

	// Mul performs '*' operation for two variants.
	//	Parameters:
	//		- value1: The first operand for this operation.
	//		- value2: The second operand for this operation.
	//	Returns: A result variant object.
	Mul(value1 *Variant, value2 *Variant) (*Variant, error)

	// Div performs '/' operation for two variants.
	//	Parameters:
	//		- value1: The first operand for this operation.
	//		- value2: The second operand for this operation.
	//	Returns: A result variant object.
	Div(value1 *Variant, value2 *Variant) (*Variant, error)

	// Mod performs '%' operation for two variants.
	//	Parameters:
	//		- value1: The first operand for this operation.
	//		- value2: The second operand for this operation.
	//	Returns: A result variant object.
	Mod(value1 *Variant, value2 *Variant) (*Variant, error)

	// Pow performs '^' operation for two variants.
	//	Parameters:
	//		- value1: The first operand for this operation.
	//		- value2: The second operand for this operation.
	//	Returns: A result variant object.
	Pow(value1 *Variant, value2 *Variant) (*Variant, error)

	// And performs AND operation for two variants.
	//	Parameters:
	//		- value1: The first operand for this operation.
	//		- value2: The second operand for this operation.
	//	Returns: A result variant object.
	And(value1 *Variant, value2 *Variant) (*Variant, error)

	// Or performs OR operation for two variants.
	//	Parameters:
	//		- value1: The first operand for this operation.
	//		- value2: The second operand for this operation.
	//	Returns: A result variant object.
	Or(value1 *Variant, value2 *Variant) (*Variant, error)

	// Xor performs XOR operation for two variants.
	//	Parameters:
	//		- value1: The first operand for this operation.
	//		- value2: The second operand for this operation.
	//	Returns: A result variant object.
	Xor(value1 *Variant, value2 *Variant) (*Variant, error)

	// Lsh performs << operation for two variants.
	//	Parameters:
	//		- value1: The first operand for this operation.
	//		- value2: The second operand for this operation.
	//	Returns: A result variant object.
	Lsh(value1 *Variant, value2 *Variant) (*Variant, error)

	// Rsh performs >> operation for two variants.
	//	Parameters:
	//		- value1: The first operand for this operation.
	//		- value2: The second operand for this operation.
	//	Returns: A result variant object.
	Rsh(value1 *Variant, value2 *Variant) (*Variant, error)

	// Not performs NOT operation for a variant.
	//	Parameters:
	//		- value: The operand for this operation.
	//	Returns: A result variant object.
	Not(value *Variant) (*Variant, error)

	// Negative performs unary '-' operation for a variant.
	//	Parameters:
	//		- value: The operand for this operation.
	//	Returns: A result variant object.
	Negative(value *Variant) (*Variant, error)

	// Equal performs '=' operation for two variants.
	//	Parameters:
	//		- value1: The first operand for this operation.
	//		- value2: The second operand for this operation.
	//	Returns: A result variant object.
	Equal(value1 *Variant, value2 *Variant) (*Variant, error)

	// NotEqual performs '<>' operation for two variants.
	// Parameters:
	//   - value1: The first operand for this operation.
	//   - value2: The second operand for this operation.
	// Returns: A result variant object.
	NotEqual(value1 *Variant, value2 *Variant) (*Variant, error)

	// More performs '>' operation for two variants.
	//	Parameters:
	//		- value1: The first operand for this operation.
	//		- value2: The second operand for this operation.
	//	Returns: A result variant object.
	More(value1 *Variant, value2 *Variant) (*Variant, error)

	// Less performs '<' operation for two variants.
	//	Parameters:
	//		- value1: The first operand for this operation.
	//		- value2: The second operand for this operation.
	//	Returns: A result variant object.
	Less(value1 *Variant, value2 *Variant) (*Variant, error)

	// MoreEqual performs '>=' operation for two variants.
	//	Parameters:
	//		- value1: The first operand for this operation.
	//		- value2: The second operand for this operation.
	//	Returns: A result variant object.
	MoreEqual(value1 *Variant, value2 *Variant) (*Variant, error)

	// LessEqual performs '<=' operation for two variants.
	//	Parameters:
	//		- value1: The first operand for this operation.
	//		- value2: The second operand for this operation.
	//	Returns: A result variant object.
	LessEqual(value1 *Variant, value2 *Variant) (*Variant, error)

	// In performs IN operation for two variants.
	//	Parameters:
	//		- value1: The first operand for this operation.
	//		- value2: The second operand for this operation.
	//	Returns: A result variant object.
	In(value1 *Variant, value2 *Variant) (*Variant, error)

	// GetElement performs [] operation for two variants.
	//	Parameters:
	//		- value1: The first operand for this operation.
	//		- value2: The second operand for this operation.
	//	Returns: A result variant object.
	GetElement(value1 *Variant, value2 *Variant) (*Variant, error)
}
