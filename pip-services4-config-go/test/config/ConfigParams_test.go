package test_config

import (
	"testing"

	pconfig "github.com/pip-services4/pip-services4-go/pip-services4-components-go/config"
	mustache "github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/mustache"
	"github.com/stretchr/testify/assert"
)

func TestConfigSections(t *testing.T) {
	config := pconfig.NewConfigParamsFromTuples(
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

	section2 := pconfig.NewConfigParamsFromTuples(
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
	config := pconfig.NewConfigParamsFromString(
		"Queue=TestQueue;Endpoint=sb://cvctestbus.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=K70UpCUXN1Q5RFykll6/gz4Et14iJrYFnGPlwiFBlow=")
	assert.Equal(t, config.Len(), 4)

	val, ok := config.Get("Queue")
	assert.True(t, ok)
	assert.Equal(t, "TestQueue", val)
}

func TestProcessTemplates(t *testing.T) {
	template := mustache.NewMustacheTemplate()
	template.SetTemplate("{{#if A}}{{B}}{{/if}}")
	params := map[string]string{"A": "true", "B": "XYZ"}

	result, err := template.EvaluateWithVariables(params)

	assert.Nil(t, err)
	assert.Equal(t, "XYZ", result)
}
