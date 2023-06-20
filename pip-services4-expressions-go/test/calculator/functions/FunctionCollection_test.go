package test_calculator_functions

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/calculator/functions"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/variants"
	"github.com/stretchr/testify/assert"
)

func testFunc(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {
	return variants.EmptyVariant(), nil
}

func TestFunctionsCollectionAddRemoveFunctions(t *testing.T) {
	collection := functions.NewFunctionCollection()

	func1 := functions.NewDelegatedFunction("ABC", testFunc)
	collection.Add(func1)
	assert.Equal(t, 1, collection.Length())

	func2 := functions.NewDelegatedFunction("XYZ", testFunc)
	collection.Add(func2)
	assert.Equal(t, 2, collection.Length())

	index := collection.FindIndexByName("abc")
	assert.Equal(t, 0, index)

	f := collection.FindByName("Xyz")
	assert.Equal(t, func2, f)

	collection.Remove(0)
	assert.Equal(t, 1, collection.Length())

	collection.RemoveByName("XYZ")
	assert.Equal(t, 0, collection.Length())
}
