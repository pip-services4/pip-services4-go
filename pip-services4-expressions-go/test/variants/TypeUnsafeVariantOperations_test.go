package test_variants

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/variants"
	"github.com/stretchr/testify/assert"
)

func TestUnsafeOperations(t *testing.T) {
	a := variants.NewVariant(123)
	manager := variants.NewTypeUnsafeVariantOperations()

	b, _ := manager.Convert(a, variants.Float)
	assert.Equal(t, variants.Float, b.Type())
	assert.Equal(t, float32(123.0), b.AsFloat())

	c := variants.NewVariant(2)
	v, _ := manager.Add(b, c)
	assert.Equal(t, float32(125.0), v.AsFloat())
	v, _ = manager.Sub(b, c)
	assert.Equal(t, float32(121.0), v.AsFloat())
	v, _ = manager.Equal(a, b)
	assert.True(t, v.AsBoolean())
}
