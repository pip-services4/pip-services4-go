package tokenizers

import (
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/io"
)

// ITokenizer a tokenizer divides a string into tokens. This class is highly customizable with regard
// to exactly how this division occurs, but it also has defaults that are suitable for many
// languages. This class assumes that the character values read from the string lie in
// the range 0-255. For example, the Unicode value of a capital A is 65,
// so <code> System.out.println((char)65); </code> prints out a capital A.
// <p>
// The behavior of a tokenizer depends on its character state table. This table is an array
// of 256 <code>TokenizerState</code> states. The state table decides which state to enter
// upon reading a character from the input string.
// <p>
// For example, by default, upon reading an 'A', a tokenizer will enter a "word" state.
// This means the tokenizer will ask a <code>WordState</code> object to consume the 'A',
// along with the characters after the 'A' that form a word. The state's responsibility
// is to consume characters and return a complete token.
// <p>
// The default table sets a SymbolState for every character from 0 to 255,
// and then overrides this with:<blockquote><pre>
// From    To     State
// 0     ' '    whitespaceState
// 'a'    'z'    wordState
// 'A'    'Z'    wordState
// 160     255    wordState
// '0'    '9'    numberState
// '-'    '-'    numberState
// '.'    '.'    numberState
// '"'    '"'    quoteState
// '\”   '\”    quoteState
// '/'    '/'    slashState
// </pre></blockquote>
// In addition to allowing modification of the state table, this class makes each of the states
// above available. Some of these states are customizable. For example, wordState allows customization
// of what characters can be part of a word, after the first character.
type ITokenizer interface {
	// SkipUnknown gets skip unknown characters flag.
	SkipUnknown() bool

	// SetSkipUnknown sets skip unknown characters flag.
	SetSkipUnknown(value bool)

	// SkipWhitespaces gets skip whitespaces flag.
	SkipWhitespaces() bool

	// SetSkipWhitespaces sets skip whitespaces flag.
	SetSkipWhitespaces(value bool)

	// SkipComments gets skip comments flag.
	SkipComments() bool

	// SetSkipComments sets skip comments flag.
	SetSkipComments(value bool)

	// SkipEof gets skip End-Of-File token at the end of stream flag.
	SkipEof() bool

	// SetSkipEof sets skip End-Of-File token at the end of stream flag.
	SetSkipEof(value bool)

	// MergeWhitespaces gets merges whitespaces flag.
	MergeWhitespaces() bool

	// SetMergeWhitespaces sets merges whitespaces flag.
	SetMergeWhitespaces(value bool)

	// UnifyNumbers gets unifies numbers: "Integers" and "Floats" makes just "Numbers" flag
	UnifyNumbers() bool

	// SetUnifyNumbers sets unifies numbers: "Integers" and "Floats" makes just "Numbers" flag
	SetUnifyNumbers(value bool)

	// DecodeStrings gets decodes quoted strings flag.
	DecodeStrings() bool

	// SetDecodeStrings sets decodes quoted strings flag.
	SetDecodeStrings(value bool)

	// CommentState gets a token state to process comments.
	CommentState() ICommentState

	// NumberState gets a token state to process numbers.
	NumberState() INumberState

	// QuoteState gets a token state to process quoted strings.
	QuoteState() IQuoteState

	// SymbolState gets a token state to process symbols (single like "=" or muti-character like "<>")
	SymbolState() ISymbolState

	// WhitespaceState gets a token state to process white space delimiters.
	WhitespaceState() IWhitespaceState

	// WordState gets a token state to process words or indentificators.
	WordState() IWordState

	// Reader gets the stream scanner to tokenize.
	Reader() io.IScanner

	// SetReader sets the stream scanner to tokenize.
	SetReader(scanner io.IScanner)

	// HasNextToken checks if there is the next token exist.
	//	Returns: <code>true</code> if scanner has the next token.
	HasNextToken() bool

	// NextToken gets the next token from the scanner.
	//	Returns: Next token of <code>null</code> if there are no more tokens left.
	NextToken() *Token

	// TokenizeStream tokenizes a textual stream into a list of token structures.
	//	Parameters:
	//		- scanner: A textual stream to be tokenized.
	//	Returns: A list of token structures.
	TokenizeStream(scanner io.IScanner) []*Token

	// TokenizeBuffer tokenizes a string buffer into a list of tokens structures.
	//	Parameters:
	//		- buffer: A string buffer to be tokenized.
	//	Returns: A list of token structures.
	TokenizeBuffer(buffer string) []*Token

	// TokenizeStreamToStrings tokenizes a textual stream into a list of strings.
	//	Parameters:
	//		- scanner: A textual stream to be tokenized.
	//	Returns: A list of token strings.
	TokenizeStreamToStrings(scanner io.IScanner) []string

	// TokenizeBufferToStrings tokenizes a string buffer into a list of strings.
	//	Parameters:
	//		- buffer: A string buffer to be tokenized.
	//	Returns: A list of token strings.
	TokenizeBufferToStrings(buffer string) []string
}
