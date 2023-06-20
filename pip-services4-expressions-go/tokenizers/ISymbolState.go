package tokenizers

// ISymbolState defines an interface for tokenizer state that processes delimiters.
type ISymbolState interface {
	ITokenizerState

	// Add a multi-character symbol.
	//	Parameters:
	//		- value: The symbol to add, such as "=:="
	Add(value string, tokenType int)
}
