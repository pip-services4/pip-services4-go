package variants

import "github.com/pip-services4/pip-services4-go/pip-services4-commons-go/errors"

type IVariantOperationsOverrides interface {
	Convert(value *Variant, newType VariantType) (*Variant, error)
}

// AbstractVariantOperations implements an abstract variant operations manager object.
type AbstractVariantOperations struct {
	Overrides IVariantOperationsOverrides
}

func InheritAbstractVariantOperations(overrides IVariantOperationsOverrides) *AbstractVariantOperations {
	c := AbstractVariantOperations{
		Overrides: overrides,
	}
	return &c
}

// typeToString convert variant type to string representation
//	Parameters:
//		- value: a variant type to be converted.
//	Returns: a string representation of the type.
func typeToString(value VariantType) string {
	switch value {
	case Null:
		return "Null"
	case Integer:
		return "Integer"
	case Long:
		return "Long"
	case Float:
		return "Float"
	case Double:
		return "Double"
	case String:
		return "String"
	case Boolean:
		return "Boolean"
	case DateTime:
		return "DateTime"
	case TimeSpan:
		return "TimeSpan"
	case Object:
		return "Object"
	case Array:
		return "Array"
	default:
		return "Unknown"
	}
}

// Add performs '+' operation for two variants.
//	Parameters:
//		- value1: The first operand for this operation.
//		- value2: The second operand for this operation.
//	Returns: A result variant object.
func (c *AbstractVariantOperations) Add(
	value1 *Variant, value2 *Variant) (*Variant, error) {

	result := EmptyVariant()

	// Processes VariantType.Null values.
	if value1.Type() == Null || value2.Type() == Null {
		return result, nil
	}

	// Converts second operant to the type of the first operand.
	var err error
	value2, err = c.Overrides.Convert(value2, value1.Type())
	if err != nil {
		return nil, err
	}

	// Performs operation.
	switch value1.Type() {
	case Integer:
		result.SetAsInteger(value1.AsInteger() + value2.AsInteger())
		return result, nil
	case Long:
		result.SetAsLong(value1.AsLong() + value2.AsLong())
		return result, nil
	case Float:
		result.SetAsFloat(value1.AsFloat() + value2.AsFloat())
		return result, nil
	case Double:
		result.SetAsDouble(value1.AsDouble() + value2.AsDouble())
		return result, nil
	case TimeSpan:
		result.SetAsTimeSpan(value1.AsTimeSpan() + value2.AsTimeSpan())
		return result, nil
	case String:
		result.SetAsString(value1.AsString() + value2.AsString())
		return result, nil
	}

	err = errors.NewUnsupportedError("", "OP_NOT_SUPPORTED",
		"Operation '+' is not supported for type "+typeToString(value1.Type()))
	return nil, err
}

// Sub performs '-' operation for two variants.
//	Parameters:
//		- value1: The first operand for this operation.
//		- value2: The second operand for this operation.
//	Returns: A result variant object.
func (c *AbstractVariantOperations) Sub(
	value1 *Variant, value2 *Variant) (*Variant, error) {

	result := EmptyVariant()

	// Processes VariantType.Null values.
	if value1.Type() == Null || value2.Type() == Null {
		return result, nil
	}

	// Converts second operant to the type of the first operand.
	var err error
	value2, err = c.Overrides.Convert(value2, value1.Type())
	if err != nil {
		return nil, err
	}

	// Performs operation.
	switch value1.Type() {
	case Integer:
		result.SetAsInteger(value1.AsInteger() - value2.AsInteger())
		return result, nil
	case Long:
		result.SetAsLong(value1.AsLong() - value2.AsLong())
		return result, nil
	case Float:
		result.SetAsFloat(value1.AsFloat() - value2.AsFloat())
		return result, nil
	case Double:
		result.SetAsDouble(value1.AsDouble() - value2.AsDouble())
		return result, nil
	case TimeSpan:
		result.SetAsTimeSpan(value1.AsTimeSpan() - value2.AsTimeSpan())
		return result, nil
	case DateTime:
		result.SetAsTimeSpan(value1.AsDateTime().Sub(value2.AsDateTime()))
		return result, nil
	}

	err = errors.NewUnsupportedError("", "OP_NOT_SUPPORTED",
		"Operation '-' is not supported for type "+typeToString(value1.Type()))
	return nil, err
}

