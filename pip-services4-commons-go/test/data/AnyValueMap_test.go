package test_data

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	"github.com/stretchr/testify/assert"
)

func TestAnyValueMapNew(t *testing.T) {
	mp := data.NewEmptyAnyValueMap()
	_, ok := mp.GetAsObject("key1")
	assert.False(t, ok)

	mp = data.NewAnyValueMapFromValue(map[string]any{
		"key1": 1,
		"key2": "A",
	})
	val, ok := mp.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, int64(1), val)

	val, ok = mp.Get("key2")
	assert.True(t, ok)
	assert.Equal(t, "A", val)

	val, ok = mp.GetAsNullableString("key3")
	assert.False(t, ok)
	assert.Equal(t, "", val)

	mp = data.NewAnyValueMapFromMaps(map[string]any{
		"key1": 1,
		"key2": "A",
	})
	val, ok = mp.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, 1, val)

	val, ok = mp.Get("key2")
	assert.True(t, ok)
	assert.Equal(t, "A", val)

	mp = data.NewAnyValueMapFromTuples(
		"key1", 1,
		"key2", "A",
	)
	val, ok = mp.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, 1, val)

	val, ok = mp.Get("key2")
	assert.True(t, ok)
	assert.Equal(t, "A", val)
}

func TestAnyValueMapGetAndSet(t *testing.T) {
	mp := data.NewEmptyAnyValueMap()
	_, ok := mp.GetAsObject("key1")
	assert.False(t, ok)

	mp.SetAsObject("key1", 1)
	assert.Equal(t, 1, mp.GetAsInteger("key1"))
	assert.True(t, 1.0-mp.GetAsFloat("key1") < 0.001)
	assert.Equal(t, "1", mp.GetAsString("key1"))

	mp.Put("key2", "1")
	assert.Equal(t, 1, mp.GetAsInteger("key2"))
	assert.True(t, 1.0-mp.GetAsFloat("key2") < 0.001)
	assert.Equal(t, "1", mp.GetAsString("key2"))

	mp.Remove("key2")
	_, ok = mp.GetAsObject("key2")
	assert.False(t, ok)
}

func TestAnyValueMapGetAsNullable(t *testing.T) {
	array := data.NewEmptyAnyValueMap()
	assert.Equal(t, 0, array.Len())

	_, ok := array.GetAsNullableInteger("")
	assert.False(t, ok)

	_, ok = array.GetAsNullableLong("")
	assert.False(t, ok)

	_, ok = array.GetAsNullableFloat("")
	assert.False(t, ok)

	_, ok = array.GetAsNullableDouble("")
	assert.False(t, ok)

	_, ok = array.GetAsNullableBoolean("")
	assert.False(t, ok)

	_, ok = array.GetAsNullableString("")
	assert.False(t, ok)

	_, ok = array.GetAsNullableDuration("")
	assert.False(t, ok)

	_, ok = array.GetAsNullableDateTime("")
	assert.False(t, ok)
}
