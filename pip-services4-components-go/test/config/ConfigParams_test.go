package test_config

import (
	"encoding/json"
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	conf "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	"github.com/stretchr/testify/assert"
)

func TestConfigSections(t *testing.T) {
	config := conf.NewConfigParamsFromTuples(
		"Section1.Key1", "Value1",
		"Section1.Key2", "Value2",
		"Section1.Key3", "Value3",
	)

	assert.Equal(t, config.Len(), 3)

	val, ok := config.Get("Section1.Key1")
	assert.True(t, ok)
	assert.Equal(t, "Value1", val)

	val, ok = config.Get("Section1.Key2")
	assert.True(t, ok)
	assert.Equal(t, "Value2", val)

	val, ok = config.Get("Section1.Key3")
	assert.True(t, ok)
	assert.Equal(t, "Value3", val)

	val, ok = config.Get("Section1.Key4")
	assert.False(t, ok)
	assert.Equal(t, "", val)

	section2 := conf.NewConfigParamsFromTuples(
		"Key1", "ValueA",
		"Key2", "ValueB",
	)

	config.AddSection("Section2", section2)
	assert.Equal(t, config.Len(), 5)

	val, ok = config.Get("Section2.Key1")
	assert.True(t, ok)
	assert.Equal(t, "ValueA", val)

	val, ok = config.Get("Section2.Key2")
	assert.True(t, ok)
	assert.Equal(t, "ValueB", val)

	section1 := config.GetSection("Section1")
	assert.Equal(t, section1.Len(), 3)

	val, ok = section1.Get("Key1")
	assert.True(t, ok)
	assert.Equal(t, "Value1", val)

	val, ok = section1.Get("Key2")
	assert.True(t, ok)
	assert.Equal(t, "Value2", val)

	val, ok = section1.Get("Key3")
	assert.True(t, ok)
	assert.Equal(t, "Value3", val)
}

func TestConfigFromString(t *testing.T) {
	config := conf.NewConfigParamsFromString("Queue=TestQueue;Endpoint=sb://cvctestbus.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=K70UpCUXN1Q5RFykll6/gz4Et14iJrYFnGPlwiFBlow=")
	assert.Equal(t, config.Len(), 4)

	val, ok := config.Get("Queue")
	assert.True(t, ok)
	assert.Equal(t, "TestQueue", val)
}

func TestConfigJsonSerialization(t *testing.T) {
	json1 := []byte("{\"key1\":\"1\",\"key2\":\"A\"}")

	var value *data.StringValueMap
	err := json.Unmarshal(json1, &value)
	assert.Empty(t, err)

	val, ok := value.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "1", val)

	val, ok = value.Get("key2")
	assert.True(t, ok)
	assert.Equal(t, "A", val)

	json2, err2 := json.Marshal(value)
	assert.Empty(t, err2)
	assert.Equal(t, json1, json2)
}

// func TestConfigFromObject(t *testing.T) {
// 	value := data.NewAnyValueMapFromTuples(
// 		"field1", conf.NewConfigParamsFromString("field11=123;field12=ABC"),
// 		"field2", data.NewAnyValueArrayFromValue(
// 			123, "ABC", conf.NewConfigParamsFromString("field21=543;field22=XYZ"),
// 		),
// 		"field3", true,
// 	)

// 	config := conf.NewConfigParamsFromValue(value)
// 	assert.Equal(t, config.Len(), 7)
// 	assert.Equal(t, config.GetAsInteger("field1.field11"), 123)
// 	assert.Equal(t, config.Get("field1.field12"), "ABC")
// 	assert.Equal(t, config.GetAsInteger("field2.0"), 123)
// 	assert.Equal(t, config.Get("field2.1"), "ABC")
// 	assert.Equal(t, config.GetAsInteger("field2.2.field21"), 543)
// 	assert.Equal(t, config.Get("field2.2.field22"), "XYZ")
// 	assert.Equal(t, config.GetAsBoolean("field3"), true)
// }