// Mul performs '*' operation for two variants.
//	Parameters:
//		- value1: The first operand for this operation.
//		- value2: The second operand for this operation.
//	Returns: A result variant object.
func (c *AbstractVariantOperations) Mul(
	value1 *Variant, value2 *Variant) (*Variant, error) {

	result := EmptyVariant()

	// Processes VariantType.Null values.
	if value1.Type() == Null || value2.Type() == Null {
		return result, nil
	}

	// Converts second operant to the type of the first operand.
	var err error
	value2, err = c.Overrides.Convert(value2, value1.Type())
	if err != nil {
		return nil, err
	}

	// Performs operation.
	switch value1.Type() {
	case Integer:
		result.SetAsInteger(value1.AsInteger() * value2.AsInteger())
		return result, nil
	case Long:
		result.SetAsLong(value1.AsLong() * value2.AsLong())
		return result, nil
	case Float:
		result.SetAsFloat(value1.AsFloat() * value2.AsFloat())
		return result, nil
	case Double:
		result.SetAsDouble(value1.AsDouble() * value2.AsDouble())
		return result, nil
	}

	err = errors.NewUnsupportedError("", "OP_NOT_SUPPORTED",
		"Operation '*' is not supported for type "+typeToString(value1.Type()))
	return nil, err
}

// Div performs '/' operation for two variants.
//	Parameters:
//		- value1: The first operand for this operation.
//		- value2: The second operand for this operation.
//	Returns: A result variant object.
func (c *AbstractVariantOperations) Div(
	value1 *Variant, value2 *Variant) (*Variant, error) {

	result := EmptyVariant()

	// Processes VariantType.Null values.
	if value1.Type() == Null || value2.Type() == Null {
		return result, nil
	}

	// Converts second operant to the type of the first operand.
	var err error
	value2, err = c.Overrides.Convert(value2, value1.Type())
	if err != nil {
		return nil, err
	}

	// Performs operation.
	switch value1.Type() {
	case Integer:
		result.SetAsInteger(value1.AsInteger() / value2.AsInteger())
		return result, nil
	case Long:
		result.SetAsLong(value1.AsLong() / value2.AsLong())
		return result, nil
	case Float:
		result.SetAsFloat(value1.AsFloat() / value2.AsFloat())
		return result, nil
	case Double:
		result.SetAsDouble(value1.AsDouble() / value2.AsDouble())
		return result, nil
	}

	err = errors.NewUnsupportedError("", "OP_NOT_SUPPORTED",
		"Operation '/' is not supported for type "+typeToString(value1.Type()))
	return nil, err
}

// Mod performs '%' operation for two variants.
//	Parameters:
//		- value1: The first operand for this operation.
//		- value2: The second operand for this operation.
//	Returns: A result variant object.
func (c *AbstractVariantOperations) Mod(
	value1 *Variant, value2 *Variant) (*Variant, error) {

	result := EmptyVariant()

	// Processes VariantType.Null values.
	if value1.Type() == Null || value2.Type() == Null {
		return result, nil
	}

	// Converts second operant to the type of the first operand.
	var err error
	value2, err = c.Overrides.Convert(value2, value1.Type())
	if err != nil {
		return nil, err
	}

	// Performs operation.
	switch value1.Type() {
	case Integer:
		result.SetAsInteger(value1.AsInteger() % value2.AsInteger())
		return result, nil
	case Long:
		result.SetAsLong(value1.AsLong() % value2.AsLong())
		return result, nil
	}

	err = errors.NewUnsupportedError("", "OP_NOT_SUPPORTED",
		"Operation '%' is not supported for type "+typeToString(value1.Type()))
	return nil, err
}

