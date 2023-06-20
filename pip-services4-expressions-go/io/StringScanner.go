package io

// StringScanner scan characters in a string that allows tokenizers
// to look ahead through stream to perform lexical analysis.
type StringScanner struct {
	content  []rune
	position int
	line     int
	column   int
}

// NewStringScanner creates an instance of this class.
//	Parameters:
//		content a text content to be read
func NewStringScanner(content string) *StringScanner {
	c := StringScanner{
		content:  []rune(content),
		position: -1,
		line:     1,
		column:   0,
	}
	return &c
}

// charAt returns character from a specified position in the stream
//	Parameters:
//		position a position to read character
//	Returns: a character from the specified position or EOF (-1)
func (c *StringScanner) charAt(position int) rune {
	if position < 0 || position >= len(c.content) {
		return -1
	}

	return c.content[position]
}

// isLine checks if the current character represents a new line
//	Parameters:
//		charBefore the character before the current one
//		charAt the current character
//		charAfter the character after the current one
//	Returns: <code>true</code> if the current character is a new line, or <code>false</code> otherwise.
func (c *StringScanner) isLine(charBefore rune, charAt rune, charAfter rune) bool {
	if charAt != '\n' && charAt != '\r' {
		return false
	}
	if charAt == '\r' && (charBefore == '\n' || charAfter == '\n') {
		return false
	}
	return true
}

// isColumn сhecks if the current character represents a column
//	Parameters:
//		charAt the current character
//	Returns: <code>true</code> if the current character is a column, or <code>false</code> otherwise.
func (c *StringScanner) isColumn(charAt rune) bool {
	if charAt == '\n' || charAt == '\r' {
		return false
	}
	return true
}

// Read character from the top of the stream.
//	ReturnsЖ a read character or <code>-1</code> if stream processed to the end.</returns>
func (c *StringScanner) Read() rune {
	// Skip if we are at the end
	if (c.position + 1) > len(c.content) {
		return -1
	}

	// Update the current position
	c.position++

	if c.position >= len(c.content) {
		return -1
	}

	// Update line and columns
	charBefore := c.charAt(c.position - 1)
	charAt := c.charAt(c.position)
	charAfter := c.charAt(c.position + 1)

	if c.isLine(charBefore, charAt, charAfter) {
		c.line++
		c.column = 0
	}
	if c.isColumn(charAt) {
		c.column++
	}

	return charAt
}

// Line gets the current line number
//	Returns: the current line number in the stream
func (c *StringScanner) Line() int {
	return c.line
}

// Column gets the column in the current line
//	Returns: the column in the current line in the stream
func (c *StringScanner) Column() int {
	return c.column
}

// Peek returns the character from the top of the stream without moving the stream pointer.
//	Returns: a character from the top of the stream or <code>-1</code> if stream is empty.</returns>
func (c *StringScanner) Peek() rune {
	return c.charAt(c.position + 1)
}

// PeekLine gets the next character line number
//	Returns: the next character line number in the stream
func (c *StringScanner) PeekLine() int {
	charBefore := c.charAt(c.position)
	charAt := c.charAt(c.position + 1)
	charAfter := c.charAt(c.position + 2)

	if c.isLine(charBefore, charAt, charAfter) {
		return c.line + 1
	}
	return c.line
}

// PeekColumn gets the next character column
//	Returns: the next character column in the stream
func (c *StringScanner) PeekColumn() int {
	charBefore := c.charAt(c.position)
	charAt := c.charAt(c.position + 1)
	charAfter := c.charAt(c.position + 2)

	if c.isLine(charBefore, charAt, charAfter) {
		return 0
	}

	if c.isColumn(charAt) {
		return c.column + 1
	}
	return c.column
}

// Unread puts the specified character to the top of the stream.
func (c *StringScanner) Unread() {
	// Skip if we are at the beginning
	if c.position < -1 {
		return
	}

	// Update the current position
	c.position--

	// Update line and columns (optimization)
	if c.column > 0 {
		c.column--
		return
	}

	// Update line and columns (full version)
	c.line = 1
	c.column = 0

	charBefore := rune(-1)
	charAt := rune(-1)
	charAfter := c.charAt(0)

	for position := 0; position <= c.position; position++ {
		charBefore = charAt
		charAt = charAfter
		charAfter = c.charAt(position + 1)

		if c.isLine(charBefore, charAt, charAfter) {
			c.line++
			c.column = 0
		}
		if c.isColumn(charAt) {
			c.column++
		}
	}
}

// UnreadMany puts the specified number of characters to the top of the stream.
//	Parameters:
//		- count A number of characters to be unread
func (c *StringScanner) UnreadMany(count int) {
	for count > 0 {
		c.Unread()
		count--
	}
}

// Reset scanner to the initial position
func (c *StringScanner) Reset() {
	c.position = -1
	c.line = 1
	c.column = 0
}
