package test_refer

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-components-go/refer"
	"github.com/stretchr/testify/assert"
)

func TestDescriptorMatch(t *testing.T) {
	descriptor := refer.NewDescriptor("pip-dummies", "controller", "default", "default", "1.0")

	// Check match by individual fields
	assert.True(t, descriptor.Match(refer.NewDescriptor("", "", "", "", "")))
	assert.True(t, descriptor.Match(refer.NewDescriptor("pip-dummies", "controller", "", "", "")))
	assert.True(t, descriptor.Match(refer.NewDescriptor("", "", "default", "", "")))
	assert.True(t, descriptor.Match(refer.NewDescriptor("", "", "", "", "1.0")))

	// Check match by individual "*" fields
	assert.True(t, descriptor.Match(refer.NewDescriptor("pip-dummies", "*", "*", "*", "*")))
	assert.True(t, descriptor.Match(refer.NewDescriptor("*", "controller", "*", "*", "*")))
	assert.True(t, descriptor.Match(refer.NewDescriptor("*", "*", "default", "*", "*")))
	assert.True(t, descriptor.Match(refer.NewDescriptor("*", "*", "*", "*", "1.0")))

	// Check match by all values
	assert.True(t, descriptor.Match(refer.NewDescriptor("pip-dummies", "controller", "default", "default", "")))
	assert.True(t, descriptor.Match(refer.NewDescriptor("", "controller", "default", "default", "1.0")))
	assert.True(t, descriptor.Match(refer.NewDescriptor("pip-dummies", "controller", "default", "default", "1.0")))

	// Check mismatch by individual fields
	assert.False(t, descriptor.Match(refer.NewDescriptor("", "cache", "", "", "")))
	assert.False(t, descriptor.Match(refer.NewDescriptor("pip-commons", "controller", "", "", "")))
	assert.False(t, descriptor.Match(refer.NewDescriptor("", "", "special", "", "")))
	assert.False(t, descriptor.Match(refer.NewDescriptor("", "", "", "", "2.0")))
}

func TestDescriptorToString(t *testing.T) {
	descriptor1 := refer.NewDescriptor("pip-dummies", "controller", "default", "default", "1.0")
	assert.Equal(t, "pip-dummies:controller:default:default:1.0", descriptor1.String())

	descriptor2 := refer.NewDescriptor("", "controller", "", "", "")
	assert.Equal(t, "*:controller:*:*:*", descriptor2.String())
}

func TestDescriptorFromString(t *testing.T) {
	descriptor, err := refer.ParseDescriptorFromString("")
	assert.Nil(t, err)
	assert.Nil(t, descriptor)

	descriptor, err = refer.ParseDescriptorFromString("pip-dummies:controller:default:default:1.0")
	assert.Nil(t, err)
	assert.True(t, descriptor.ExactMatch(refer.NewDescriptor("pip-dummies", "controller", "default", "default", "1.0")))

	descriptor, err = refer.ParseDescriptorFromString("xxx")
	assert.NotNil(t, err)
	assert.Nil(t, descriptor)
}

func TestDescriptorEquals(t *testing.T) {
	descriptor1 := refer.NewDescriptor("pip-commons", "controller", "", "default", "")
	descriptor2 := refer.NewDescriptor("", "controller", "", "", "1.0")

	match := descriptor1.Equals(descriptor2)
	assert.True(t, match)

	descriptor3 := refer.NewDescriptor("", "persistence", "", "", "")
	match = descriptor1.Equals(descriptor3)
	assert.False(t, match)
}
