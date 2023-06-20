package random

import (
	"math/rand"
)

// Boolean Random generator for boolean values.
//	Example:
//		value1 := random.Boolean.Next();      // Possible result: true
//		value2 := random.Boolean.Chance(1,3); // Possible result: false
var Boolean = &_TRandomBoolean{}

type _TRandomBoolean struct{}

// Chance calculates "chance" out of "max chances". Example: 1 chance out of 3 chances (or 33.3%)
//	Parameters:
//		- chance: number  - a chance proportional to maxChances.
//		- maxChances: number - a maximum number of chances
//	Returns: bool
func (c *_TRandomBoolean) Chance(chances int, maxChances int) bool {
	if chances < 0 {
		chances = 0
	}
	if maxChances < 0 {
		maxChances = 0
	}
	if chances == 0 && maxChances == 0 {
		return false
	}
	if maxChances < chances {
		maxChances = chances
	}
	start := (maxChances - chances) / 2
	end := start + chances
	hit := rand.Intn(maxChances + 1)
	return hit >= start && hit <= end
}

// Next generates a random boolean value.
//	Returns: bool
func (c *_TRandomBoolean) Next() bool {
	return rand.Float32()*100 < 50
}
