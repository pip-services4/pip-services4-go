package test_calculator_functions

import (
	"testing"
	"time"

	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/calculator/functions"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/variants"
	"github.com/stretchr/testify/assert"
)

func TestDefaultFunctionsCollection(t *testing.T) {
	collection := functions.NewDefaultFunctionCollection()
	parameters := []*variants.Variant{
		variants.VariantFromInteger(1),
		variants.VariantFromInteger(2),
		variants.VariantFromInteger(3),
	}
	operations := variants.NewTypeUnsafeVariantOperations()

	f := collection.FindByName("sum")
	assert.NotNil(t, f)

	result, err := f.Calculate(parameters, operations)
	assert.Nil(t, err)
	assert.Equal(t, variants.Integer, result.Type())
	assert.Equal(t, 6, result.AsInteger())
}

func TestDefaultFunctionsCollectionDateFunctions(t *testing.T) {
	collection := functions.NewDefaultFunctionCollection()
	parameters := []*variants.Variant{}
	operations := variants.NewTypeUnsafeVariantOperations()

	f := collection.FindByName("now")
	assert.NotNil(t, f)

	result, err := f.Calculate(parameters, operations)
	assert.Nil(t, err)
	assert.Equal(t, variants.DateTime, result.Type())

	parameters = []*variants.Variant{
		variants.VariantFromInteger(1975),
		variants.VariantFromInteger(4),
		variants.VariantFromInteger(8),
	}

	f = collection.FindByName("date")
	assert.NotNil(t, f)

	result, err = f.Calculate(parameters, operations)
	assert.Nil(t, err)
	assert.Equal(t, variants.DateTime, result.Type())
	date := time.Date(1975, time.Month(4), 8, 0, 0, 0, 0, time.Local)
	assert.Equal(t, date, result.AsDateTime())
}
