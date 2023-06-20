package random

import (
	"math"
	"math/rand"
)

// Long random generator for integer values.
//	Example:
//		value1 := random.Long.Next(5, 10);     // Possible result: 7
//		value2 := random.Long.Update(10, 3);   // Possible result: 9
var Long = &_TRandomLong{}

type _TRandomLong struct{}

// Next generates a integer in the range ['min', 'max']. If 'max' is omitted, then the range
// will be set to [0, 'min'].
//	Parameters:
//		- min: int64 - minimum value of the integer that will be generated. If 'max' is omitted,
//			then 'max' is set to 'min' and 'min' is set to 0.
//		- max: int64 - maximum value of the int that will be generated. Defaults to 'min' if omitted.
//	Returns: generated random int64 value.
func (c *_TRandomLong) Next(min int64, max int64) int64 {
	if max-min <= 0 {
		return min
	}

	return min + rand.Int63n(max-min)
}

// Update updates (drifts) a integer value within specified range defined
//	Parameters:
//		- value: int - a integer value to drift.
//		- interval:int - a range. Default: 10% of the value
//	Returns: int
func (c *_TRandomLong) Update(value int64, interval int64) int64 {
	if interval <= 0 {
		interval = int64(math.Trunc(0.1 * float64(value)))
	}
	minValue := value - interval
	maxValue := value + interval
	return c.Next(minValue, maxValue)
}

// Sequence generates a random sequence of integers starting from 0 like: [0,1,2,3...??]
//	Parameters:
//		- min: int64 - minimum value of the integer that will be generated. If 'max'
//			is omitted, then 'max' is set to 'min' and 'min' is set to 0.
//		- max: int64 - maximum value of the int that will be generated. Defaults to 'min' if omitted.
//	Returns: generated array of int64.
func (c *_TRandomLong) Sequence(min int64, max int64) []int64 {
	if min < 0 {
		min = 0
	}
	if max < min {
		max = min
	}

	count := c.Next(min, max)

	result := make([]int64, count)
	for i := range result {
		result[i] = int64(i)
	}

	return result
}