// Pow performs '^' operation for two variants.
//	Parameters:
//		- value1: The first operand for this operation.
//		- value2: The second operand for this operation.
//	Returns: A result variant object.
func (c *AbstractVariantOperations) Pow(
	value1 *Variant, value2 *Variant) (*Variant, error) {
	result := EmptyVariant()

	// Processes VariantType.Null values.
	if value1.Type() == Null || value2.Type() == Null {
		return result, nil
	}

	// Performs operation.
	switch value1.Type() {
	case Integer:
	case Long:
	case Float:
	case Double:
		// Converts second operant to the type of the first operand.
		var err error
		value1, err = c.Overrides.Convert(value1, Double)
		if err != nil {
			return nil, err
		}

		value2, err = c.Overrides.Convert(value2, Double)
		if err != nil {
			return nil, err
		}

		result.SetAsDouble(value1.AsDouble() * value2.AsDouble())
		return result, nil
	}

	err := errors.NewUnsupportedError("",
		"OP_NOT_SUPPORTED", "Operation '^' is not supported for type "+typeToString(value1.Type()))
	return nil, err
}

// And performs AND operation for two variants.
//	Parameters:
//		- value1: The first operand for this operation.
//		- value2: The second operand for this operation.
//	Returns: A result variant object.
func (c *AbstractVariantOperations) And(
	value1 *Variant, value2 *Variant) (*Variant, error) {
	result := EmptyVariant()

	// Processes VariantType.Null values.
	if value1.Type() == Null || value2.Type() == Null {
		return result, nil
	}

	// Converts second operant to the type of the first operand.
	var err error
	value2, err = c.Overrides.Convert(value2, value1.Type())
	if err != nil {
		return nil, err
	}

	// Performs operation.
	switch value1.Type() {
	case Integer:
		result.SetAsInteger(value1.AsInteger() & value2.AsInteger())
		return result, nil
	case Long:
		result.SetAsLong(value1.AsLong() & value2.AsLong())
		return result, nil
	case Boolean:
		result.SetAsBoolean(value1.AsBoolean() && value2.AsBoolean())
		return result, nil
	}

	err = errors.NewUnsupportedError("",
		"OP_NOT_SUPPORTED", "Operation AND is not supported for type "+typeToString(value1.Type()))
	return nil, err
}

// Or performs OR operation for two variants.
//	Parameters:
//		- value1: The first operand for this operation.
//		- value2: The second operand for this operation.
//	Returns: A result variant object.
func (c *AbstractVariantOperations) Or(
	value1 *Variant, value2 *Variant) (*Variant, error) {

	result := EmptyVariant()

	// Processes VariantType.Null values.
	if value1.Type() == Null || value2.Type() == Null {
		return result, nil
	}

	// Converts second operant to the type of the first operand.
	var err error
	value2, err = c.Overrides.Convert(value2, value1.Type())
	if err != nil {
		return nil, err
	}

	// Performs operation.
	switch value1.Type() {
	case Integer:
		result.SetAsInteger(value1.AsInteger() | value2.AsInteger())
		return result, nil
	case Long:
		result.SetAsLong(value1.AsLong() | value2.AsLong())
		return result, nil
	case Boolean:
		result.SetAsBoolean(value1.AsBoolean() || value2.AsBoolean())
		return result, nil
	}

	err = errors.NewUnsupportedError("", "OP_NOT_SUPPORTED",
		"Operation OR is not supported for type "+typeToString(value1.Type()))
	return nil, err
}

// Xor performs XOR operation for two variants.
//	Parameters:
//		- value1: The first operand for this operation.
//		- value2: The second operand for this operation.
//	Returns: A result variant object.
func (c *AbstractVariantOperations) Xor(
	value1 *Variant, value2 *Variant) (*Variant, error) {
	result := EmptyVariant()

	// Processes VariantType.Null values.
	if value1.Type() == Null || value2.Type() == Null {
		return result, nil
	}

	// Converts second operant to the type of the first operand.
	var err error
	value2, err = c.Overrides.Convert(value2, value1.Type())
	if err != nil {
		return nil, err
	}

	// Performs operation.
	switch value1.Type() {
	case Integer:
		result.SetAsInteger(value1.AsInteger() ^ value2.AsInteger())
		return result, nil
	case Long:
		result.SetAsLong(value1.AsLong() ^ value2.AsLong())
		return result, nil
	case Boolean:
		result.SetAsBoolean((value1.AsBoolean() && !value2.AsBoolean()) ||
			(!value1.AsBoolean() && value2.AsBoolean()))
		return result, nil
	}

	err = errors.NewUnsupportedError("",
		"OP_NOT_SUPPORTED", "Operation XOR is not supported for type "+typeToString(value1.Type()))
	return nil, err
}

