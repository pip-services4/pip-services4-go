package tokenizers

import (
	"strings"

	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/io"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers/utilities"
)

// MustacheSpecialState implements a quote string state object for Mustache templates.
type MustacheSpecialState struct {
}

func NewMustacheSpecialState() *MustacheSpecialState {
	c := &MustacheSpecialState{}
	return c
}

// NextToken gets the next token from the stream started from the character linked to this state.
//
//	Parameters:
//		- scanner: A textual string to be tokenized.
//		- tokenizer: A tokenizer class that controls the process.
//	Returns: The next token from the top of the stream.
func (c *MustacheSpecialState) NextToken(scanner io.IScanner,
	tokenizer tokenizers.ITokenizer) *tokenizers.Token {

	line := scanner.PeekLine()
	column := scanner.PeekColumn()
	tokenValue := strings.Builder{}
	nextSymbol := scanner.Read()

	for !utilities.CharValidator.IsEof(nextSymbol) {
		if nextSymbol == '{' {
			chr := scanner.Peek()
			if chr == '{' {
				scanner.Unread()
				break
			}
		}

		tokenValue.WriteRune(nextSymbol)
		nextSymbol = scanner.Read()
	}

	return tokenizers.NewToken(tokenizers.Special, tokenValue.String(), line, column)
}
