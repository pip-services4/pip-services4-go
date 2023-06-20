package tokenizers

import (
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers/generic"
)

// ExpressionSymbolState implements a symbol state object.
type ExpressionSymbolState struct {
	*generic.GenericSymbolState
}

// NewExpressionSymbolState constructs an instance of this class.
func NewExpressionSymbolState() *ExpressionSymbolState {
	c := &ExpressionSymbolState{
		GenericSymbolState: generic.NewGenericSymbolState(),
	}

	c.Add("<=", tokenizers.Symbol)
	c.Add(">=", tokenizers.Symbol)
	c.Add("<>", tokenizers.Symbol)
	c.Add("!=", tokenizers.Symbol)
	c.Add(">>", tokenizers.Symbol)
	c.Add("<<", tokenizers.Symbol)

	return c
}