// Lsh performs '<<' operation for two variants.
//	Parameters:
//		- value1: The first operand for this operation.
//		- value2: The second operand for this operation.
//	Returns: A result variant object.
func (c *AbstractVariantOperations) Lsh(
	value1 *Variant, value2 *Variant) (*Variant, error) {
	result := EmptyVariant()

	// Processes VariantType.Null values.
	if value1.Type() == Null || value2.Type() == Null {
		return result, nil
	}

	// Converts second operant to the type of the first operand.
	var err error
	value2, err = c.Overrides.Convert(value2, Integer)
	if err != nil {
		return nil, err
	}

	// Performs operation.
	switch value1.Type() {
	case Integer:
		result.SetAsInteger(value1.AsInteger() << value2.AsInteger())
		return result, nil
	case Long:
		result.SetAsLong(value1.AsLong() << value2.AsInteger())
		return result, nil
	}

	err = errors.NewUnsupportedError("",
		"OP_NOT_SUPPORTED", "Operation '<<' is not supported for type "+typeToString(value1.Type()))
	return nil, err
}

// Rsh performs '>>' operation for two variants.
//	Parameters:
//		- value1: The first operand for this operation.
//		- value2: The second operand for this operation.
//	Returns: A result variant object.
func (c *AbstractVariantOperations) Rsh(
	value1 *Variant, value2 *Variant) (*Variant, error) {

	result := EmptyVariant()

	// Processes VariantType.Null values.
	if value1.Type() == Null || value2.Type() == Null {
		return result, nil
	}

	// Converts second operant to the type of the first operand.
	var err error
	value2, err = c.Overrides.Convert(value2, Integer)
	if err != nil {
		return nil, err
	}

	// Performs operation.
	switch value1.Type() {
	case Integer:
		result.SetAsInteger(value1.AsInteger() >> value2.AsInteger())
		return result, nil
	case Long:
		result.SetAsLong(value1.AsLong() >> value2.AsInteger())
		return result, nil
	}

	err = errors.NewUnsupportedError("", "OP_NOT_SUPPORTED",
		"Operation '>>' is not supported for type "+typeToString(value1.Type()))
	return nil, err
}

// Not performs NOT operation for a variant.
//	Parameters:
//		- value: The operand for this operation.
//	Returns: A result variant object.
func (c *AbstractVariantOperations) Not(value *Variant) (*Variant, error) {

	result := EmptyVariant()

	// Processes VariantType.Null values.
	if value.Type() == Null {
		result.SetAsBoolean(true)
		return result, nil
	}

	// Performs operation.
	switch value.Type() {
	case Integer:
		result.SetAsInteger(^value.AsInteger())
		return result, nil
	case Long:
		result.SetAsLong(^value.AsLong())
		return result, nil
	case Boolean:
		result.SetAsBoolean(!value.AsBoolean())
		return result, nil
	}

	err := errors.NewUnsupportedError("", "OP_NOT_SUPPORTED",
		"Operation NOT is not supported for type "+typeToString(value.Type()))
	return nil, err
}

// Negative performs unary '-' operation for a variant.
//	Parameters:
//		- value: The operand for this operation.
//	Returns: A result variant object.
func (c *AbstractVariantOperations) Negative(value *Variant) (*Variant, error) {
	result := EmptyVariant()

	// Processes VariantType.Null values.
	if value.Type() == Null {
		return result, nil
	}

	// Performs operation.
	switch value.Type() {
	case Integer:
		result.SetAsInteger(-value.AsInteger())
		return result, nil
	case Long:
		result.SetAsLong(-value.AsLong())
		return result, nil
	case Float:
		result.SetAsFloat(-value.AsFloat())
		return result, nil
	case Double:
		result.SetAsDouble(-value.AsDouble())
		return result, nil
	}

	err := errors.NewUnsupportedError("", "OP_NOT_SUPPORTED",
		"Operation unary '-' is not supported for type "+typeToString(value.Type()))
	return nil, err
}

