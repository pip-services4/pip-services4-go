package test_calculator_variables

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/calculator/variables"
	"github.com/stretchr/testify/assert"
)

func TestVariableCollectionAddRemoveVariables(t *testing.T) {
	collection := variables.NewVariableCollection()

	var1 := variables.EmptyVariable("ABC")
	collection.Add(var1)
	assert.Equal(t, 1, collection.Length())

	var2 := variables.EmptyVariable("XYZ")
	collection.Add(var2)
	assert.Equal(t, 2, collection.Length())

	index := collection.FindIndexByName("abc")
	assert.Equal(t, 0, index)

	v := collection.FindByName("Xyz")
	assert.Equal(t, var2, v)

	var3 := collection.Locate("ghi")
	assert.NotNil(t, var3)
	assert.Equal(t, "ghi", var3.Name())
	assert.Equal(t, 3, collection.Length())

	collection.Remove(0)
	assert.Equal(t, 2, collection.Length())

	collection.RemoveByName("GHI")
	assert.Equal(t, 1, collection.Length())
}
