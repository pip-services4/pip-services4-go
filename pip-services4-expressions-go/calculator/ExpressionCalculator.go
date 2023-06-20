package calculator

import (
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/calculator/errors"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/calculator/functions"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/calculator/parsers"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/calculator/variables"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/tokenizers"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/variants"
)

// ExpressionCalculator implements an expression calculator class.
type ExpressionCalculator struct {
	defaultVariables  variables.IVariableCollection
	defaultFunctions  functions.IFunctionCollection
	variantOperations variants.IVariantOperations
	parser            *parsers.ExpressionParser
	autoVariables     bool
}

// NewExpressionCalculator constructs this class with default parameters.
func NewExpressionCalculator() *ExpressionCalculator {
	c := &ExpressionCalculator{
		defaultVariables:  variables.NewVariableCollection(),
		defaultFunctions:  functions.NewDefaultFunctionCollection(),
		variantOperations: variants.NewTypeUnsafeVariantOperations(),
		parser:            parsers.NewExpressionParser(),
		autoVariables:     true,
	}
	return c
}

// ExpressionCalculatorFromExpression constructs this class and assigns expression string.
//
//	Parameters:
//		- expression: The expression string.
func ExpressionCalculatorFromExpression(expression string) (*ExpressionCalculator, error) {
	c := NewExpressionCalculator()
	err := c.SetExpression(expression)
	return c, err
}

func ExpressionCalculatorFromTokens(originalTokens []*tokenizers.Token) *ExpressionCalculator {
	c := NewExpressionCalculator()
	c.SetOriginalTokens(originalTokens)
	return c
}

// Expression gets the expression string.
func (c *ExpressionCalculator) Expression() string {
	return c.parser.Expression()
}

// SetExpression sets the expression string.
func (c *ExpressionCalculator) SetExpression(value string) error {
	err := c.parser.SetExpression(value)
	if err != nil {
		return err
	}

	if c.autoVariables {
		c.CreateVariables(c.defaultVariables)
	}

	return nil
}

func (c *ExpressionCalculator) OriginalTokens() []*tokenizers.Token {
	return c.parser.OriginalTokens()
}

func (c *ExpressionCalculator) SetOriginalTokens(value []*tokenizers.Token) {
	c.parser.SetOriginalTokens(value)
	if c.autoVariables {
		c.CreateVariables(c.defaultVariables)
	}
}

// AutoVariables gets the flag to turn on auto creation of variables for specified expression.
func (c *ExpressionCalculator) AutoVariables() bool {
	return c.autoVariables
}

// SetAutoVariables sets the flag to turn on auto creation of variables for specified expression.
func (c *ExpressionCalculator) SetAutoVariables(value bool) {
	c.autoVariables = value
}

// VariantOperations gets the manager for operations on variant values.
func (c *ExpressionCalculator) VariantOperations() variants.IVariantOperations {
	return c.variantOperations
}

// SetVariantOperations sets the manager for operations on variant values.
func (c *ExpressionCalculator) SetVariantOperations(value variants.IVariantOperations) {
	c.variantOperations = value
}

// DefaultVariables the list with default variables.
func (c *ExpressionCalculator) DefaultVariables() variables.IVariableCollection {
	return c.defaultVariables
}

// DefaultFunctions the list with default functions.
func (c *ExpressionCalculator) DefaultFunctions() functions.IFunctionCollection {
	return c.defaultFunctions
}

// InitialTokens the list of original expression tokens.
func (c *ExpressionCalculator) InitialTokens() []*parsers.ExpressionToken {
	return c.parser.InitialTokens()
}

// ResultTokens the list of processed expression tokens.
func (c *ExpressionCalculator) ResultTokens() []*parsers.ExpressionToken {
	return c.parser.ResultTokens()
}

// CreateVariables populates the specified variables list with variables from parsed expression.
//
//	Parameters:
//		- variables: The list of variables to be populated.
func (c *ExpressionCalculator) CreateVariables(vars variables.IVariableCollection) {
	for _, variableName := range c.parser.VariableNames() {
		if vars.FindByName(variableName) == nil {
			vars.Add(variables.EmptyVariable(variableName))
		}
	}
}

// Clear cleans up this calculator from all data.
func (c *ExpressionCalculator) Clear() {
	c.parser.Clear()
	c.defaultVariables.Clear()
}

// Evaluate this expression using default variables and functions.
//
//	Returns: An evaluated expression value.
func (c *ExpressionCalculator) Evaluate() (*variants.Variant, error) {
	return c.EvaluateUsingVariablesAndFunctions(nil, nil)
}

// EvaluateUsingVariables evaluates this expression using specified variables.
//
//	Parameters:
//		- variables: The list of variables
//	Returns: An evaluated expression value.
func (c *ExpressionCalculator) EvaluateUsingVariables(
	vars variables.IVariableCollection) (*variants.Variant, error) {
	return c.EvaluateUsingVariablesAndFunctions(vars, nil)
}

