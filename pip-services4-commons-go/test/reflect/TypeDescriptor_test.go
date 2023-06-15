package test_reflect

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/reflect"
	"github.com/stretchr/testify/assert"
)

func TestTypeDescriptorFromString(t *testing.T) {
	descriptor, err := reflect.ParseTypeDescriptorFromString("")
	assert.Nil(t, descriptor)
	assert.Nil(t, err)

	descriptor, err = reflect.ParseTypeDescriptorFromString("xxx,yyy")
	assert.Equal(t, "xxx", descriptor.Name())
	assert.Equal(t, "yyy", descriptor.Package())
	assert.Nil(t, err)

	descriptor, err = reflect.ParseTypeDescriptorFromString("xxx")
	assert.Equal(t, "xxx", descriptor.Name())
	assert.Equal(t, "", descriptor.Package())

	descriptor, err = reflect.ParseTypeDescriptorFromString("xxx,yyy,zzz")
	assert.NotNil(t, err)
}
