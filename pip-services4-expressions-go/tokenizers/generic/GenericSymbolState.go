package generic

import (
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/io"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers"
)

// GenericSymbolState the idea of a symbol is a character that stands on its own, such as an ampersand or a parenthesis.
// For example, when tokenizing the expression <code>(isReady)& (isWilling) </code>, a typical
// tokenizer would return 7 tokens, including one for each parenthesis and one for the ampersand.
// Thus a series of symbols such as <code>)&( </code> becomes three tokens, while a series of letters
// such as <code>isReady</code> becomes a single word token.
// <p/>
// Multi-character symbols are an exception to the rule that a symbol is a standalone character.
// For example, a tokenizer may want less-than-or-equals to tokenize as a single token. This class
// provides a method for establishing which multi-character symbols an object of this class should
// treat as single symbols. This allows, for example, <code>"cat &lt;= dog"</code> to tokenize as
// three tokens, rather than splitting the less-than and equals symbols into separate tokens.
// <p/>
// By default, this state recognizes the following multi-character symbols:
// <code>!=, :-, &lt;=, &gt;=</code>
type GenericSymbolState struct {
	symbols *SymbolRootNode
}

func NewGenericSymbolState() *GenericSymbolState {
	c := &GenericSymbolState{
		symbols: NewSymbolRootNode(),
	}
	return c
}

// NextToken returns a symbol token from a scanner.
//
//	Returns: a symbol token from a scanner.
func (c *GenericSymbolState) NextToken(
	scanner io.IScanner, tokenizer tokenizers.ITokenizer) *tokenizers.Token {
	return c.symbols.NextToken(scanner)
}

// Add a multi-character symbol.
//
//	Parameters:
//		- value: the symbol to add, such as "=:="
func (c *GenericSymbolState) Add(value string, tokenType int) {
	c.symbols.Add(value, tokenType)
}
