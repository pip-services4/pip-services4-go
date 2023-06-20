package test_validate

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
	"github.com/stretchr/testify/assert"
)

func TestEmptySchema(t *testing.T) {
	schema := validate.NewSchema()
	results := schema.Validate(nil)
	assert.Equal(t, 0, len(results))
}

func TestSchemaRequired(t *testing.T) {
	schema := validate.NewSchema().MakeRequired()
	results := schema.Validate(nil)
	assert.Equal(t, 1, len(results))
}
