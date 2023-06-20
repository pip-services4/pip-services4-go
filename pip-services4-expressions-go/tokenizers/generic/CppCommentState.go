package generic

import (
	"strings"

	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/io"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers/utilities"
)

// CppCommentState this state will either delegate to a comment-handling state,
// or return a token with just a slash in it.
type CppCommentState struct{}

func NewCppCommentState() *CppCommentState {
	c := &CppCommentState{}
	return c
}

// GetMultiLineComment ignore everything up to a closing star and slash, and then return the tokenizer's next token.
func (c *CppCommentState) GetMultiLineComment(scanner io.IScanner) string {
	result := strings.Builder{}

	lastSymbol := rune(0)
	nextSymbol := scanner.Read()
	for !utilities.CharValidator.IsEof(nextSymbol) {
		result.WriteRune(nextSymbol)
		if lastSymbol == '*' && nextSymbol == '/' {
			break
		}
		lastSymbol = nextSymbol

		nextSymbol = scanner.Read()
	}

	return result.String()
}

// GetSingleLineComment ignore everything up to an end-of-line and return the tokenizer's next token.
func (c *CppCommentState) GetSingleLineComment(scanner io.IScanner) string {

	result := strings.Builder{}

	nextSymbol := scanner.Read()
	for !utilities.CharValidator.IsEof(nextSymbol) && !utilities.CharValidator.IsEol(nextSymbol) {
		result.WriteRune(nextSymbol)

		nextSymbol = scanner.Read()
	}

	if utilities.CharValidator.IsEol(nextSymbol) {
		scanner.Unread()
	}

	return result.String()
}

// NextToken either delegate to a comment-handling state, or return a token with just a slash in it.
//
//	Returns: either just a slash token, or the results of delegating to a comment-handling state.
func (c *CppCommentState) NextToken(
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
	} else if secondSymbol == '/' {
		str := c.GetSingleLineComment(scanner)
		return tokenizers.NewToken(tokenizers.Comment, "//"+str, line, column)
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
