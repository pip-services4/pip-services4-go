package test_data

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	"github.com/stretchr/testify/assert"
)

func TestAnyValueArrayCreate(t *testing.T) {
	array := data.NewEmptyAnyValueArray()
	assert.Equal(t, 0, array.Len())

	array = data.NewAnyValueArray([]any{1, 2, 3})
	assert.Equal(t, 3, array.Len())
	assert.Equal(t, "1,2,3", array.String())

	array = data.NewAnyValueArrayFromString("Fatal,Error,Info,", ",", true)
	assert.Equal(t, 3, array.Len())

	array = data.NewAnyValueArray([]any{1, 2, 3})
	assert.Equal(t, 3, array.Len())
	assert.True(t, array.Contains(1))

	array = data.NewAnyValueArrayFromValue([]any{1, 2, 3})
	assert.Equal(t, 3, array.Len())
	val, ok := array.Get(0)
	assert.True(t, ok)
	assert.Equal(t, int64(1), val)

	val, ok = array.GetAsNullableLong(4)
	assert.False(t, ok)
	assert.Equal(t, int64(0), val)
}

func TestAnyValueArrayGetAsNullable(t *testing.T) {
	array := data.NewEmptyAnyValueArray()
	assert.Equal(t, 0, array.Len())

	_, ok := array.GetAsNullableInteger(0)
	assert.False(t, ok)

	_, ok = array.GetAsNullableLong(0)
	assert.False(t, ok)

	_, ok = array.GetAsNullableFloat(0)
	assert.False(t, ok)

	_, ok = array.GetAsNullableDouble(0)
	assert.False(t, ok)

	_, ok = array.GetAsNullableBoolean(0)
	assert.False(t, ok)

	_, ok = array.GetAsNullableString(0)
	assert.False(t, ok)

	_, ok = array.GetAsNullableDuration(0)
	assert.False(t, ok)

	_, ok = array.GetAsNullableDateTime(0)
	assert.False(t, ok)
}
