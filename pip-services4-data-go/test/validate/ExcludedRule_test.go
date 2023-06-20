package test_validate

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-data-go/validate"
	"github.com/stretchr/testify/assert"
)

func TestExcludedRule(t *testing.T) {
	schema := validate.NewSchema().
		WithRule(validate.NewExcludedRule("AAA", "BBB", "CCC", nil))

	results := schema.Validate("AAA")
	assert.Equal(t, 1, len(results))

	results = schema.Validate("ABC")
	assert.Equal(t, 0, len(results))
}
