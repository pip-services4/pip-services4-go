package test_convert

import (
	"testing"

	"github.com/pip-services4/pip-services4-commons-go/convert"
	"github.com/stretchr/testify/assert"
)

func TestToNullableArray(t *testing.T) {
	arr, ok := convert.ArrayConverter.ToNullableArray(nil)
	assert.False(t, ok)
	assert.Nil(t, arr)

	arr, ok = convert.ArrayConverter.ToNullableArray(2)
	assert.True(t, ok)
	assert.Len(t, arr, 1)
	assert.Equal(t, int64(2), arr[0])

	array := []int{1, 2}
	arr, ok = convert.ArrayConverter.ToNullableArray(array)
	assert.True(t, ok)
	assert.Len(t, arr, 2)
	assert.Equal(t, int64(1), arr[0])
	assert.Equal(t, int64(2), arr[1])

	stringArray := []string{"ab", "cd"}
	arr, ok = convert.ArrayConverter.ToNullableArray(stringArray)
	assert.True(t, ok)
	assert.Len(t, arr, 2)
	assert.Equal(t, "ab", arr[0])
	assert.Equal(t, "cd", arr[1])
}

func TestToArray(t *testing.T) {
	arr := convert.ArrayConverter.ToArray(nil)
	assert.Len(t, arr, 0)

	arr = convert.ArrayConverter.ToArray(2)
	assert.Len(t, arr, 1)
	assert.Equal(t, int64(2), arr[0])

	array := []int{1, 2}
	arr = convert.ArrayConverter.ToArray(array)
	assert.Len(t, arr, 2)
	assert.Equal(t, int64(1), arr[0])
	assert.Equal(t, int64(2), arr[1])

	stringArray := []string{"ab", "cd"}
	arr = convert.ArrayConverter.ToArray(stringArray)
	assert.Len(t, arr, 2)
	assert.Equal(t, "ab", arr[0])
	assert.Equal(t, "cd", arr[1])
}
