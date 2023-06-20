package tokenizers

// IWhitespaceState defines an interface for tokenizer state that processes whitespaces (' ', '\t')
type IWhitespaceState interface {
	ITokenizerState

	// SetWhitespaceChars establish the given characters as whitespace to ignore.
	//	Parameters:
	//		- fromSymbol: First character index of the interval.
	//		- toSymbol: Last character index of the interval.
	//		- enable: <code>true</code> if this state should ignore characters in the given range.
	SetWhitespaceChars(fromSymbol rune, toSymbol rune, enable bool)

	// ClearWhitespaceChars clears definitions of whitespace characters.
	ClearWhitespaceChars()
}
