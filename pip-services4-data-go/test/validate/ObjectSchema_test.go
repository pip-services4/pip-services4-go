package test_validate

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	"github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
	"github.com/stretchr/testify/assert"
)

func TestObjectSchemaWithOptionalProperties(t *testing.T) {
	schema := validate.NewObjectSchema().
		WithOptionalProperty("intField", nil).
		WithOptionalProperty("stringField1", nil).
		WithOptionalProperty("stringField2", nil).
		WithOptionalProperty("intArrayField", nil).
		WithOptionalProperty("stringArrayField", nil).
		WithOptionalProperty("mapField", nil).
		WithOptionalProperty("subObjectField", nil).
		WithOptionalProperty("subArrayField", nil)

	obj := &TestClass{}
	results := schema.Validate(obj)
	assert.Equal(t, 0, len(results))
}

func TestObjectSchemaWithRequiredProperties(t *testing.T) {
	schema := validate.NewObjectSchema().
		WithRequiredProperty("intField", nil).
		WithRequiredProperty("stringField1", nil).
		WithRequiredProperty("stringField2", nil).
		WithRequiredProperty("intArrayField", nil).
		WithRequiredProperty("stringArrayField", nil).
		WithRequiredProperty("mapField", nil).
		WithRequiredProperty("subObjectField", nil).
		WithRequiredProperty("subArrayField", nil)

	obj := &TestClass{}
	results := schema.Validate(obj)
	assert.Equal(t, 0, len(results))

	obj = &TestClass{
		IntArrayField:    []int{1, 2, 3},
		StringArrayField: []string{"A", "B", "C"},
		MapField:         map[string]any{},
		SubObjectField:   &SubTestClass{},
		SubArrayField:    []*SubTestClass{},
	}
	results = schema.Validate(obj)
	assert.Equal(t, 0, len(results))
}

func TestObjectSchemaTypes(t *testing.T) {
	schema := validate.NewObjectSchema().
		WithRequiredProperty("intField", convert.Integer).
		WithRequiredProperty("stringField1", convert.String).
		WithRequiredProperty("stringField2", convert.String).
		WithRequiredProperty("intArrayField", convert.Array).
		WithRequiredProperty("stringArrayField", convert.Array).
		WithRequiredProperty("mapField", convert.Map).
		WithRequiredProperty("subObjectField", convert.Object).
		WithRequiredProperty("subArrayField", convert.Array)

	obj := &TestClass{
		IntArrayField:    []int{1, 2, 3},
		StringArrayField: []string{"A", "B", "C"},
		MapField:         map[string]any{},
		SubObjectField:   &SubTestClass{},
		SubArrayField:    []*SubTestClass{},
	}
	results := schema.Validate(obj)
	assert.Equal(t, 0, len(results))
}

func TestObjectSubSchema(t *testing.T) {
	subSchema := validate.NewObjectSchema().
		WithRequiredProperty("id", convert.String).
		WithRequiredProperty("floatField", convert.Float)

	schema := validate.NewObjectSchema().
		WithRequiredProperty("intField", convert.Integer).
		WithRequiredProperty("stringField1", convert.String).
		WithRequiredProperty("stringField2", convert.String).
		WithRequiredProperty("intArrayField", convert.Array).
		WithRequiredProperty("stringArrayField", convert.Array).
		WithRequiredProperty("mapField", convert.Map).
		WithRequiredProperty("subObjectField", subSchema).
		WithRequiredProperty("subArrayField", convert.Array)

	obj := &TestClass{
		IntArrayField:    []int{1, 2, 3},
		StringArrayField: []string{"A", "B", "C"},
		MapField:         map[string]any{},
		SubObjectField:   &SubTestClass{},
		SubArrayField:    []*SubTestClass{},
	}
	results := schema.Validate(obj)
	assert.Equal(t, 0, len(results))
}
