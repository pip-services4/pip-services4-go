package test_refer

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/stretchr/testify/assert"
)

func TestSimpleReference(t *testing.T) {
	ref := refer.NewReference("ABC", 123)
	assert.Equal(t, "ABC", ref.Locator())
	assert.Equal(t, 123, ref.Component())
	assert.True(t, ref.Match("ABC"))
	assert.True(t, ref.Match(123))
	assert.False(t, ref.Match("XYZ"))
	assert.False(t, ref.Match(321))
}

func TestDescriptorReference(t *testing.T) {
	var descriptor1 = refer.NewDescriptor("pip-dummies", "controller", "default", "default", "1.0")
	var descriptor2 = refer.NewDescriptor("pip-dummies", "controller", "default", "default", "1.0")

	ref := refer.NewReference(descriptor1, 123)
	assert.Equal(t, descriptor1, ref.Locator())
	assert.Equal(t, 123, ref.Component())
	assert.True(t, ref.Match(descriptor1))
	assert.True(t, ref.Match(descriptor2))
	assert.True(t, ref.Match(123))
	assert.False(t, ref.Match("XYZ"))
	assert.False(t, ref.Match(321))
}
