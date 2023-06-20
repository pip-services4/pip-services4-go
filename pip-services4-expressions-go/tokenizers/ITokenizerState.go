package tokenizers

import (
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/io"
)

// ITokenizerState a tokenizerState returns a token, given a scanner, an initial character read from the scanner,
// and a tokenizer that is conducting an overall tokenization of the scanner. The tokenizer will
// typically have a character state table that decides which state to use, depending on an initial
// character. If a single character is insufficient, a state such as <code>SlashState</code>
// will read a second character, and may delegate to another state, such as <code>SlashStarState</code>.
// This prospect of delegation is the reason that the <code>nextToken()</code>
// method has a tokenizer argument.
type ITokenizerState interface {
	// NextToken gets the next token from the stream started from the character linked to this state.
	//	Parameters:
	//		- scanner: A textual string to be tokenized.
	//		- tokenizer: A tokenizer class that controls the process.
	//	Returns: The next token from the top of the stream.
	NextToken(scanner io.IScanner, tokenizer ITokenizer) *Token
}
