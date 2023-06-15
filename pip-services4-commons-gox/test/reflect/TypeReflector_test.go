package test_reflect

import (
	refl "reflect"
	"testing"

	"github.com/pip-services4/pip-services4-commons-go/reflect"
	"github.com/stretchr/testify/assert"
)

func TestTypeReflectorCreate(t *testing.T) {
	typ := refl.TypeOf(TestClass{})
	obj, err := reflect.TypeReflector.CreateInstanceByType(typ)
	assert.NotNil(t, obj)
	assert.Nil(t, err)

	typ = refl.TypeOf((*TestClass)(nil))
	obj, err = reflect.TypeReflector.CreateInstanceByType(typ)
	assert.NotNil(t, obj)
	assert.Nil(t, err)
}