// EvaluateUsingVariablesAndFunctions evaluates this expression using specified variables and functions.
//
//	Parameters:
//		- variables: The list of variables
//		- functions: The list of functions.
//	Returns: An evaluated expression value.
func (c *ExpressionCalculator) EvaluateUsingVariablesAndFunctions(
	vars variables.IVariableCollection, funcs functions.IFunctionCollection) (*variants.Variant, error) {

	stack := NewCalculationStack()
	if vars == nil {
		vars = c.defaultVariables
	}
	if funcs == nil {
		funcs = c.defaultFunctions
	}

	for _, token := range c.ResultTokens() {
		if ok, err := c.evaluateConstant(token, stack); ok || err != nil {
			if err != nil {
				return nil, err
			}
		} else if ok, err := c.evaluateVariable(token, stack, vars); ok || err != nil {
			if err != nil {
				return nil, err
			}
		} else if ok, err := c.evaluateFunction(token, stack, funcs); ok || err != nil {
			if err != nil {
				return nil, err
			}
		} else if ok, err := c.evaluateLogical(token, stack); ok || err != nil {
			if err != nil {
				return nil, err
			}
		} else if ok, err := c.evaluateArithmetical(token, stack); ok || err != nil {
			if err != nil {
				return nil, err
			}
		} else if ok, err := c.evaluateBoolean(token, stack); ok || err != nil {
			if err != nil {
				return nil, err
			}
		} else if ok, err := c.evaluateOther(token, stack); ok || err != nil {
			if err != nil {
				return nil, err
			}
		} else {
			err := errors.NewExpressionError("", "INTERNAL", "Internal error", token.Line(), token.Column())
			return nil, err
		}
	}

	if stack.Length() != 1 {
		err := errors.NewExpressionError("", "INTERNAL", "Internal error", 0, 0)
		return nil, err
	}

	return stack.Pop(), nil
}

func (c *ExpressionCalculator) evaluateConstant(
	token *parsers.ExpressionToken, stack *CalculationStack) (bool, error) {
	if token.Type() == parsers.Constant {
		stack.Push(token.Value())
		return true, nil
	}
	return false, nil
}

func (c *ExpressionCalculator) evaluateVariable(
	token *parsers.ExpressionToken, stack *CalculationStack,
	vars variables.IVariableCollection) (bool, error) {

	if token.Type() == parsers.Variable {
		variable := vars.FindByName(token.Value().AsString())
		if variable == nil {
			err := errors.NewExpressionError("", "VAR_NOT_FOUND",
				"Variable "+token.Value().AsString()+" was not found.",
				token.Line(), token.Column())
			return false, err
		}

		stack.Push(variable.Value())
		return true, nil
	}

	return false, nil
}

func (c *ExpressionCalculator) evaluateFunction(
	token *parsers.ExpressionToken, stack *CalculationStack,
	funcs functions.IFunctionCollection) (bool, error) {

	if token.Type() == parsers.Function {
		function := funcs.FindByName(token.Value().AsString())
		if function == nil {
			err := errors.NewExpressionError("", "FUNC_NOT_FOUND",
				"Function "+token.Value().AsString()+" was not found.",
				token.Line(), token.Column())
			return false, err
		}

		// Prepare parameters
		parameters := []*variants.Variant{}
		paramCount := stack.Pop().AsInteger()
		for paramCount > 0 {
			parameters = append([]*variants.Variant{stack.Pop()}, parameters...)
			paramCount = paramCount - 1
		}

		functionResult, err := function.Calculate(parameters, c.variantOperations)
		if err != nil {
			return false, err
		}

		stack.Push(functionResult)

		return true, nil
	}

	return false, nil
}

func (c *ExpressionCalculator) evaluateLogical(
	token *parsers.ExpressionToken, stack *CalculationStack) (bool, error) {

	switch token.Type() {
	case parsers.And:
		{
			value2 := stack.Pop()
			value1 := stack.Pop()
			result, err := c.variantOperations.And(value1, value2)
			if err != nil {
				return false, err
			}
			stack.Push(result)
			return true, nil
		}
	case parsers.Or:
		{
			value2 := stack.Pop()
			value1 := stack.Pop()
			result, err := c.variantOperations.Or(value1, value2)
			if err != nil {
				return false, err
			}
			stack.Push(result)
			return true, nil
		}
	case parsers.Xor:
		{
			value2 := stack.Pop()
			value1 := stack.Pop()
			result, err := c.variantOperations.Xor(value1, value2)
			if err != nil {
				return false, err
			}
			stack.Push(result)
			return true, nil
		}
	case parsers.Not:
		{
			value := stack.Pop()
			result, err := c.variantOperations.Not(value)
			if err != nil {
				return false, err
			}
			stack.Push(result)
			return true, nil
		}
	}

	return false, nil
}

