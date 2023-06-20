package csv

import "github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers/generic"

// CsvWordState implements a word state to tokenize CSV stream.
type CsvWordState struct {
	*generic.GenericWordState
}

// NewCsvWordState constructs this object with specified parameters.
//	Parameters:
//		- fieldSeparators: Separators for fields in CSV stream.
//		- quoteSymbol: Delimiters character to quote strings.
func NewCsvWordState(fieldSeparators []rune, quoteSymbols []rune) *CsvWordState {
	c := &CsvWordState{
		GenericWordState: generic.NewGenericWordState(),
	}

	c.ClearWordChars()
	c.SetWordChars(0x0000, 0xffff, true)

	c.SetWordChars(CR, CR, false)
	c.SetWordChars(LF, LF, false)

	for _, fieldSeparator := range fieldSeparators {
		c.SetWordChars(fieldSeparator, fieldSeparator, false)
	}

	for _, quoteSymbol := range quoteSymbols {
		c.SetWordChars(quoteSymbol, quoteSymbol, false)
	}

	return c
}
