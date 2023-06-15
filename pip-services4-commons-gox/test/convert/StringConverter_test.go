package test_convert

import (
	"testing"

	"github.com/pip-services4/pip-services4-commons-go/convert"
	"github.com/stretchr/testify/assert"
)

func TestToString(t *testing.T) {
	str, ok := convert.StringConverter.ToNullableString(nil)
	assert.False(t, ok)
	assert.Equal(t, "", str)

	assert.Equal(t, "xyz", convert.StringConverter.ToString("xyz"))
	assert.Equal(t, "16030862614303175036", convert.StringConverter.ToString((uint64)(16030862614303175036)))
	assert.Equal(t, "123", convert.StringConverter.ToString(123))
	assert.Equal(t, "true", convert.StringConverter.ToString(true))

	value := struct{ prop string }{"xyz"}
	assert.Equal(t, "{xyz}", convert.StringConverter.ToString(value))

	array1 := []string{"A", "B", "C"}
	assert.Equal(t, "A,B,C", convert.StringConverter.ToString(array1))

	array2 := []int32{1, 2, 3}
	assert.Equal(t, "1,2,3", convert.StringConverter.ToString(array2))

	assert.Equal(t, "xyz", convert.StringConverter.ToStringWithDefault(nil, "xyz"))
}
