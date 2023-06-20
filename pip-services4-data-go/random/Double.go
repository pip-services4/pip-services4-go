package random

import (
	"math/rand"
)

// Double random generator for double values.
//	Example:
//		value1 := random.Double.Next(5, 10);     // Possible result: 7.3
//		value2 := random.Double.Update(10, 3);   // Possible result: 9.2
var Double = &_TRandomDouble{}

type _TRandomDouble struct{}

// Next generates a random double value in the range ['min', 'max'].
//	Parameters:
//		- min: float64 - minimum range value
//		- max: float64 - max range value
//	Returns: float64 - a random value.
func (c *_TRandomDouble) Next(min float64, max float64) float64 {
	if max-min <= 0 {
		return min
	}

	return min + rand.Float64()*(max-min)
}

// Update updates (drifts) a double value within specified range defined
//	Parameters:
//		- value: float64 - value to drift.
//		- interval: float64 - a range to drift. Default: 10% of the value
//	Returns: float64
func (c *_TRandomDouble) Update(value float64, interval float64) float64 {
	if interval <= 0 {
		interval = 0.1 * value
	}
	minValue := value - interval
	maxValue := value + interval
	return c.Next(minValue, maxValue)
}
