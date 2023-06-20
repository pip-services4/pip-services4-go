package tokenizers

// Types (categories) of tokens such as "number", "symbol" or "word".
const (
	Unknown = iota
	Eof
	Eol
	Float
	Integer
	HexDecimal
	Number
	Symbol
	Quoted
	Word
	Keyword
	Whitespace
	Comment
	Special
)
