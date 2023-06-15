package test_reflect

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/reflect"
	"github.com/stretchr/testify/assert"
)

func TestRecursiveHasProperty(t *testing.T) {
	obj := convert.JsonConverter.ToMap("{ \"value1\": 123, \"value2\": { \"value21\": 111, \"value22\": 222 }, \"value3\": [ 444, { \"value311\": 555 } ] }")

	has := reflect.RecursiveObjectReader.HasProperty(obj, "")
	assert.False(t, has)

	has = reflect.RecursiveObjectReader.HasProperty(obj, "value1")
	assert.True(t, has)

	has = reflect.RecursiveObjectReader.HasProperty(obj, "value2")
	assert.True(t, has)

	has = reflect.RecursiveObjectReader.HasProperty(obj, "value2.value21")
	assert.True(t, has)

	has = reflect.RecursiveObjectReader.HasProperty(obj, "value2.value31")
	assert.False(t, has)

	has = reflect.RecursiveObjectReader.HasProperty(obj, "value2.value21.value211")
	assert.False(t, has)

	has = reflect.RecursiveObjectReader.HasProperty(obj, "valueA.valueB.valueC")
	assert.False(t, has)

	has = reflect.RecursiveObjectReader.HasProperty(obj, "value3")
	assert.True(t, has)

	has = reflect.RecursiveObjectReader.HasProperty(obj, "value3.0")
	assert.True(t, has)

	has = reflect.RecursiveObjectReader.HasProperty(obj, "value3.0.value311")
	assert.False(t, has)

	has = reflect.RecursiveObjectReader.HasProperty(obj, "value3.1")
	assert.True(t, has)

	has = reflect.RecursiveObjectReader.HasProperty(obj, "value3.1.value311")
	assert.True(t, has)

	has = reflect.RecursiveObjectReader.HasProperty(obj, "value3.2")
	assert.False(t, has)
}

func TestRecursiveGetProperty(t *testing.T) {
	obj := convert.JsonConverter.ToMap("{ \"value1\": 123, \"value2\": { \"value21\": 111, \"value22\": 222 }, \"value3\": [ 444, { \"value311\": 555 } ] }")

	value := reflect.RecursiveObjectReader.GetProperty(obj, "")
	assert.Nil(t, value)

	value = reflect.RecursiveObjectReader.GetProperty(obj, "value1")
	assert.Equal(t, float64(123), value)

	value = reflect.RecursiveObjectReader.GetProperty(obj, "value2")
	assert.NotNil(t, value)

	value = reflect.RecursiveObjectReader.GetProperty(obj, "value2.value21")
	assert.Equal(t, float64(111), value)

	value = reflect.RecursiveObjectReader.GetProperty(obj, "value2.value31")
	assert.Nil(t, value)

	value = reflect.RecursiveObjectReader.GetProperty(obj, "value2.value21.value211")
	assert.Nil(t, value)

	value = reflect.RecursiveObjectReader.GetProperty(obj, "valueA.valueB.valueC")
	assert.Nil(t, value)

	value = reflect.RecursiveObjectReader.GetProperty(obj, "value3")
	assert.NotNil(t, value)

	value = reflect.RecursiveObjectReader.GetProperty(obj, "value3.0")
	assert.Equal(t, float64(444), value)

	value = reflect.RecursiveObjectReader.GetProperty(obj, "value3.0.value311")
	assert.Nil(t, value)

	value = reflect.RecursiveObjectReader.GetProperty(obj, "value3.1")
	assert.NotNil(t, value)

	value = reflect.RecursiveObjectReader.GetProperty(obj, "value3.1.value311")
	assert.Equal(t, float64(555), value)

	value = reflect.RecursiveObjectReader.GetProperty(obj, "value3.2")
	assert.Nil(t, value)
}

func TestRecursiveGetPropertyName(t *testing.T) {
	obj := convert.JsonConverter.ToMap("{ \"value1\": 123, \"value2\": { \"value21\": 111, \"value22\": 222 }, \"value3\": [ 444, { \"value311\": 555 } ] }")

	contains := func(values []string, value string) bool {
		for _, v := range values {
			if v == value {
				return true
			}
		}
		return false
	}

	names := reflect.RecursiveObjectReader.GetPropertyNames(obj)
	assert.Equal(t, 5, len(names))
	assert.True(t, contains(names, "value1"))
	assert.True(t, contains(names, "value2.value21"))
	assert.True(t, contains(names, "value2.value22"))
	assert.True(t, contains(names, "value3.0"))
	assert.True(t, contains(names, "value3.1.value311"))
}

func TestRecursiveGetProperties(t *testing.T) {
	obj := convert.JsonConverter.ToMap("{ \"value1\": 123, \"value2\": { \"value21\": 111, \"value22\": 222 }, \"value3\": [ 444, { \"value311\": 555 } ] }")

	values := reflect.RecursiveObjectReader.GetProperties(obj)
	assert.Equal(t, 5, len(values))
	assert.Equal(t, float64(123), values["value1"])
	assert.Equal(t, float64(111), values["value2.value21"])
	assert.Equal(t, float64(222), values["value2.value22"])
	assert.Equal(t, float64(444), values["value3.0"])
	assert.Equal(t, float64(555), values["value3.1.value311"])
}
