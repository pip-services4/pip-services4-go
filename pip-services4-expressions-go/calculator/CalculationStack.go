package calculator

import "github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/variants"

// CalculationStack implements a stack of Variant values.
type CalculationStack struct {
	values []*variants.Variant
}

func NewCalculationStack() *CalculationStack {
	c := &CalculationStack{
		values: []*variants.Variant{},
	}
	return c
}

func (c *CalculationStack) Length() int {
	return len(c.values)
}

func (c *CalculationStack) Push(value *variants.Variant) {
	c.values = append(c.values, value)
}

func (c *CalculationStack) Pop() *variants.Variant {
	if len(c.values) == 0 {
		panic("Stack is empty.")
	}
	result := c.values[len(c.values)-1]
	c.values = c.values[:len(c.values)-1]
	return result
}

func (c *CalculationStack) PeekAt(index int) *variants.Variant {
	if index < 0 || index >= len(c.values) {
		panic("Stack index is out of bounds.")
	}

	return c.values[index]
}

func (c *CalculationStack) Peek() *variants.Variant {
	if len(c.values) == 0 {
		panic("Stack is empty.")
	}
	return c.values[len(c.values)-1]
}

func (c *CalculationStack) Clear() {
	c.values = []*variants.Variant{}
}
