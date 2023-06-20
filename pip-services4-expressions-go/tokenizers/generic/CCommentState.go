package generic

import (
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/io"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers/utilities"
)

// CCommentState this state will either delegate to a comment-handling state,
// or return a token with just a slash in it.
type CCommentState struct {
	*CppCommentState
}

func NewCCommentState() *CCommentState {
	c := &CCommentState{
		CppCommentState: NewCppCommentState(),
	}
	return c
}

// NextToken either delegate to a comment-handling state, or return a token with just a slash in it.
//
//	Returns: Either just a slash token, or the results of delegating to a comment-handling state.
func (c *CCommentState) NextToken(
	scanner io.IScanner, tokenizer tokenizers.ITokenizer) *tokenizers.Token {

	firstSymbol := scanner.Read()
	line := scanner.Line()
	column := scanner.Column()

	if firstSymbol != '/' {
		scanner.Unread()
		panic("Incorrect usage of CppCommentState.")
	}

	secondSymbol := scanner.Read()
	if secondSymbol == '*' {
		str := c.GetMultiLineComment(scanner)
		return tokenizers.NewToken(tokenizers.Comment, "/*"+str, line, column)
	} else {
		if !utilities.CharValidator.IsEof(secondSymbol) {
			scanner.Unread()
		}
		if !utilities.CharValidator.IsEof(firstSymbol) {
			scanner.Unread()
		}
		return tokenizer.SymbolState().NextToken(scanner, tokenizer)
	}
}
