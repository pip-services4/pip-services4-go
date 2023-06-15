package test_reflect

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/reflect"
	"github.com/stretchr/testify/assert"
)

func TestReaderHasProperty(t *testing.T) {
	obj := NewTestClass()
	assert.True(t, reflect.ObjectReader.HasProperty(obj, "RootPublicProperty"))
	assert.True(t, reflect.ObjectReader.HasProperty(obj, "PublicField"))
	assert.False(t, reflect.ObjectReader.HasProperty(obj, "privateField"))

	dict := map[string]any{
		"1": "AAA",
		"A": 111,
	}
	assert.True(t, reflect.ObjectReader.HasProperty(dict, "1"))
	assert.True(t, reflect.ObjectReader.HasProperty(dict, "A"))
	assert.False(t, reflect.ObjectReader.HasProperty(dict, "B"))

	list := []any{"BBB", 222}
	assert.True(t, reflect.ObjectReader.HasProperty(list, "0"))
	assert.True(t, reflect.ObjectReader.HasProperty(list, "1"))
	assert.False(t, reflect.ObjectReader.HasProperty(list, "3"))
}

func TestReaderGetProperty(t *testing.T) {
	obj := NewTestClass()
	assert.Equal(t, true, reflect.ObjectReader.GetProperty(obj, "RootPublicProperty"))
	assert.Equal(t, "BBB", reflect.ObjectReader.GetProperty(obj, "PublicField"))
	assert.Nil(t, reflect.ObjectReader.GetProperty(obj, "privateField"))

	dict := map[string]any{
		"1": "AAA",
		"A": 111,
	}
	assert.Equal(t, "AAA", reflect.ObjectReader.GetProperty(dict, "1"))
	assert.Equal(t, 111, reflect.ObjectReader.GetProperty(dict, "A"))
	assert.Nil(t, reflect.ObjectReader.GetProperty(dict, "B"))

	list := []any{"BBB", 222}
	assert.Equal(t, "BBB", reflect.ObjectReader.GetProperty(list, "0"))
	assert.Equal(t, 222, reflect.ObjectReader.GetProperty(list, "1"))
	assert.Nil(t, reflect.ObjectReader.GetProperty(list, "3"))
}

func TestReaderGetPropertyNames(t *testing.T) {
	obj := NewTestClass()
	assert.Equal(t, 5, len(reflect.ObjectReader.GetPropertyNames(obj)))

	dict := map[string]any{
		"1": "AAA",
		"A": 111,
	}
	assert.Equal(t, 2, len(reflect.ObjectReader.GetPropertyNames(dict)))

	list := []any{"BBB", 222}
	assert.Equal(t, 2, len(reflect.ObjectReader.GetPropertyNames(list)))
}

func TestReaderGetProperties(t *testing.T) {
	obj := NewTestClass()
	values := reflect.ObjectReader.GetProperties(obj)
	assert.Equal(t, 5, len(values))
	assert.Equal(t, true, values["RootPublicProperty"])
	assert.Equal(t, "BBB", values["PublicField"])

	dict := map[string]any{
		"1": "AAA",
		"A": 111,
	}
	values = reflect.ObjectReader.GetProperties(dict)
	assert.Equal(t, 2, len(values))
	assert.Equal(t, "AAA", values["1"])
	assert.Equal(t, 111, values["A"])

	list := []any{"BBB", 222}
	values = reflect.ObjectReader.GetProperties(list)
	assert.Equal(t, 2, len(values))
	assert.Equal(t, "BBB", values["0"])
	assert.Equal(t, 222, values["1"])
}
