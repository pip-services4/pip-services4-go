package test_random

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-data-go/random"
	"github.com/stretchr/testify/assert"
)

func TestArrayPick(t *testing.T) {
	array1 := []any{}
	value1 := random.Array.Pick(array1)
	assert.Nil(t, value1)

	array2 := []any{nil, nil}
	value2 := random.Array.Pick(array2)
	assert.Nil(t, value2)

	array3 := []int{}
	assert.Nil(t, random.Array.Pick(array3))

	array4 := []int{1, 2}
	value4 := random.Array.Pick(array4).(int)
	assert.True(t, value4 == 1 || value4 == 2)
}
