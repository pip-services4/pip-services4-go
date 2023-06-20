package test_calculator

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/calculator"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/variants"
	"github.com/stretchr/testify/assert"
)

func TestExpressionCalculator(t *testing.T) {
	calculator := calculator.NewExpressionCalculator()

	err := calculator.SetExpression("2 + 2")
	assert.Nil(t, err)
	result, err1 := calculator.Evaluate()
	assert.Nil(t, err1)
	assert.Equal(t, variants.Integer, result.Type())
	assert.Equal(t, 4, result.AsInteger())

	err = calculator.SetExpression("A + b / (3 - Max(-123, 1)*2)")
	assert.Nil(t, err)
	assert.Equal(t, "A", calculator.DefaultVariables().FindByName("a").Name())
	assert.Equal(t, "b", calculator.DefaultVariables().FindByName("b").Name())
	calculator.DefaultVariables().FindByName("a").SetValue(variants.VariantFromString("xyz"))
	calculator.DefaultVariables().FindByName("b").SetValue(variants.VariantFromInteger(123))
	result, err1 = calculator.Evaluate()
	assert.Nil(t, err1)
	assert.Equal(t, variants.String, result.Type())
	assert.Equal(t, "xyz123", result.AsString())

	err = calculator.SetExpression("'abc'[1]")
	assert.Nil(t, err)
	result, err1 = calculator.Evaluate()
	assert.Nil(t, err1)
	assert.Equal(t, variants.String, result.Type())
	assert.Equal(t, "b", result.AsString())

	err = calculator.SetExpression("1 > 2")
	assert.Nil(t, err)
	result, err1 = calculator.Evaluate()
	assert.Nil(t, err1)
	assert.Equal(t, variants.Boolean, result.Type())
	assert.False(t, result.AsBoolean())

	err = calculator.SetExpression("2 IN ARRAY(1,2,3)")
	assert.Nil(t, err)
	result, err1 = calculator.Evaluate()
	assert.Nil(t, err1)
	assert.Equal(t, variants.Boolean, result.Type())
	assert.True(t, result.AsBoolean())

	err = calculator.SetExpression("5 NOT IN ARRAY(1,2,3)")
	assert.Nil(t, err)
	result, err1 = calculator.Evaluate()
	assert.Nil(t, err1)
	assert.Equal(t, variants.Boolean, result.Type())
	assert.True(t, result.AsBoolean())
}
