package test_convert

import (
	"testing"

	"github.com/pip-services4/pip-services4-commons-go/convert"
	"github.com/stretchr/testify/assert"
)

func TestObjectToMap(t *testing.T) {
	mp, ok := convert.MapConverter.ToNullableMap(nil)
	assert.False(t, ok)
	assert.Nil(t, mp)

	v1 := struct{ value1, value2 float64 }{123, 234}
	mp = convert.MapConverter.ToMap(v1)
	assert.Len(t, mp, 2)
	assert.Equal(t, 123., mp["value1"])
	assert.Equal(t, 234., mp["value2"])

	v2 := map[string]interface{}{"value1": 123}
	mp = convert.MapConverter.ToMap(v2)
	assert.Len(t, mp, 1)
	assert.Equal(t, int64(123), mp["value1"])
}

func TestToNullableMap(t *testing.T) {
	mp, ok := convert.MapConverter.ToNullableMap(nil)
	assert.False(t, ok)
	assert.Nil(t, mp)
	mp, ok = convert.MapConverter.ToNullableMap(5)
	assert.False(t, ok)
	assert.Nil(t, mp)

	array := []int{1, 2}

	mp, ok = convert.MapConverter.ToNullableMap(array)
	assert.True(t, ok)
	assert.Len(t, mp, 2)
	assert.Equal(t, int64(1), mp["0"])
	assert.Equal(t, int64(2), mp["1"])

	values := []string{"ab", "cd"}
	mp, ok = convert.MapConverter.ToNullableMap(values)
	assert.True(t, ok)
	assert.Len(t, mp, 2)
	assert.Equal(t, "ab", mp["0"])
	assert.Equal(t, "cd", mp["1"])

	hash := map[int]string{}
	hash[8] = "title 8"
	hash[11] = "title 11"
	mp, ok = convert.MapConverter.ToNullableMap(hash)
	assert.True(t, ok)
	assert.Len(t, mp, 2)
	assert.Equal(t, "title 8", mp["8"])
	assert.Equal(t, "title 11", mp["11"])
}
