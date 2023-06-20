package tokenizers

// IWordState defines an interface for tokenizer state that processes words, identificators or keywords
type IWordState interface {
	ITokenizerState

	// SetWordChars establish characters in the given range as valid characters for part of a word after
	// the first character. Note that the tokenizer must determine which characters are valid
	// as the beginning character of a word.
	//	Parameters:
	//		- fromSymbol: First character index of the interval.
	//		- toSymbol: Last character index of the interval.
	//		- enable: <code>true</code> if this state should use characters in the given range.
	SetWordChars(fromSymbol rune, toSymbol rune, enable bool)

	// ClearWordChars clears definitions of word chars.
	ClearWordChars()
}
