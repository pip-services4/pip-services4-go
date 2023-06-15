package test_reflect

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/reflect"
	"github.com/stretchr/testify/assert"
)

func TestRecursiveSetProperty(t *testing.T) {
	obj := convert.JsonConverter.ToMap("{ \"value1\": 123, \"value2\": { \"value21\": 111, \"value22\": 222 }, \"value3\": [ 444, { \"value311\": 555 } ] }")

	//reflect.RecursiveObjectWriter.SetProperty(obj, "", null);
	reflect.RecursiveObjectWriter.SetProperty(obj, "value1", "AAA")
	reflect.RecursiveObjectWriter.SetProperty(obj, "value2", "BBB")
	reflect.RecursiveObjectWriter.SetProperty(obj, "value3.1.value312", "CCC")
	//reflect.RecursiveObjectWriter.SetProperty(obj, "value3.3", "DDD")
	reflect.RecursiveObjectWriter.SetProperty(obj, "value4.1", "EEE")

	values := reflect.RecursiveObjectReader.GetProperties(obj)
	assert.Equal(t, 6, len(values))
	assert.Equal(t, "AAA", values["value1"])
	assert.Equal(t, "BBB", values["value2"])
	//assert.Equal(t, 111, values["value2.value21"])
	assert.Equal(t, float64(444), values["value3.0"])
	assert.Equal(t, float64(555), values["value3.1.value311"])
	assert.Equal(t, "CCC", values["value3.1.value312"])
	//assert.Nil(t, values["value3.2"])
	//assert.Equal(t, "DDD", values["value3.3"])
	//assert.Nil(t, values["value4.0"])
	//assert.Equal(t, "EEE", values["value4.1"])
}

func TestRecursiveSetProperties(t *testing.T) {
	obj := convert.JsonConverter.ToMap("{ \"value1\": 123, \"value2\": { \"value21\": 111, \"value22\": 222 }, \"value3\": [ 444, { \"value311\": 555 } ] }")

	values := map[string]any{
		"value1":            "AAA",
		"value2":            "BBB",
		"value3.1.value312": "CCC",
		"value3.3":          "DDD",
		"value4.1":          "EEE",
	}
	reflect.RecursiveObjectWriter.SetProperties(obj, values)

	resultValues := reflect.RecursiveObjectReader.GetProperties(obj)
	assert.Equal(t, 6, len(resultValues))
	assert.Equal(t, "AAA", resultValues["value1"])
	assert.Equal(t, "BBB", resultValues["value2"])
	//assert.Equal(t, 111, resultValues["value2.value21"])
	assert.Equal(t, float64(444), resultValues["value3.0"])
	assert.Equal(t, float64(555), resultValues["value3.1.value311"])
	assert.Equal(t, "CCC", resultValues["value3.1.value312"])
	assert.Nil(t, resultValues["value3.2"])
	//assert.Equal(t, "DDD", resultValues["value3.3"])
	//assert.Nil(t, resultValues["value4.0"])
	//assert.Equal(t, "EEE", resultValues["value4.1"])
}
