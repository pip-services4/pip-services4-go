package test_variants

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/variants"
	"github.com/stretchr/testify/assert"
)

func TestSafeOperations(t *testing.T) {
	a := variants.NewVariant(123)
	manager := variants.NewTypeSafeVariantOperations()

	b, _ := manager.Convert(a, variants.Float)
	assert.Equal(t, variants.Float, b.Type())
	assert.Equal(t, float32(123.0), b.AsFloat())

	c := variants.NewVariant(2)
	v, _ := manager.Add(a, c)
	assert.Equal(t, 125, v.AsInteger())
	v, _ = manager.Sub(a, c)
	assert.Equal(t, 121, v.AsInteger())
	v, _ = manager.Equal(a, c)
	assert.False(t, v.AsBoolean())

	array := []*variants.Variant{
		variants.NewVariant("aaa"),
		variants.NewVariant("bbb"),
		variants.NewVariant("ccc"),
		variants.NewVariant("ddd"),
	}
	d := variants.NewVariant(array)
	v, _ = manager.In(d, variants.NewVariant("ccc"))
	assert.True(t, v.AsBoolean())
	v, _ = manager.In(d, variants.NewVariant("eee"))
	assert.False(t, v.AsBoolean())
	v, _ = manager.GetElement(d, variants.NewVariant(1))
	assert.Equal(t, "bbb", v.AsString())
}
