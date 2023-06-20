package parsers

import "github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/variants"

// ExpressionToken defines an expression token holder.
type ExpressionToken struct {
	typ    int
	value  *variants.Variant
	line   int
	column int
}

// NewExpressionToken creates an instance of this token and initializes it with specified values.
//	Parameters:
//		- typ: The type of this token.
//		- value: The value of this token.
//		- line: The line number where the token is.
//		- column: The column number where the token is.
func NewExpressionToken(typ int, value *variants.Variant, line int, column int) *ExpressionToken {
	if value == nil {
		value = variants.EmptyVariant()
	}

	c := &ExpressionToken{
		typ:    typ,
		value:  value,
		line:   line,
		column: column,
	}
	return c
}

// Type of this token.
func (c *ExpressionToken) Type() int {
	return c.typ
}

// Value of this token.
func (c *ExpressionToken) Value() *variants.Variant {
	return c.value
}

// Line number where the token is.
func (c *ExpressionToken) Line() int {
	return c.line
}

// Column number where the token is.
func (c *ExpressionToken) Column() int {
	return c.column
}
