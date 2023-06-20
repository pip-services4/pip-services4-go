package random

import (
	"math"
	"math/rand"
)

// Integer Random generator for integer values.
//	Example:
//		value1 := random.Integer.Next(5, 10);     // Possible result: 7
//		value2 := random.Integer.Update(10, 3);   // Possible result: 9
var Integer = &_TRandomInteger{}

type _TRandomInteger struct{}

// Next generates a integer in the range ['min', 'max']. If 'max' is omitted, then the range
// will be set to [0, 'min'].
//	Parameters:
//		- min: int - minimum value of the integer that will be generated. If 'max' is omitted,
//			then 'max' is set to 'min' and 'min' is set to 0.
//		- max: int - maximum value of the int that will be generated. Defaults to 'min' if omitted.
//	Returns: generated random integer value.
func (c *_TRandomInteger) Next(min int, max int) int {
	if max-min <= 0 {
		return min
	}
	return min + rand.Intn(max-min)
}

// Update updates (drifts) a integer value within specified range defined
//	Parameters:
//		- value: int - a integer value to drift.
//		- interval:int - a range. Default: 10% of the value
//	Returns: int
func (c *_TRandomInteger) Update(value int, interval int) int {
	if interval <= 0 {
		interval = int(math.Trunc(0.1 * float64(value)))
	}
	minValue := value - interval
	maxValue := value + interval
	return c.Next(minValue, maxValue)
}

// Sequence generates a random sequence of integers starting from 0 like: [0,1,2,3...??]
//	Parameters:
//		- min: int - minimum value of the integer that will be generated. If 'max'
//			is omitted, then 'max' is set to 'min' and 'min' is set to 0.
//		- max: int - maximum value of the int that will be generated. Defaults to 'min' if omitted.
//	Returns: generated array of integers.
func (c *_TRandomInteger) Sequence(min int, max int) []int {
	if min < 0 {
		min = 0
	}
	if max < min {
		max = min
	}

	count := c.Next(min, max)

	result := make([]int, count)
	for i := range result {
		result[i] = i
	}

	return result
}
