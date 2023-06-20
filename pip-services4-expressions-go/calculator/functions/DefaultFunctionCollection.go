package functions

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/calculator/errors"
	"github.com/pip-services4/pip-services4-go/pip-services4-expressions-go/variants"
)

// DefaultFunctionCollection implements a list filled with standard functions.
type DefaultFunctionCollection struct {
	*FunctionCollection
}

// NewDefaultFunctionCollection constructs this list and fills it with the standard functions.
func NewDefaultFunctionCollection() *DefaultFunctionCollection {
	c := &DefaultFunctionCollection{
		FunctionCollection: NewFunctionCollection(),
	}

	c.Add(NewDelegatedFunction("Ticks", ticksFunctionCalculator))
	c.Add(NewDelegatedFunction("TimeSpan", timeSpanFunctionCalculator))
	c.Add(NewDelegatedFunction("Now", nowFunctionCalculator))
	c.Add(NewDelegatedFunction("Date", dateFunctionCalculator))
	c.Add(NewDelegatedFunction("DayOfWeek", dayOfWeekFunctionCalculator))
	c.Add(NewDelegatedFunction("Min", minFunctionCalculator))
	c.Add(NewDelegatedFunction("Max", maxFunctionCalculator))
	c.Add(NewDelegatedFunction("Sum", sumFunctionCalculator))
	c.Add(NewDelegatedFunction("If", ifFunctionCalculator))
	c.Add(NewDelegatedFunction("Choose", chooseFunctionCalculator))
	c.Add(NewDelegatedFunction("E", eFunctionCalculator))
	c.Add(NewDelegatedFunction("Pi", piFunctionCalculator))
	c.Add(NewDelegatedFunction("Rnd", rndFunctionCalculator))
	c.Add(NewDelegatedFunction("Random", rndFunctionCalculator))
	c.Add(NewDelegatedFunction("Abs", absFunctionCalculator))
	c.Add(NewDelegatedFunction("Acos", acosFunctionCalculator))
	c.Add(NewDelegatedFunction("Asin", asinFunctionCalculator))
	c.Add(NewDelegatedFunction("Atan", atanFunctionCalculator))
	c.Add(NewDelegatedFunction("Exp", expFunctionCalculator))
	c.Add(NewDelegatedFunction("Log", logFunctionCalculator))
	c.Add(NewDelegatedFunction("Ln", logFunctionCalculator))
	c.Add(NewDelegatedFunction("Log10", log10FunctionCalculator))
	c.Add(NewDelegatedFunction("Ceil", ceilFunctionCalculator))
	c.Add(NewDelegatedFunction("Ceiling", ceilFunctionCalculator))
	c.Add(NewDelegatedFunction("Floor", floorFunctionCalculator))
	c.Add(NewDelegatedFunction("Round", roundFunctionCalculator))
	c.Add(NewDelegatedFunction("Trunc", truncFunctionCalculator))
	c.Add(NewDelegatedFunction("Truncate", truncFunctionCalculator))
	c.Add(NewDelegatedFunction("Cos", cosFunctionCalculator))
	c.Add(NewDelegatedFunction("Sin", sinFunctionCalculator))
	c.Add(NewDelegatedFunction("Tan", tanFunctionCalculator))
	c.Add(NewDelegatedFunction("Sqr", sqrtFunctionCalculator))
	c.Add(NewDelegatedFunction("Sqrt", sqrtFunctionCalculator))
	c.Add(NewDelegatedFunction("Empty", emptyFunctionCalculator))
	c.Add(NewDelegatedFunction("Null", nullFunctionCalculator))
	c.Add(NewDelegatedFunction("Contains", containsFunctionCalculator))
	c.Add(NewDelegatedFunction("Array", arrayFunctionCalculator))

	return c
}

// checkParamCount checks if parameters contains the correct number of function parameters
// (must be stored on the top of the parameters).
//
//	Parameters:
//		- parameters: A list with function parameters.
//		- expectedParamCount: The expected number of function parameters.
func checkParamCount(parameters []*variants.Variant, expectedParamCount int) error {
	paramCount := len(parameters)
	if expectedParamCount != paramCount {
		err := errors.NewExpressionError("", "WRONG_PARAM_COUNT",
			"Expected "+strconv.Itoa(expectedParamCount)+
				" parameters but was found "+strconv.Itoa(paramCount), 0, 0)
		return err
	}
	return nil
}

