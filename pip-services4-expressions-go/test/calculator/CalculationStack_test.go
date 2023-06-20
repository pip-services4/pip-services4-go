package test_calculator

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/calculator"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/variants"
	"github.com/stretchr/testify/assert"
)

func TestCalculationStackOperations(t *testing.T) {
	stack := calculator.NewCalculationStack()

	stack.Push(variants.VariantFromInteger(1))
	assert.Equal(t, 1, stack.Length())

	stack.Push(variants.VariantFromInteger(2))
	assert.Equal(t, 2, stack.Length())

	v := stack.Peek()
	assert.Equal(t, 2, v.AsInteger())
	v = stack.PeekAt(0)
	assert.Equal(t, 1, v.AsInteger())

	v = stack.Pop()
	assert.Equal(t, 2, v.AsInteger())
	assert.Equal(t, 1, stack.Length())

	v = stack.Pop()
	assert.Equal(t, 1, v.AsInteger())
	assert.Equal(t, 0, stack.Length())
}
