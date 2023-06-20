package test_variants

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/variants"
	"github.com/stretchr/testify/assert"
)

func TestVariants(t *testing.T) {
	a := variants.NewVariant(123)
	assert.Equal(t, variants.Integer, a.Type())
	assert.Equal(t, 123, a.AsInteger())
	assert.Equal(t, 123, a.AsObject())

	b := variants.NewVariant("xyz")
	assert.Equal(t, variants.String, b.Type())
	assert.Equal(t, "xyz", b.AsString())
	assert.Equal(t, "xyz", b.AsObject())
}
