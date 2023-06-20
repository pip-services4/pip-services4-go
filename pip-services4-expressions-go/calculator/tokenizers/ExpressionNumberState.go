package tokenizers

import (
	"strings"

	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/io"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers/generic"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers/utilities"
)

// ExpressionNumberState implements an Expression-specific number state object.
type ExpressionNumberState struct {
	*generic.GenericNumberState
}

func NewExpressionNumberState() *ExpressionNumberState {
	c := &ExpressionNumberState{
		GenericNumberState: generic.NewGenericNumberState(),
	}
	return c
}

// NextToken gets the next token from the stream started from the character linked to this state.
//
//	Parameters:
//		- scanner: A textual string to be tokenized.
//		- tokenizer: A tokenizer class that controls the process.
//	Returns: The next token from the top of the stream.
func (c *ExpressionNumberState) NextToken(scanner io.IScanner,
	tokenizer tokenizers.ITokenizer) *tokenizers.Token {
	nextChar := scanner.Peek()
	line := scanner.PeekLine()
	column := scanner.PeekColumn()

	// Process leading minus.
	if nextChar == '-' {
		return tokenizer.SymbolState().NextToken(scanner, tokenizer)
	}

	// Process numbers using base class algorithm.
	token := c.GenericNumberState.NextToken(scanner, tokenizer)

	// Exit if number was not detected.
	if token.Type() != tokenizers.Integer && token.Type() != tokenizers.Float {
		return token
	}

	// Exit if number is not in scientific format.
	nextChar = scanner.Peek()

	if nextChar != 'e' && nextChar != 'E' {
		return token
	}

	nextChar = scanner.Read()
	tokenValue := strings.Builder{}
	tokenValue.WriteRune(nextChar)

	// Process '-' or '+' in mantissa
	nextChar = scanner.Peek()

	if nextChar == '-' || nextChar == '+' {
		nextChar = scanner.Read()
		tokenValue.WriteRune(nextChar)
		nextChar = scanner.Peek()
	}

	// Exit if mantissa has no digits.
	if !utilities.CharValidator.IsDigit(nextChar) {
		scanner.UnreadMany(tokenValue.Len())
		return token
	}

	// Process matissa digits
	for utilities.CharValidator.IsDigit(nextChar) {
		nextChar = scanner.Read()
		tokenValue.WriteRune(nextChar)
		nextChar = scanner.Peek()
	}

	return tokenizers.NewToken(tokenizers.Float, token.Value()+tokenValue.String(), line, column)
}
