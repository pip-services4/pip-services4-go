package test_build

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/build"
	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/stretchr/testify/assert"
)

func newObject() interface{} {
	return "ABC"
}

func TestFactoryByType(t *testing.T) {
	factory := build.NewFactory()
	descriptor := refer.NewDescriptor("test", "object", "default", "*", "1.0")

	factory.RegisterType(descriptor, newObject)

	locator := factory.CanCreate(descriptor)
	assert.NotNil(t, locator)
	locator = factory.CanCreate("123")
	assert.Nil(t, locator)

	obj, err := factory.Create(descriptor)
	assert.Nil(t, err)
	assert.Equal(t, "ABC", obj)
	obj, err = factory.Create("123")
	assert.NotNil(t, err)
	assert.Nil(t, obj)
}

func TestFactory(t *testing.T) {
	factory := build.NewFactory()
	descriptor := refer.NewDescriptor("test", "object", "default", "*", "1.0")

	factory.Register(descriptor, func(locator any) any {
		name := ""
		descriptor, ok := locator.(*refer.Descriptor)
		if ok {
			name = descriptor.String()
		}
		t.Log("Factory component name:", name)
		return newObject()
	},
	)

	locator := factory.CanCreate(descriptor)
	assert.NotNil(t, locator)

	obj, err := factory.Create(descriptor)
	assert.Nil(t, err)
	assert.Equal(t, "ABC", obj)

}
