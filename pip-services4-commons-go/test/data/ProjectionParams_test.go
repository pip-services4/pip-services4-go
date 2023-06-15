package test_data

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
	"github.com/stretchr/testify/assert"
)

func TestProjectionParamsFromNull(t *testing.T) {
	parameters := data.NewProjectionParamsFromValue(nil)

	assert.Equal(t, 0, parameters.Len())
}

func TestProjectionParamsFromValue(t *testing.T) {
	parameters := data.NewProjectionParamsFromValue([]any{"field1", "field2", "field3"})

	assert.Equal(t, 3, parameters.Len())
	val, ok := parameters.Get(0)
	assert.True(t, ok)
	assert.Equal(t, "field1", val)

	val, ok = parameters.Get(1)
	assert.True(t, ok)
	assert.Equal(t, "field2", val)

	val, ok = parameters.Get(2)
	assert.True(t, ok)
	assert.Equal(t, "field3", val)
}

func TestParseProjectionParams(t *testing.T) {
	parameters := data.ParseProjectionParams("field1", "field2", "field3")

	assert.Equal(t, 3, parameters.Len())
	val, ok := parameters.Get(0)
	assert.True(t, ok)
	assert.Equal(t, "field1", val)

	val, ok = parameters.Get(1)
	assert.True(t, ok)
	assert.Equal(t, "field2", val)

	val, ok = parameters.Get(2)
	assert.True(t, ok)
	assert.Equal(t, "field3", val)

	parameters = data.ParseProjectionParams("field1,field2, field3")

	assert.Equal(t, 3, parameters.Len())
	val, ok = parameters.Get(0)
	assert.True(t, ok)
	assert.Equal(t, "field1", val)

	val, ok = parameters.Get(1)
	assert.True(t, ok)
	assert.Equal(t, "field2", val)

	val, ok = parameters.Get(2)
	assert.True(t, ok)
	assert.Equal(t, "field3", val)

	parameters = data.ParseProjectionParams("object1(field1)", "object2(field1, field2)", "field3")

	assert.Equal(t, 4, parameters.Len())

	val, ok = parameters.Get(0)
	assert.True(t, ok)
	assert.Equal(t, "object1.field1", val)

	val, ok = parameters.Get(1)
	assert.True(t, ok)
	assert.Equal(t, "object2.field1", val)

	val, ok = parameters.Get(2)
	assert.True(t, ok)
	assert.Equal(t, "object2.field2", val)

	val, ok = parameters.Get(3)
	assert.True(t, ok)
	assert.Equal(t, "field3", val)

	parameters = data.ParseProjectionParams("object1(object2(field1,field2,object3(field1)))")

	assert.Equal(t, 3, parameters.Len())

	val, ok = parameters.Get(0)
	assert.True(t, ok)
	assert.Equal(t, "object1.object2.field1", val)

	val, ok = parameters.Get(1)
	assert.True(t, ok)
	assert.Equal(t, "object1.object2.field2", val)

	val, ok = parameters.Get(2)
	assert.True(t, ok)
	assert.Equal(t, "object1.object2.object3.field1", val)

	parameters = data.ParseProjectionParams("object1(field1, object2(field1, field2, field3, field4), field3)", "field2")

	assert.Equal(t, 7, parameters.Len())
	val, ok = parameters.Get(0)
	assert.True(t, ok)
	assert.Equal(t, "object1.field1", val)
	val, ok = parameters.Get(1)
	assert.True(t, ok)
	assert.Equal(t, "object1.object2.field1", val)
	val, ok = parameters.Get(2)
	assert.True(t, ok)
	assert.Equal(t, "object1.object2.field2", val)
	val, ok = parameters.Get(3)
	assert.True(t, ok)
	assert.Equal(t, "object1.object2.field3", val)
	val, ok = parameters.Get(4)
	assert.True(t, ok)
	assert.Equal(t, "object1.object2.field4", val)
	val, ok = parameters.Get(5)
	assert.True(t, ok)
	assert.Equal(t, "object1.field3", val)
	val, ok = parameters.Get(6)
	assert.True(t, ok)
	assert.Equal(t, "field2", val)

	parameters = data.ParseProjectionParams("object1(field1, object2(field1), field3)", "field2")

	assert.Equal(t, 4, parameters.Len())
	val, ok = parameters.Get(0)
	assert.True(t, ok)
	assert.Equal(t, "object1.field1", val)
	val, ok = parameters.Get(1)
	assert.True(t, ok)
	assert.Equal(t, "object1.object2.field1", val)
	val, ok = parameters.Get(2)
	assert.True(t, ok)
	assert.Equal(t, "object1.field3", val)
	val, ok = parameters.Get(3)
	assert.True(t, ok)
	assert.Equal(t, "field2", val)

	parameters = data.ParseProjectionParams("object1(field1, object2(field1, field2, object3(field1), field4), field3)", "field2")

	assert.Equal(t, 7, parameters.Len())

	val, ok = parameters.Get(0)
	assert.True(t, ok)
	assert.Equal(t, "object1.field1", val)
	val, ok = parameters.Get(1)
	assert.True(t, ok)
	assert.Equal(t, "object1.object2.field1", val)
	val, ok = parameters.Get(2)
	assert.True(t, ok)
	assert.Equal(t, "object1.object2.field2", val)
	val, ok = parameters.Get(3)
	assert.True(t, ok)
	assert.Equal(t, "object1.object2.object3.field1", val)
	val, ok = parameters.Get(4)
	assert.True(t, ok)
	assert.Equal(t, "object1.object2.field4", val)
	val, ok = parameters.Get(5)
	assert.True(t, ok)
	assert.Equal(t, "object1.field3", val)
	val, ok = parameters.Get(6)
	assert.True(t, ok)
	assert.Equal(t, "field2", val)

	parameters = data.ParseProjectionParams("object1(object2(object3(field1)), field2)", "field2")

	assert.Equal(t, 3, parameters.Len())
	val, ok = parameters.Get(0)
	assert.True(t, ok)
	assert.Equal(t, "object1.object2.object3.field1", val)
	val, ok = parameters.Get(1)
	assert.True(t, ok)
	assert.Equal(t, "object1.field2", val)
	val, ok = parameters.Get(2)
	assert.True(t, ok)
	assert.Equal(t, "field2", val)

	parameters = data.ParseProjectionParams("field1,object1(field1),object2(field1,field2),object3(field1),field2,field3")

	assert.Equal(t, 7, parameters.Len())
	val, ok = parameters.Get(0)
	assert.True(t, ok)
	assert.Equal(t, "field1", val)
	val, ok = parameters.Get(1)
	assert.True(t, ok)
	assert.Equal(t, "object1.field1", val)
	val, ok = parameters.Get(2)
	assert.True(t, ok)
	assert.Equal(t, "object2.field1", val)
	val, ok = parameters.Get(3)
	assert.True(t, ok)
	assert.Equal(t, "object2.field2", val)
	val, ok = parameters.Get(4)
	assert.True(t, ok)
	assert.Equal(t, "object3.field1", val)
	val, ok = parameters.Get(5)
	assert.True(t, ok)
	assert.Equal(t, "field2", val)
	val, ok = parameters.Get(6)
	assert.True(t, ok)
	assert.Equal(t, "field3", val)
}
