package csv

import (
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/io"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers/generic"
)

// CsvSymbolState implements a symbol state to tokenize delimiters in CSV streams.
type CsvSymbolState struct {
	*generic.GenericSymbolState
}

// NewCsvSymbolState constructs this object with specified parameters.
func NewCsvSymbolState() *CsvSymbolState {
	c := &CsvSymbolState{
		GenericSymbolState: generic.NewGenericSymbolState(),
	}

	c.Add("\n", tokenizers.Eol)
	c.Add("\r", tokenizers.Eol)
	c.Add("\r\n", tokenizers.Eol)
	c.Add("\n\r", tokenizers.Eol)

	return c
}

func (c *CsvSymbolState) NextToken(
	scanner io.IScanner, tokenizer tokenizers.ITokenizer) *tokenizers.Token {

	// Optimization...
	nextSymbol := scanner.Read()
	line := scanner.Line()
	column := scanner.Column()

	if nextSymbol != '\n' && nextSymbol != '\r' {
		return tokenizers.NewToken(tokenizers.Symbol, string(nextSymbol), line, column)
	} else {
		scanner.Unread()
		return c.GenericSymbolState.NextToken(scanner, tokenizer)
	}
}
