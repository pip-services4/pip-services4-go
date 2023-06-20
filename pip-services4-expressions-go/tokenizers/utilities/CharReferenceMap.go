package utilities

// CharReferenceMap this class keeps references associated with specific characters
type CharReferenceMap struct {
	initialInterval []any
	otherIntervals  []*CharReferenceInterval
}

func NewCharReferenceMap() *CharReferenceMap {
	c := &CharReferenceMap{}
	c.Clear()
	return c
}

func (c *CharReferenceMap) AddDefaultInterval(reference any) {
	c.AddInterval(0x0000, 0xfffe, reference)
}

func (c *CharReferenceMap) AddInterval(start rune, end rune, reference any) {
	if start > end {
		panic("Start must be less or equal End")
	}
	if end >= 0xffff {
		end = 0xfffe
	}

	for index := int(start); index < 0x0100 && index <= int(end); index++ {
		c.initialInterval[index] = reference
	}
	if end >= 0x0100 {
		if start < 0x0100 {
			start = 0x100
		}
		c.otherIntervals = append(
			[]*CharReferenceInterval{NewCharReferenceInterval(start, end, reference)},
			c.otherIntervals...)
	}
}

func (c *CharReferenceMap) Clear() {
	c.initialInterval = make([]any, 0x0100)
	for index := 0x0000; index < 0x0100; index++ {
		c.initialInterval[index] = nil
	}
	c.otherIntervals = []*CharReferenceInterval{}
}

func (c *CharReferenceMap) Lookup(symbol rune) any {
	if symbol < 0 {
		return nil
	} else if symbol < 0x0100 {
		return c.initialInterval[int(symbol)]
	} else {
		for _, interval := range c.otherIntervals {
			if interval.InRange(symbol) {
				return interval.Reference
			}
		}
		return nil
	}
}