// Equal performs '=' operation for two variants.
//	Parameters:
//		- value1: The first operand for this operation.
//		- value2: The second operand for this operation.
//	Returns: A result variant object.
func (c *AbstractVariantOperations) Equal(
	value1 *Variant, value2 *Variant) (*Variant, error) {

	result := EmptyVariant()

	// Processes VariantType.Null values.
	if value1.Type() == Null && value2.Type() == Null {
		result.SetAsBoolean(true)
		return result, nil
	}
	if value1.Type() == Null || value2.Type() == Null {
		result.SetAsBoolean(false)
		return result, nil
	}

	// Converts second operant to the type of the first operand.
	var err error
	value2, err = c.Overrides.Convert(value2, value1.Type())
	if err != nil {
		return nil, err
	}

	// Performs operation.
	switch value1.Type() {
	case Integer:
		result.SetAsBoolean(value1.AsInteger() == value2.AsInteger())
		return result, nil
	case Long:
		result.SetAsBoolean(value1.AsLong() == value2.AsLong())
		return result, nil
	case Float:
		result.SetAsBoolean(value1.AsFloat() == value2.AsFloat())
		return result, nil
	case Double:
		result.SetAsBoolean(value1.AsDouble() == value2.AsDouble())
		return result, nil
	case String:
		result.SetAsBoolean(value1.AsString() == value2.AsString())
		return result, nil
	case Boolean:
		result.SetAsBoolean(value1.AsBoolean() == value2.AsBoolean())
		return result, nil
	case TimeSpan:
		result.SetAsBoolean(value1.AsTimeSpan() == value2.AsTimeSpan())
		return result, nil
	case DateTime:
		date1 := value1.AsDateTime()
		date2 := value2.AsDateTime()
		result.SetAsBoolean(date1.Equal(date2))
		return result, nil
	case Object:
		result.SetAsBoolean(value1.AsObject() == value2.AsObject())
		return result, nil
	}

	err = errors.NewUnsupportedError("", "OP_NOT_SUPPORTED",
		"Operation '=' is not supported for type "+typeToString(value1.Type()))
	return nil, err
}

// NotEqual performs '<>' operation for two variants.
//	Parameters:
//		- value1: The first operand for this operation.
//		- value2: The second operand for this operation.
//	Returns: A result variant object.
func (c *AbstractVariantOperations) NotEqual(
	value1 *Variant, value2 *Variant) (*Variant, error) {

	result := EmptyVariant()

	// Processes VariantType.Null values.
	if value1.Type() == Null && value2.Type() == Null {
		result.SetAsBoolean(false)
		return result, nil
	}
	if value1.Type() == Null || value2.Type() == Null {
		result.SetAsBoolean(true)
		return result, nil
	}

	// Converts second operant to the type of the first operand.
	var err error
	value2, err = c.Overrides.Convert(value2, value1.Type())
	if err != nil {
		return nil, err
	}

	// Performs operation.
	switch value1.Type() {
	case Integer:
		result.SetAsBoolean(value1.AsInteger() != value2.AsInteger())
		return result, nil
	case Long:
		result.SetAsBoolean(value1.AsLong() != value2.AsLong())
		return result, nil
	case Float:
		result.SetAsBoolean(value1.AsFloat() != value2.AsFloat())
		return result, nil
	case Double:
		result.SetAsBoolean(value1.AsDouble() != value2.AsDouble())
		return result, nil
	case String:
		result.SetAsBoolean(value1.AsString() != value2.AsString())
		return result, nil
	case Boolean:
		result.SetAsBoolean(value1.AsBoolean() != value2.AsBoolean())
		return result, nil
	case TimeSpan:
		result.SetAsBoolean(value1.AsTimeSpan() != value2.AsTimeSpan())
		return result, nil
	case DateTime:
		date1 := value1.AsDateTime()
		date2 := value2.AsDateTime()
		result.SetAsBoolean(!date1.Equal(date2))
		return result, nil
	case Object:
		result.SetAsBoolean(value1.AsObject() != value2.AsObject())
		return result, nil
	}

	err = errors.NewUnsupportedError("", "OP_NOT_SUPPORTED",
		"Operation '<>' is not supported for type "+typeToString(value1.Type()))
	return nil, err
}

