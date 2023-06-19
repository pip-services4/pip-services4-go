package test_exec

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/exec"
	"github.com/stretchr/testify/assert"
)

func TestGetParams(t *testing.T) {
	obj := convert.JsonConverter.ToMap("{ \"value1\": 123, \"value2\": { \"value21\": 111, \"value22\": 222 } }")
	params := exec.NewParametersFromValue(obj)

	value, ok := params.Get("")
	assert.False(t, ok)
	assert.Nil(t, value)

	value = params.GetAsInteger("value1")
	assert.NotNil(t, value)
	assert.Equal(t, 123, value)

	value, ok = params.Get("value2")
	assert.True(t, ok)
	assert.NotNil(t, value)

	boolVal := params.Contains("value3")
	assert.False(t, boolVal)

	value = params.GetAsInteger("value2.value21")
	assert.NotNil(t, value)
	assert.Equal(t, 111, value)

	boolVal = params.Contains("value2.value31")
	assert.False(t, boolVal)

	boolVal = params.Contains("value2.value21.value211")
	assert.False(t, boolVal)

	boolVal = params.Contains("valueA.valueB.valueC")
	assert.False(t, boolVal)
}

func TestParamsContainKey(t *testing.T) {
	obj := convert.JsonConverter.ToMap("{ \"value1\": 123, \"value2\": { \"value21\": 111, \"value22\": 222 } }")
	params := exec.NewParametersFromValue(obj)

	has := params.Contains("")
	assert.False(t, has)

	has = params.Contains("value1")
	assert.True(t, has)

	has = params.Contains("value2")
	assert.True(t, has)

	has = params.Contains("value3")
	assert.False(t, has)

	has = params.Contains("value2.value21")
	assert.True(t, has)

	has = params.Contains("value2.value31")
	assert.False(t, has)

	has = params.Contains("value2.value21.value211")
	assert.False(t, has)

	has = params.Contains("valueA.valueB.valueC")
	assert.False(t, has)
}

func TestSetParams(t *testing.T) {
	params := exec.NewEmptyParameters()

	params.Put("", 123)
	assert.Equal(t, 0, params.Len())

	params.Put("field1", 123)
	assert.Equal(t, 1, params.Len())
	assert.Equal(t, 123, params.GetAsInteger("field1"))

	params.Put("field2", "ABC")
	assert.Equal(t, 2, params.Len())

	val, ok := params.Get("field2")
	assert.True(t, ok)
	assert.Equal(t, "ABC", val)

	params.Put("field2.field1", 123)

	val, ok = params.Get("field2")
	assert.True(t, ok)
	assert.Equal(t, "ABC", val)

	// params.Put("field3.field31", 456)
	// assert.Equal(t, 3, params.Len())
	// subConfig := params.GetAsMap("field3")
	// assert.NotNil(t, subConfig)
	// assert.Equal(t, 456, subConfig.GetAsInteger("field31"))

	// params.Put("field3.field32", "XYZ")
	// assert.Equal(t, "XYZ", params.Get("field3.field32"))
}

func TestParamsDefaults(t *testing.T) {
	result := exec.NewParametersFromTuples(
		"value1", 123,
		"value2", 234,
	)
	defaults := exec.NewParametersFromTuples(
		"value2", 432,
		"value3", 345,
	)
	result = result.SetDefaults(defaults, false)
	assert.Equal(t, 3, result.Len())

	val, ok := result.Get("value1")
	assert.True(t, ok)
	assert.Equal(t, 123, val)

	val, ok = result.Get("value2")
	assert.True(t, ok)
	assert.Equal(t, 234, val)

	val, ok = result.Get("value3")
	assert.True(t, ok)
	assert.Equal(t, 345, val)
}

func TestParamsOverrideRecursive(t *testing.T) {
	obj := convert.JsonConverter.ToMap("{ \"value1\": 123, \"value2\": { \"value21\": 111, \"value22\": 222 } }")
	result := exec.NewParametersFromValue(obj)
	obj = convert.JsonConverter.ToMap("{ \"value2\": { \"value22\": 777, \"value23\": 333 }, \"value3\": 345 }")
	defaults := exec.NewParametersFromValue(obj)
	result = result.SetDefaults(defaults, true)

	assert.Equal(t, 3, result.Len())

	val, ok := result.Get("value1")
	assert.True(t, ok)
	assert.Equal(t, float64(123), val)

	val, ok = result.Get("value3")
	assert.True(t, ok)
	assert.Equal(t, float64(345), val)

	deepResult := result.GetAsMap("value2")
	assert.Equal(t, 3, deepResult.Len())
	assert.Equal(t, 111, deepResult.GetAsInteger("value21"))
	assert.Equal(t, 222, deepResult.GetAsInteger("value22"))
	assert.Equal(t, 333, deepResult.GetAsInteger("value23"))
}

func TestParamsOverrideNulls(t *testing.T) {
	obj := convert.JsonConverter.ToMap("{ \"value1\": 123, \"value2\": 234 }")
	result := exec.NewParametersFromValue(obj)
	result = result.Override(nil, true)

	assert.Equal(t, 2, result.Len())

	val, ok := result.Get("value1")
	assert.True(t, ok)
	assert.Equal(t, float64(123), val)

	val, ok = result.Get("value2")
	assert.True(t, ok)
	assert.Equal(t, float64(234), val)
}

func TestParamsAssignTo(t *testing.T) {
	value := NewTestClass()
	obj := convert.JsonConverter.ToMap("{ \"value1\": 123, \"value2\": \"ABC\", \"value3\": 456 }")
	newValues := exec.NewParametersFromValue(obj)

	newValues.AssignTo(value)
	// assert.NotNil(t, value.value1)
	// assert.Equal(t, 123, value.value1)
	// assert.NotNil(t, value.value2)
	// assert.Equal(t, "ABC", value.value2)
}

func TestParamsFromConfig(t *testing.T) {
	params := config.NewConfigParamsFromTuples(
		"field1.field11", 123,
		"field2", "ABC",
		"field1.field12", "XYZ",
	)

	parameters := exec.NewParametersFromConfig(params)
	assert.Equal(t, 2, parameters.Len())

	val, ok := parameters.Get("field2")
	assert.True(t, ok)
	assert.Equal(t, "ABC", val)

	value := parameters.GetAsMap("field1")
	assert.Equal(t, 2, value.Len())

	val, ok = value.Get("field11")
	assert.True(t, ok)
	assert.Equal(t, "123", val)

	val, ok = value.Get("field12")
	assert.True(t, ok)
	assert.Equal(t, "XYZ", val)
}
