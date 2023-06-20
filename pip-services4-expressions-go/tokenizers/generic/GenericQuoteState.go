package generic

import (
	"strings"

	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/io"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers/utilities"
)

// GenericQuoteState a quoteState returns a quoted string token from a scanner. This state will collect characters
// until it sees a match to the character that the tokenizer used to switch to this state.
// For example, if a tokenizer uses a double-quote character to enter this state,
// then <code>nextToken()</code> will search for another double-quote until it finds one
// or finds the end of the scanner.
type GenericQuoteState struct{}

func NewGenericQuoteState() *GenericQuoteState {
	c := &GenericQuoteState{}
	return c
}

// NextToken return a quoted string token from a scanner. This method will collect
// characters until it sees a match to the character that the tokenizer used
// to switch to this state.
//
//	Returns: a quoted string token from a scanner.
func (c *GenericQuoteState) NextToken(
	scanner io.IScanner, tokenizer tokenizers.ITokenizer) *tokenizers.Token {

	firstSymbol := scanner.Read()
	line := scanner.Line()
	column := scanner.Column()
	tokenValue := strings.Builder{}
	tokenValue.WriteRune(firstSymbol)

	nextSymbol := scanner.Read()
	for !utilities.CharValidator.IsEof(nextSymbol) {
		tokenValue.WriteRune(nextSymbol)
		if nextSymbol == firstSymbol {
			break
		}

		nextSymbol = scanner.Read()
	}

	return tokenizers.NewToken(tokenizers.Quoted, tokenValue.String(), line, column)
}

// EncodeString encodes a string value.
//
//	Parameters:
//		- value: A string value to be encoded.
//		- quoteSymbol: A string quote character.
//	Returns: An encoded string.
func (c *GenericQuoteState) EncodeString(value string, quoteSymbol rune) string {
	result := strings.Builder{}
	result.WriteRune(quoteSymbol)
	result.WriteString(value)
	result.WriteRune(quoteSymbol)
	return result.String()
}

// DecodeString decodes a string value.
//
//	Parameters:
//		- value: A string value to be decoded.
//		- quoteSymbol: A string quote character.
//	Returns>: An decoded string.
func (c *GenericQuoteState) DecodeString(value string, quoteSymbol rune) string {
	runes := []rune(value)
	if len(runes) >= 2 && runes[0] == quoteSymbol && runes[len(runes)-1] == quoteSymbol {
		return string(runes[1 : len(runes)-1])
	}
	return value
}
