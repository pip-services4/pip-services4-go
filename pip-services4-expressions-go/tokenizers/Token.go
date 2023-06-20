package tokenizers

// Token a token represents a logical chunk of a string. For example, a typical tokenizer would break
// the string "1.23 &lt;= 12.3" into three tokens: the number 1.23, a less-than-or-equal symbol,
// and the number 12.3. A token is a receptacle, and relies on a tokenizer to decide precisely how
// to divide a string into tokens.
type Token struct {
	typ    int
	value  string
	line   int
	column int
}

// NewToken constructs this token with type and value.
//	Parameters:
//		- typ: The type of this token.
//		- value: The token string value.
//		- line: The line number where the token is.
//		- column: The column number where the token is.
//	Returns: Created token
func NewToken(typ int, value string, line int, column int) *Token {
	return &Token{
		typ:    typ,
		value:  value,
		line:   line,
		column: column,
	}
}

// Type token type.
func (c *Token) Type() int {
	return c.typ
}

// Value token value.
func (c *Token) Value() string {
	return c.value
}

// Line the line number where the token is.
func (c *Token) Line() int {
	return c.line
}

// Column the column number where the token is.
func (c *Token) Column() int {
	return c.column
}

func (c *Token) Equals(obj *Token) bool {
	if obj != nil {
		return c.typ == obj.typ && c.value == obj.value
	}
	return false
}
