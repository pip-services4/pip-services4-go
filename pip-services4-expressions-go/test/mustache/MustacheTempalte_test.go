package test_mustache

import (
	"testing"

	mustache "github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/mustache"
	"github.com/stretchr/testify/assert"
)

func TestMustacheTemplate1(t *testing.T) {
	template := mustache.NewMustacheTemplate()
	template.SetTemplate("Hello, {{{NAME}}}{{ #if ESCLAMATION }}!{{/if}}{{{^ESCLAMATION}}}.{{{/ESCLAMATION}}}")
	variables := map[string]string{
		"NAME":        "Alex",
		"ESCLAMATION": "1",
	}
	result, err := template.EvaluateWithVariables(variables)
	assert.Nil(t, err)
	assert.Equal(t, "Hello, Alex!", result)

	defaultVariables := map[string]string{
		"name":        "Mike",
		"esclamation": "",
	}
	template.SetDefaultVariables(defaultVariables)

	result, err = template.Evaluate()
	assert.Nil(t, err)
	assert.Equal(t, "Hello, Mike.", result)
}