// More performs '>' operation for two variants.
//	Parameters:
//		- value1: The first operand for this operation.
//		- value2: The second operand for this operation.
//	Returns: A result variant object.
func (c *AbstractVariantOperations) More(
	value1 *Variant, value2 *Variant) (*Variant, error) {

	result := EmptyVariant()

	// Processes VariantType.Null values.
	if value1.Type() == Null || value2.Type() == Null {
		return result, nil
	}

	// Converts second operant to the type of the first operand.
	var err error
	value2, err = c.Overrides.Convert(value2, value1.Type())
	if err != nil {
		return nil, err
	}

	// Performs operation.
	switch value1.Type() {
	case Integer:
		result.SetAsBoolean(value1.AsInteger() > value2.AsInteger())
		return result, nil
	case Long:
		result.SetAsBoolean(value1.AsLong() > value2.AsLong())
		return result, nil
	case Float:
		result.SetAsBoolean(value1.AsFloat() > value2.AsFloat())
		return result, nil
	case Double:
		result.SetAsBoolean(value1.AsDouble() > value2.AsDouble())
		return result, nil
	case String:
		result.SetAsBoolean(value1.AsString() > value2.AsString())
		return result, nil
	case TimeSpan:
		result.SetAsBoolean(value1.AsTimeSpan() > value2.AsTimeSpan())
		return result, nil
	case DateTime:
		result.SetAsBoolean(value1.AsDateTime().After(value2.AsDateTime()))
		return result, nil
	}

	err = errors.NewUnsupportedError("", "OP_NOT_SUPPORTED",
		"Operation '>' is not supported for type "+typeToString(value1.Type()))
	return nil, err
}

// Less performs '<' operation for two variants.
//	Parameters:
//		- value1: The first operand for this operation.
//		- value2: The second operand for this operation.
//	Returns: A result variant object.
func (c *AbstractVariantOperations) Less(
	value1 *Variant, value2 *Variant) (*Variant, error) {

	result := EmptyVariant()

	// Processes VariantType.Null values.
	if value1.Type() == Null || value2.Type() == Null {
		return result, nil
	}

	// Converts second operant to the type of the first operand.
	var err error
	value2, err = c.Overrides.Convert(value2, value1.Type())
	if err != nil {
		return nil, err
	}

	// Performs operation.
	switch value1.Type() {
	case Integer:
		result.SetAsBoolean(value1.AsInteger() < value2.AsInteger())
		return result, nil
	case Long:
		result.SetAsBoolean(value1.AsLong() < value2.AsLong())
		return result, nil
	case Float:
		result.SetAsBoolean(value1.AsFloat() < value2.AsFloat())
		return result, nil
	case Double:
		result.SetAsBoolean(value1.AsDouble() < value2.AsDouble())
		return result, nil
	case String:
		result.SetAsBoolean(value1.AsString() < value2.AsString())
		return result, nil
	case TimeSpan:
		result.SetAsBoolean(value1.AsTimeSpan() < value2.AsTimeSpan())
		return result, nil
	case DateTime:
		result.SetAsBoolean(value1.AsDateTime().Before(value2.AsDateTime()))
		return result, nil
	}

	err = errors.NewUnsupportedError("", "OP_NOT_SUPPORTED",
		"Operation '<' is not supported for type "+typeToString(value1.Type()))
	return nil, err
}

// MoreEqual performs '>=' operation for two variants.
//	Parameters:
//		- value1: The first operand for this operation.
//		- value2: The second operand for this operation.
//	Returns: A result variant object.
func (c *AbstractVariantOperations) MoreEqual(
	value1 *Variant, value2 *Variant) (*Variant, error) {

	result := EmptyVariant()

	// Processes VariantType.Null values.
	if value1.Type() == Null || value2.Type() == Null {
		return result, nil
	}

	// Converts second operant to the type of the first operand.
	var err error
	value2, err = c.Overrides.Convert(value2, value1.Type())
	if err != nil {
		return nil, err
	}

	// Performs operation.
	switch value1.Type() {
	case Integer:
		result.SetAsBoolean(value1.AsInteger() >= value2.AsInteger())
		return result, nil
	case Long:
		result.SetAsBoolean(value1.AsLong() >= value2.AsLong())
		return result, nil
	case Float:
		result.SetAsBoolean(value1.AsFloat() >= value2.AsFloat())
		return result, nil
	case Double:
		result.SetAsBoolean(value1.AsDouble() >= value2.AsDouble())
		return result, nil
	case String:
		result.SetAsBoolean(value1.AsString() >= value2.AsString())
		return result, nil
	case TimeSpan:
		result.SetAsBoolean(value1.AsTimeSpan() >= value2.AsTimeSpan())
		return result, nil
	case DateTime:
		date1 := value1.AsDateTime()
		date2 := value2.AsDateTime()
		result.SetAsBoolean(date1.After(date2) || date1.Equal(date2))
		return result, nil
	}

	err = errors.NewUnsupportedError("", "OP_NOT_SUPPORTED",
		"Operation '>=' is not supported for type "+typeToString(value1.Type()))
	return nil, err
}

