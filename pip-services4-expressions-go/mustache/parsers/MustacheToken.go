package parsers

// MustacheToken defines an mustache token holder.
type MustacheToken struct {
	typ    int
	value  string
	tokens []*MustacheToken
	line   int
	column int
}

// NewMustacheToken creates an instance of a mustache token
//	Parameters:
//		- typ: The type of this token.
//		- value: The value of this token.
//		- line: The line number where the token is.
//		- column: The column number where the token is.
func NewMustacheToken(typ int, value string, line int, column int) *MustacheToken {
	c := &MustacheToken{
		typ:    typ,
		value:  value,
		line:   line,
		column: column,
	}
	return c
}

// Type the type of this token.
func (c *MustacheToken) Type() int {
	return c.typ
}

// Value the value of this token.
func (c *MustacheToken) Value() string {
	return c.value
}

// Tokens gets a list of subtokens is this token a section.
func (c *MustacheToken) Tokens() []*MustacheToken {
	return c.tokens
}

// SetTokens sets a list of subtokens is this token a section.
func (c *MustacheToken) SetTokens(tokens []*MustacheToken) {
	c.tokens = tokens
}

// Line the line number where the token is.
func (c *MustacheToken) Line() int {
	return c.line
}

// Column the column number where the token is.
func (c *MustacheToken) Column() int {
	return c.column
}