// getParameter gets function parameter by it's index.
//
//	Parameters:
//		- parameters: A list with function parameters.
//		- paramIndex: Index for the function parameter (0 for the first parameter).
//	Returns: Function parameter value.
func getParameter(parameters []*variants.Variant, paramIndex int) *variants.Variant {
	return parameters[paramIndex]
}

func ticksFunctionCalculator(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {
	err := checkParamCount(parameters, 0)
	if err != nil {
		return nil, err
	}

	result := variants.VariantFromLong(time.Now().Unix())

	return result, nil
}

func timeSpanFunctionCalculator(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {
	paramCount := len(parameters)
	if paramCount != 1 && paramCount != 3 && paramCount != 4 && paramCount != 5 {
		err := errors.NewExpressionError("", "WRONG_PARAM_COUNT", "Expected 1, 3, 4 or 5 parameters", 0, 0)
		return nil, err
	}

	result := variants.EmptyVariant()

	if paramCount == 1 {
		value, err := variantOperations.Convert(getParameter(parameters, 0), variants.Long)
		if err != nil {
			return nil, err
		}

		result.SetAsTimeSpan(time.Millisecond * time.Duration(value.AsLong()))
	} else if paramCount > 2 {
		value1, err1 := variantOperations.Convert(getParameter(parameters, 0), variants.Long)
		if err1 != nil {
			return nil, err1
		}

		value2, err2 := variantOperations.Convert(getParameter(parameters, 1), variants.Long)
		if err2 != nil {
			return nil, err2
		}

		value3, err3 := variantOperations.Convert(getParameter(parameters, 2), variants.Long)
		if err3 != nil {
			return nil, err3
		}

		value4 := variants.VariantFromLong(0)
		if paramCount > 3 {
			value4, err1 = variantOperations.Convert(getParameter(parameters, 3), variants.Long)
			if err1 != nil {
				return nil, err1
			}
		}

		value5 := variants.VariantFromLong(0)
		if paramCount > 4 {
			value5, err1 = variantOperations.Convert(getParameter(parameters, 4), variants.Long)
			if err1 != nil {
				return nil, err1
			}
		}

		ticks := (((value1.AsLong()*24+value2.AsLong())*60+value3.AsLong())*60+value4.AsLong())*1000 + value5.AsLong()
		result.SetAsTimeSpan(time.Millisecond * time.Duration(ticks))
	}

	return result, nil
}

func nowFunctionCalculator(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {
	err := checkParamCount(parameters, 0)
	if err != nil {
		return nil, err
	}

	result := variants.VariantFromDateTime(time.Now())

	return result, nil
}

func dateFunctionCalculator(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {

	paramCount := len(parameters)
	if paramCount < 1 || paramCount > 7 {
		err := errors.NewExpressionError("", "WRONG_PARAM_COUNT", "Expected from 1 to 7 parameters", 0, 0)
		return nil, err
	}

	if paramCount == 1 {
		value, err := variantOperations.Convert(getParameter(parameters, 0), variants.Long)
		if err != nil {
			return nil, err
		}
		date := time.Unix(value.AsLong(), 0)
		result := variants.VariantFromDateTime(date)
		return result, nil
	}

	value1, err := variantOperations.Convert(getParameter(parameters, 0), variants.Integer)
	if err != nil {
		return nil, err
	}

	value2 := variants.VariantFromInteger(1)
	if paramCount > 1 {
		value2, err = variantOperations.Convert(getParameter(parameters, 1), variants.Integer)
		if err != nil {
			return nil, err
		}
	}

	value3 := variants.VariantFromInteger(1)
	if paramCount > 2 {
		value3, err = variantOperations.Convert(getParameter(parameters, 2), variants.Integer)
		if err != nil {
			return nil, err
		}
	}

	value4 := variants.VariantFromInteger(0)
	if paramCount > 3 {
		value4, err = variantOperations.Convert(getParameter(parameters, 3), variants.Integer)
		if err != nil {
			return nil, err
		}
	}

	value5 := variants.VariantFromInteger(0)
	if paramCount > 4 {
		value5, err = variantOperations.Convert(getParameter(parameters, 4), variants.Integer)
		if err != nil {
			return nil, err
		}
	}

	value6 := variants.VariantFromInteger(0)
	if paramCount > 5 {
		value6, err = variantOperations.Convert(getParameter(parameters, 5), variants.Integer)
		if err != nil {
			return nil, err
		}
	}

	value7 := variants.VariantFromInteger(0)
	if paramCount > 6 {
		value7, err = variantOperations.Convert(getParameter(parameters, 6), variants.Integer)
		if err != nil {
			return nil, err
		}
	}

	date := time.Date(value1.AsInteger(), time.Month(value2.AsInteger()), value3.AsInteger(),
		value4.AsInteger(), value5.AsInteger(), value6.AsInteger(), value7.AsInteger(), time.Local)
	result := variants.VariantFromDateTime(date)
	return result, nil
}

func dayOfWeekFunctionCalculator(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {
	err := checkParamCount(parameters, 1)
	if err != nil {
		return nil, err
	}

	value, err1 := variantOperations.Convert(getParameter(parameters, 0), variants.DateTime)
	if err1 != nil {
		return nil, err1
	}

	day := value.AsDateTime().Weekday()
	result := variants.VariantFromInteger(int(day))

	return result, nil
}

func minFunctionCalculator(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {
	paramCount := len(parameters)
	if paramCount < 2 {
		err := errors.NewExpressionError("", "WRONG_PARAM_COUNT",
			"Expected at least 2 parameters", 0, 0)
		return nil, err
	}

	result := getParameter(parameters, 0)
	for i := 1; i < paramCount; i = i + 1 {
		value := getParameter(parameters, i)
		temp, err := variantOperations.More(result, value)
		if err != nil {
			return nil, err
		}
		if temp.AsBoolean() {
			result = value
		}
	}

	return result, nil
}

func maxFunctionCalculator(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {
	paramCount := len(parameters)
	if paramCount < 2 {
		err := errors.NewExpressionError("", "WRONG_PARAM_COUNT",
			"Expected at least 2 parameters", 0, 0)
		return nil, err
	}

	result := getParameter(parameters, 0)
	for i := 1; i < paramCount; i++ {
		value := getParameter(parameters, i)
		temp, err := variantOperations.Less(result, value)
		if err != nil {
			return nil, err
		}
		if temp.AsBoolean() {
			result = value
		}
	}

	return result, nil
}

func sumFunctionCalculator(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {
	paramCount := len(parameters)
	if paramCount < 2 {
		err := errors.NewExpressionError("", "WRONG_PARAM_COUNT",
			"Expected at least 2 parameters", 0, 0)
		return nil, err
	}

	result := getParameter(parameters, 0)
	for i := 1; i < paramCount; i++ {
		value := getParameter(parameters, i)
		var err error
		result, err = variantOperations.Add(result, value)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func ifFunctionCalculator(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {
	err := checkParamCount(parameters, 3)
	if err != nil {
		return nil, err
	}

	value1 := getParameter(parameters, 0)
	value2 := getParameter(parameters, 1)
	value3 := getParameter(parameters, 2)

	condition, err1 := variantOperations.Convert(value1, variants.Boolean)
	if err1 != nil {
		return nil, err1
	}

	result := value3
	if condition.AsBoolean() {
		result = value2
	}

	return result, nil
}

func chooseFunctionCalculator(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {
	paramCount := len(parameters)
	if paramCount < 3 {
		err := errors.NewExpressionError("", "WRONG_PARAM_COUNT",
			"Expected at least 3 parameters", 0, 0)
		return nil, err
	}

	value1 := getParameter(parameters, 0)
	condition, err := variantOperations.Convert(value1, variants.Integer)
	if err != nil {
		return nil, err
	}
	paramIndex := int(condition.AsInteger())

	if paramCount < paramIndex+1 {
		err := errors.NewExpressionError("", "WRONG_PARAM_COUNT",
			"Expected at least "+strconv.Itoa(paramIndex+1)+" parameters", 0, 0)
		return nil, err
	}

	result := getParameter(parameters, paramIndex)

	return result, nil
}

func eFunctionCalculator(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {
	err := checkParamCount(parameters, 0)
	if err != nil {
		return nil, err
	}

	result := variants.VariantFromFloat(math.E)

	return result, nil
}

func piFunctionCalculator(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {
	err := checkParamCount(parameters, 0)
	if err != nil {
		return nil, err
	}

	result := variants.VariantFromFloat(math.Pi)

	return result, nil
}

func rndFunctionCalculator(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {
	err := checkParamCount(parameters, 0)
	if err != nil {
		return nil, err
	}

	result := variants.VariantFromFloat(rand.Float32())

	return result, nil
}

func absFunctionCalculator(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {
	err := checkParamCount(parameters, 1)
	if err != nil {
		return nil, err
	}

	value := getParameter(parameters, 0)
	result := variants.EmptyVariant()
	switch value.Type() {
	case variants.Integer:
		result.SetAsInteger(int(math.Abs(float64(value.AsInteger()))))
		break
	case variants.Long:
		result.SetAsLong(int64(math.Abs(float64(value.AsLong()))))
		break
	case variants.Float:
		result.SetAsFloat(float32(math.Abs(float64(value.AsFloat()))))
		break
	case variants.Double:
		result.SetAsDouble(math.Abs(value.AsDouble()))
		break
	default:
		value, err = variantOperations.Convert(value, variants.Double)
		if err != nil {
			return nil, err
		}
		result.SetAsDouble(math.Abs(value.AsDouble()))
		break
	}

	return result, nil
}

func acosFunctionCalculator(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {
	err := checkParamCount(parameters, 1)
	if err != nil {
		return nil, err
	}

	value, err := variantOperations.Convert(getParameter(parameters, 0), variants.Double)
	if err == nil {
		return nil, err
	}
	result := variants.VariantFromDouble(math.Acos(value.AsDouble()))

	return result, nil
}

func asinFunctionCalculator(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {
	err := checkParamCount(parameters, 1)
	if err != nil {
		return nil, err
	}

	value, err1 := variantOperations.Convert(getParameter(parameters, 0), variants.Double)
	if err1 != nil {
		return nil, err1
	}
	result := variants.VariantFromDouble(math.Asin(value.AsDouble()))

	return result, nil
}

func atanFunctionCalculator(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {
	err := checkParamCount(parameters, 1)
	if err != nil {
		return nil, err
	}

	value, err1 := variantOperations.Convert(getParameter(parameters, 0), variants.Double)
	if err1 != nil {
		return nil, err1
	}
	result := variants.VariantFromDouble(math.Atan(value.AsDouble()))

	return result, nil
}

func expFunctionCalculator(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {
	err := checkParamCount(parameters, 1)
	if err != nil {
		return nil, err
	}

	value, err1 := variantOperations.Convert(getParameter(parameters, 0), variants.Double)
	if err1 != nil {
		return nil, err1
	}
	result := variants.VariantFromDouble(math.Exp(value.AsDouble()))

	return result, nil
}

func logFunctionCalculator(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {
	err := checkParamCount(parameters, 1)
	if err != nil {
		return nil, err
	}

	value, err1 := variantOperations.Convert(getParameter(parameters, 0), variants.Double)
	if err1 != nil {
		return nil, err1
	}
	result := variants.VariantFromDouble(math.Log(value.AsDouble()))

	return result, nil
}

func log10FunctionCalculator(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {
	err := checkParamCount(parameters, 1)
	if err != nil {
		return nil, err
	}

	value, err1 := variantOperations.Convert(getParameter(parameters, 0), variants.Double)
	if err1 != nil {
		return nil, err1
	}
	result := variants.VariantFromDouble(math.Log10(value.AsDouble()))

	return result, nil
}

func ceilFunctionCalculator(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {
	err := checkParamCount(parameters, 1)
	if err != nil {
		return nil, err
	}

	value, err1 := variantOperations.Convert(getParameter(parameters, 0), variants.Double)
	if err1 != nil {
		return nil, err1
	}
	result := variants.VariantFromDouble(math.Ceil(value.AsDouble()))

	return result, nil
}

func floorFunctionCalculator(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {
	err := checkParamCount(parameters, 1)
	if err != nil {
		return nil, err
	}

	value, err1 := variantOperations.Convert(getParameter(parameters, 0), variants.Double)
	if err1 != nil {
		return nil, err1
	}
	result := variants.VariantFromDouble(math.Floor(value.AsDouble()))

	return result, nil
}

func roundFunctionCalculator(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {
	err := checkParamCount(parameters, 1)
	if err != nil {
		return nil, err
	}

	value, err1 := variantOperations.Convert(getParameter(parameters, 0), variants.Double)
	if err1 != nil {
		return nil, err1
	}
	result := variants.VariantFromDouble(math.Round(value.AsDouble()))

	return result, nil
}

func truncFunctionCalculator(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {
	err := checkParamCount(parameters, 1)
	if err != nil {
		return nil, err
	}

	value, err1 := variantOperations.Convert(getParameter(parameters, 0), variants.Double)
	if err1 != nil {
		return nil, err1
	}
	result := variants.VariantFromLong(int64(math.Trunc(value.AsDouble())))

	return result, nil
}

func cosFunctionCalculator(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {
	err := checkParamCount(parameters, 1)
	if err != nil {
		return nil, err
	}

	value, err1 := variantOperations.Convert(getParameter(parameters, 0), variants.Double)
	if err1 != nil {
		return nil, err1
	}
	result := variants.VariantFromDouble(math.Cos(value.AsDouble()))

	return result, nil
}

func sinFunctionCalculator(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {
	err := checkParamCount(parameters, 1)
	if err != nil {
		return nil, err
	}

	value, err1 := variantOperations.Convert(getParameter(parameters, 0), variants.Double)
	if err1 != nil {
		return nil, err1
	}
	result := variants.VariantFromDouble(math.Sin(value.AsDouble()))

	return result, nil
}

func tanFunctionCalculator(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {
	err := checkParamCount(parameters, 1)
	if err != nil {
		return nil, err
	}

	value, err1 := variantOperations.Convert(getParameter(parameters, 0), variants.Double)
	if err1 != nil {
		return nil, err1
	}
	result := variants.VariantFromDouble(math.Tan(value.AsDouble()))

	return result, nil
}

func sqrtFunctionCalculator(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {
	err := checkParamCount(parameters, 1)
	if err != nil {
		return nil, err
	}

	value, err1 := variantOperations.Convert(getParameter(parameters, 0), variants.Double)
	if err1 != nil {
		return nil, err1
	}
	result := variants.VariantFromDouble(math.Sqrt(value.AsDouble()))

	return result, nil
}

func emptyFunctionCalculator(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {
	err := checkParamCount(parameters, 1)
	if err != nil {
		return nil, err
	}

	value := getParameter(parameters, 0)
	result := variants.VariantFromBoolean(value.IsEmpty())

	return result, nil
}

func nullFunctionCalculator(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {
	err := checkParamCount(parameters, 0)
	if err != nil {
		return nil, err
	}

	result := variants.EmptyVariant()

	return result, nil
}

func containsFunctionCalculator(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {
	err := checkParamCount(parameters, 2)
	if err != nil {
		return nil, err
	}

	str, err1 := variantOperations.Convert(getParameter(parameters, 0), variants.String)
	if err1 != nil {
		return nil, err1
	}
	substr, err2 := variantOperations.Convert(getParameter(parameters, 1), variants.String)
	if err2 != nil {
		return nil, err2
	}

	if str.IsEmpty() || str.IsNull() {
		return variants.VariantFromBoolean(false), nil
	}

	contains := strings.Contains(str.AsString(), substr.AsString())
	result := variants.VariantFromBoolean(contains)

	return result, nil
}

func arrayFunctionCalculator(parameters []*variants.Variant,
	variantOperations variants.IVariantOperations) (*variants.Variant, error) {
	result := variants.VariantFromArray(parameters)
	return result, nil
}