// LessEqual performs '<=' operation for two variants.
//	Parameters:
//		- value1: The first operand for this operation.
//		- value2: The second operand for this operation.
//	Returns: A result variant object.
func (c *AbstractVariantOperations) LessEqual(
	value1 *Variant, value2 *Variant) (*Variant, error) {
	result := EmptyVariant()

	// Processes VariantType.Null values.
	if value1.Type() == Null || value2.Type() == Null {
		return result, nil
	}

	// Converts second operant to the type of the first operand.
	var err error
	value2, err = c.Overrides.Convert(value2, value1.Type())
	if err != nil {
		return nil, err
	}

	// Performs operation.
	switch value1.Type() {
	case Integer:
		result.SetAsBoolean(value1.AsInteger() <= value2.AsInteger())
		return result, nil
	case Long:
		result.SetAsBoolean(value1.AsLong() <= value2.AsLong())
		return result, nil
	case Float:
		result.SetAsBoolean(value1.AsFloat() <= value2.AsFloat())
		return result, nil
	case Double:
		result.SetAsBoolean(value1.AsDouble() <= value2.AsDouble())
		return result, nil
	case String:
		result.SetAsBoolean(value1.AsString() <= value2.AsString())
		return result, nil
	case TimeSpan:
		result.SetAsBoolean(value1.AsTimeSpan() <= value2.AsTimeSpan())
		return result, nil
	case DateTime:
		date1 := value1.AsDateTime()
		date2 := value2.AsDateTime()
		result.SetAsBoolean(date1.Before(date2) || date1.Equal(date2))
		return result, nil
	}

	err = errors.NewUnsupportedError("", "OP_NOT_SUPPORTED",
		"Operation '<=' is not supported for type "+typeToString(value1.Type()))
	return nil, err
}

// In performs IN operation for two variants.
//	Parameters:
//		- value1: The first operand for this operation.
//		- value2: The second operand for this operation.
//	Returns: A result variant object.
func (c *AbstractVariantOperations) In(
	value1 *Variant, value2 *Variant) (*Variant, error) {

	result := EmptyVariant()

	// Processes VariantType.Null values.
	if value1.Type() == Null || value2.Type() == Null {
		return result, nil
	}

	// Processes null arrays.
	if value1.AsObject() == nil {
		result.SetAsBoolean(false)
		return result, nil
	}

	if value1.Type() == Array {
		array := value1.AsArray()
		for _, element := range array {
			eq, err := c.Equal(value2, element)
			if err != nil {
				return nil, err
			}
			if eq.Type() == Boolean && eq.AsBoolean() {
				result.SetAsBoolean(true)
				return result, nil
			}
		}
		result.SetAsBoolean(false)
		return result, nil
	}

	return c.Equal(value1, value2)
}

// GetElement performs [] operation for two variants.
//	Parameters:
//		- value1: The first operand for this operation.
//		- value2: The second operand for this operation.
//	Returns: A result variant object.
func (c *AbstractVariantOperations) GetElement(
	value1 *Variant, value2 *Variant) (*Variant, error) {
	result := EmptyVariant()

	// Processes VariantType.Null values.
	if value1.Type() == Null || value2.Type() == Null {
		return result, nil
	}

	var err error
	value2, err = c.Overrides.Convert(value2, Integer)
	if err != nil {
		return nil, err
	}

	index := int(value2.AsInteger())

	if value1.Type() == Array {
		return value1.GetByIndex(index), nil
	} else if value1.Type() == String {
		runes := []rune(value1.AsString())
		result.SetAsString(string(runes[value2.AsInteger()]))
		return result, nil
	}

	err = errors.NewUnsupportedError("", "OP_NOT_SUPPORTED",
		"Operation '[]' is not supported for type "+typeToString(value1.Type()))
	return nil, err
}
