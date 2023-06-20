package io

// IScanner defines scanned with ability to unread characters and count lines.
// This scanner is used by tokenizers to process input streams.
type IScanner interface {
	// Read character from the top of the stream.
	//	Returns: a read character or <code>-1</code> if stream processed to the end.</returns>
	Read() rune

	// Line gets the current line number
	//	Returns: the current line number in the stream
	Line() int

	// Column gets the column in the current line
	//	Returns: the column in the current line in the stream
	Column() int

	// Peek returns the character from the top of the stream without moving the stream pointer.
	//	Returns: a character from the top of the stream or <code>-1</code> if stream is empty.</returns>
	Peek() rune

	// PeekLine gets the next character line number
	//	Returns: the next character line number in the stream
	PeekLine() int

	// PeekColumn gets the next character column
	//	Returns the next character column in the stream
	PeekColumn() int

	// Unread puts the specified character to the top of the stream.
	Unread()

	// UnreadMany puts the specified number of characters to the top of the stream.
	//	Parameters:
	//		- count: A number of characters to be unread
	UnreadMany(count int)

	// Reset scanner to the initial position
	Reset()
}
