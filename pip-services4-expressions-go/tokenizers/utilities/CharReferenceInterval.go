package utilities

// CharReferenceInterval represents a character interval that keeps a reference.
// This class is internal and used by CharacterReferenceMap.
type CharReferenceInterval struct {
	start     rune
	end       rune
	reference any
}

func NewCharReferenceInterval(start rune, end rune, reference any) *CharReferenceInterval {
	if start > end {
		panic("Start must be less or equal End")
	}
	return &CharReferenceInterval{
		start:     start,
		end:       end,
		reference: reference,
	}
}

func (c *CharReferenceInterval) Start() rune {
	return c.start
}

func (c *CharReferenceInterval) End() rune {
	return c.end
}

func (c *CharReferenceInterval) Reference() any {
	return c.reference
}

func (c *CharReferenceInterval) InRange(symbol rune) bool {
	return symbol >= c.start && symbol <= c.end
}
