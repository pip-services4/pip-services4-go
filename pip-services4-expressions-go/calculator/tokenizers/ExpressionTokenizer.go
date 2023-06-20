package tokenizers

import (
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers/generic"
)

// ExpressionTokenizer implements a tokenizer to perform lexical analysis for expressions.
type ExpressionTokenizer struct {
	*tokenizers.AbstractTokenizer
}

func NewExpressionTokenizer() *ExpressionTokenizer {
	c := &ExpressionTokenizer{}
	c.AbstractTokenizer = tokenizers.InheritAbstractTokenizer(c)

	c.SetDecodeStrings(false)

	c.SetWhitespaceState(generic.NewGenericWhitespaceState())
	c.SetSymbolState(NewExpressionSymbolState())
	c.SetNumberState(NewExpressionNumberState())
	c.SetQuoteState(NewExpressionQuoteState())
	c.SetWordState(NewExpressionWordState())
	c.SetCommentState(generic.NewCCommentState())

	c.ClearCharacterStates()
	c.SetCharacterState(0x0000, 0xffff, c.SymbolState())
	c.SetCharacterState(0, ' ', c.WhitespaceState())

	c.SetCharacterState('a', 'z', c.WordState())
	c.SetCharacterState('A', 'Z', c.WordState())
	c.SetCharacterState(0x00c0, 0x00ff, c.WordState())
	c.SetCharacterState('_', '_', c.WordState())

	c.SetCharacterState('0', '9', c.NumberState())
	c.SetCharacterState('-', '-', c.NumberState())
	c.SetCharacterState('.', '.', c.NumberState())

	c.SetCharacterState('"', '"', c.QuoteState())
	c.SetCharacterState('\'', '\'', c.QuoteState())

	c.SetCharacterState('/', '/', c.CommentState())

	return c
}
