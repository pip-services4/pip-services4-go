package test_reflect

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/reflect"
	"github.com/stretchr/testify/assert"
)

func TestReflectorGetMethods(t *testing.T) {
	obj := NewTestClass()

	methods := reflect.MethodReflector.GetMethodNames(obj)
	assert.Equal(t, 7, len(methods))
}

func TestReflectorHasMethod(t *testing.T) {
	obj := NewTestClass()

	has := reflect.MethodReflector.HasMethod(obj, "pUblIcMeThoD")
	assert.True(t, has)
}

func TestReflectorInvokeMethod(t *testing.T) {
	obj := NewTestClass()

	result := reflect.MethodReflector.InvokeMethod(obj, "PUBLICMETHOD", 1, 2)
	assert.Equal(t, 3, result)
}
