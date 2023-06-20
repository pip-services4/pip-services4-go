package generic

import (
	"strings"

	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/io"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers/utilities"
)

// GenericNumberState a NumberState object returns a number from a scanner. This state's idea of a number allows
// an optional, initial minus sign, followed by one or more digits. A decimal point and another string
// of digits may follow these digits.
type GenericNumberState struct{}

func NewGenericNumberState() *GenericNumberState {
	c := &GenericNumberState{}
	return c
}

// NextToken gets the next token from the stream started from the character linked to this state.
//
//	Parameters:
//		- scanner: A textual string to be tokenized.
//		- tokenizer: A tokenizer class that controls the process.
//	Returns: The next token from the top of the stream.
func (c *GenericNumberState) NextToken(
	scanner io.IScanner, tokenizer tokenizers.ITokenizer) *tokenizers.Token {

	absorbedDot := false
	gotADigit := false
	tokenValue := strings.Builder{}
	nextSymbol := scanner.Read()
	line := scanner.Line()
	column := scanner.Column()

	// Parses leading minus.
	if nextSymbol == '-' {
		tokenValue.WriteRune('-')
		nextSymbol = scanner.Read()
	}

	// Parses digits before decimal separator.
	for utilities.CharValidator.IsDigit(nextSymbol) &&
		!utilities.CharValidator.IsEof(nextSymbol) {
		gotADigit = true
		tokenValue.WriteRune(nextSymbol)
		nextSymbol = scanner.Read()
	}

	// Parses part after the decimal separator.
	if nextSymbol == '.' {
		absorbedDot = true
		tokenValue.WriteRune('.')
		nextSymbol = scanner.Read()

		// Absorb all digits.
		for utilities.CharValidator.IsDigit(nextSymbol) &&
			!utilities.CharValidator.IsEof(nextSymbol) {
			gotADigit = true
			tokenValue.WriteRune(nextSymbol)
			nextSymbol = scanner.Read()
		}
	}

	// Unread last unprocessed symbol.
	if !utilities.CharValidator.IsEof(nextSymbol) {
		scanner.Unread()
	}

	// Process the result.
	if !gotADigit {
		scanner.UnreadMany(tokenValue.Len())
		if tokenizer != nil && tokenizer.SymbolState() != nil {
			return tokenizer.SymbolState().NextToken(scanner, tokenizer)
		} else {
			panic("Tokenizer must have an assigned symbol state.")
		}
	}

	tokenType := tokenizers.Integer
	if absorbedDot {
		tokenType = tokenizers.Float
	}
	return tokenizers.NewToken(tokenType, tokenValue.String(), line, column)
}
