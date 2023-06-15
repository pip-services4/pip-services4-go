package test_reflect

import (
	"testing"

	"github.com/pip-services4/pip-services4-commons-go/reflect"
	"github.com/stretchr/testify/assert"
)

func TestReflectorHasProperty(t *testing.T) {
	obj := NewTestClass()

	assert.False(t, reflect.PropertyReflector.HasProperty(obj, "123"))
	assert.False(t, reflect.PropertyReflector.HasProperty(obj, "rootPrivateField"))
	assert.False(t, reflect.PropertyReflector.HasProperty(obj, "privateField"))
	assert.True(t, reflect.PropertyReflector.HasProperty(obj, "PublicField"))
	assert.True(t, reflect.PropertyReflector.HasProperty(obj, "NestedField"))
	assert.True(t, reflect.PropertyReflector.HasProperty(obj, "RootPublicProperty"))
	assert.True(t, reflect.PropertyReflector.HasProperty(obj, "PublicProperty"))

	// check by names in tags
	assert.False(t, reflect.PropertyReflector.HasProperty(obj, "root_private_field"))
	assert.False(t, reflect.PropertyReflector.HasProperty(obj, "private_field"))
	assert.True(t, reflect.PropertyReflector.HasProperty(obj, "public_field"))
	assert.False(t, reflect.PropertyReflector.HasProperty(obj, "nested_field"))
	assert.False(t, reflect.PropertyReflector.HasProperty(obj, "root_public_property"))
	assert.False(t, reflect.PropertyReflector.HasProperty(obj, "public_property"))
}

func TestReflectorGetProperty(t *testing.T) {
	obj := NewTestClass()

	assert.Nil(t, reflect.PropertyReflector.GetProperty(obj, "123"))
	assert.Nil(t, reflect.PropertyReflector.GetProperty(obj, "rootPrivateField"))
	assert.Nil(t, reflect.PropertyReflector.GetProperty(obj, "root_private_field"))

	assert.Nil(t, reflect.PropertyReflector.GetProperty(obj, "privateField"))
	assert.Equal(t, "BBB", reflect.PropertyReflector.GetProperty(obj, "PublicField"))
	assert.Equal(t, "BBB", reflect.PropertyReflector.GetProperty(obj, "public_field"))
	assert.NotNil(t, reflect.PropertyReflector.GetProperty(obj, "NestedField"))
	assert.Equal(t, true, reflect.PropertyReflector.GetProperty(obj, "RootPublicProperty"))
	assert.Equal(t, true, reflect.PropertyReflector.GetProperty(obj, "PublicProperty"))
}

func TestReflectorGetPropertyNames(t *testing.T) {
	obj := NewTestClass()

	assert.Equal(t, 5, len(reflect.PropertyReflector.GetPropertyNames(obj)))
}

func TestReflectorGetProperties(t *testing.T) {
	obj := NewTestClass()

	properties := reflect.PropertyReflector.GetProperties(obj)
	assert.Equal(t, 5, len(properties))
	assert.Equal(t, "BBB", properties["PublicField"])
	assert.Equal(t, true, properties["RootPublicProperty"])
	assert.Equal(t, true, properties["PublicProperty"])
}

func TestReflectorSetProperty(t *testing.T) {
	obj := NewTestClass()

	assert.Equal(t, "BBB", reflect.PropertyReflector.GetProperty(obj, "PublicField"))
	reflect.PropertyReflector.SetProperty(obj, "PublicField", "XYZ")
	assert.Equal(t, "XYZ", reflect.PropertyReflector.GetProperty(obj, "PublicField"))

	reflect.PropertyReflector.SetProperty(obj, "public_field", "CCC") // set over tag name
	assert.Equal(t, "CCC", reflect.PropertyReflector.GetProperty(obj, "PublicField"))

	assert.Equal(t, true, reflect.PropertyReflector.GetProperty(obj, "PublicProperty"))
	reflect.PropertyReflector.SetProperty(obj, "PublicProperty", false)
	assert.Equal(t, false, reflect.PropertyReflector.GetProperty(obj, "PublicProperty"))
}

func TestReflectorSetProperties(t *testing.T) {
	obj := NewTestClass()

	assert.Equal(t, "BBB", reflect.PropertyReflector.GetProperty(obj, "PublicField"))
	assert.Equal(t, true, reflect.PropertyReflector.GetProperty(obj, "PublicProperty"))

	values := map[string]any{
		"PublicField":    "XYZ",
		"PublicProperty": false,
	}

	reflect.PropertyReflector.SetProperties(obj, values)
	assert.Equal(t, "XYZ", reflect.PropertyReflector.GetProperty(obj, "PublicField"))
	assert.Equal(t, false, reflect.PropertyReflector.GetProperty(obj, "PublicProperty"))

	values = map[string]any{
		"public_field":   "CCC",
		"PublicProperty": true,
	}

	reflect.PropertyReflector.SetProperties(obj, values)

	reflect.PropertyReflector.SetProperty(obj, "public_field", "CCC") // set over tag name
	assert.Equal(t, "CCC", reflect.PropertyReflector.GetProperty(obj, "PublicField"))
	assert.Equal(t, true, reflect.PropertyReflector.GetProperty(obj, "PublicProperty"))
}
