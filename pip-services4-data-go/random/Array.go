package random

import (
	"reflect"
)

// Array random generator for array objects.
//	Examples:
//		value1 := random.Array.Pick([]int{1, 2, 3, 4}) // Possible result: 3
var Array = &_TRandomArray{}

type _TRandomArray struct{}

// Pick picks a random element from specified array.
//	Parameters:
//		- values: an array of any interface
//	Returns: a randomly picked item.
func (c *_TRandomArray) Pick(value any) any {
	if value == nil {
		return nil
	}

	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Array && v.Kind() != reflect.Slice {
		return nil
	}

	ln := v.Len()
	if ln == 0 {
		return nil
	}

	index := Integer.Next(0, ln-1)

	v = v.Index(index)
	return v.Interface()
}