func (c *ExpressionCalculator) evaluateArithmetical(
	token *parsers.ExpressionToken, stack *CalculationStack) (bool, error) {

	switch token.Type() {
	case parsers.Plus:
		{
			value2 := stack.Pop()
			value1 := stack.Pop()
			result, err := c.variantOperations.Add(value1, value2)
			if err != nil {
				return false, err
			}
			stack.Push(result)
			return true, nil
		}
	case parsers.Minus:
		{
			value2 := stack.Pop()
			value1 := stack.Pop()
			result, err := c.variantOperations.Sub(value1, value2)
			if err != nil {
				return false, err
			}
			stack.Push(result)
			return true, nil
		}
	case parsers.Star:
		{
			value2 := stack.Pop()
			value1 := stack.Pop()
			result, err := c.variantOperations.Mul(value1, value2)
			if err != nil {
				return false, err
			}
			stack.Push(result)
			return true, nil
		}
	case parsers.Slash:
		{
			value2 := stack.Pop()
			value1 := stack.Pop()
			result, err := c.variantOperations.Div(value1, value2)
			if err != nil {
				return false, err
			}
			stack.Push(result)
			return true, nil
		}
	case parsers.Procent:
		{
			value2 := stack.Pop()
			value1 := stack.Pop()
			result, err := c.variantOperations.Mod(value1, value2)
			if err != nil {
				return false, err
			}
			stack.Push(result)
			return true, nil
		}
	case parsers.Power:
		{
			value2 := stack.Pop()
			value1 := stack.Pop()
			result, err := c.variantOperations.Pow(value1, value2)
			if err != nil {
				return false, err
			}
			stack.Push(result)
			return true, nil
		}
	case parsers.Unary:
		{
			value := stack.Pop()
			result, err := c.variantOperations.Negative(value)
			if err != nil {
				return false, err
			}
			stack.Push(result)
			return true, nil
		}
	case parsers.ShiftLeft:
		{
			value2 := stack.Pop()
			value1 := stack.Pop()
			result, err := c.variantOperations.Lsh(value1, value2)
			if err != nil {
				return false, err
			}
			stack.Push(result)
			return true, nil
		}
	case parsers.ShiftRight:
		{
			value2 := stack.Pop()
			value1 := stack.Pop()
			result, err := c.variantOperations.Rsh(value1, value2)
			if err != nil {
				return false, err
			}
			stack.Push(result)
			return true, nil
		}
	}

	return false, nil
}

func (c *ExpressionCalculator) evaluateBoolean(
	token *parsers.ExpressionToken, stack *CalculationStack) (bool, error) {

	switch token.Type() {
	case parsers.Equal:
		{
			value2 := stack.Pop()
			value1 := stack.Pop()
			result, err := c.variantOperations.Equal(value1, value2)
			if err != nil {
				return false, err
			}
			stack.Push(result)
			return true, nil
		}
	case parsers.NotEqual:
		{
			value2 := stack.Pop()
			value1 := stack.Pop()
			result, err := c.variantOperations.NotEqual(value1, value2)
			if err != nil {
				return false, err
			}
			stack.Push(result)
			return true, nil
		}
	case parsers.More:
		{
			value2 := stack.Pop()
			value1 := stack.Pop()
			result, err := c.variantOperations.More(value1, value2)
			if err != nil {
				return false, err
			}
			stack.Push(result)
			return true, nil
		}
	case parsers.Less:
		{
			value2 := stack.Pop()
			value1 := stack.Pop()
			result, err := c.variantOperations.Less(value1, value2)
			if err != nil {
				return false, err
			}
			stack.Push(result)
			return true, nil
		}
	case parsers.EqualMore:
		{
			value2 := stack.Pop()
			value1 := stack.Pop()
			result, err := c.variantOperations.MoreEqual(value1, value2)
			if err != nil {
				return false, err
			}
			stack.Push(result)
			return true, nil
		}
	case parsers.EqualLess:
		{
			value2 := stack.Pop()
			value1 := stack.Pop()
			result, err := c.variantOperations.LessEqual(value1, value2)
			if err != nil {
				return false, err
			}
			stack.Push(result)
			return true, nil
		}
	}

	return false, nil
}

func (c *ExpressionCalculator) evaluateOther(
	token *parsers.ExpressionToken, stack *CalculationStack) (bool, error) {

	switch token.Type() {
	case parsers.In:
		{
			value2 := stack.Pop()
			value1 := stack.Pop()
			result, err := c.variantOperations.In(value2, value1)
			if err != nil {
				return false, err
			}
			stack.Push(result)
			return true, nil
		}
	case parsers.NotIn:
		{
			value2 := stack.Pop()
			value1 := stack.Pop()
			result, err := c.variantOperations.In(value2, value1)
			if err != nil {
				return false, err
			}
			result = variants.VariantFromBoolean(!result.AsBoolean())
			stack.Push(result)
			return true, nil
		}
	case parsers.Element:
		{
			value2 := stack.Pop()
			value1 := stack.Pop()
			result, err := c.variantOperations.GetElement(value1, value2)
			if err != nil {
				return false, err
			}
			stack.Push(result)
			return true, nil
		}
	case parsers.IsNull:
		{
			stack.Push(variants.VariantFromBoolean(stack.Pop().IsNull()))
			return true, nil
		}
	case parsers.IsNotNull:
		{
			stack.Push(variants.VariantFromBoolean(!stack.Pop().IsNull()))
			return true, nil
		}
	}

	return false, nil
}
