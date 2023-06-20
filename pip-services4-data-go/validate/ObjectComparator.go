package validate

import (
	"regexp"
	"strings"

	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/data"
)

// ObjectComparator Helper class to perform comparison operations over arbitrary values.
//
//	Example:
//		ObjectComparator.Compare(2, "GT", 1);        // Result: true
//		ObjectComparator.AreEqual("A", "B");         // Result: false
var ObjectComparator = &_TObjectComparator{}

type _TObjectComparator struct{}

// Compare perform comparison operation over two arguments.
// The operation can be performed over values of any type.
//
//	Parameters:
//		- value1 any the first argument to compare
//		- operation string the comparison operation:
//			"==" ("=", "EQ"), "!= " ("<>", "NE"); "<"/">" ("LT"/"GT"), "<="/">=" ("LE"/"GE"); "LIKE".
//		- value2 any the second argument to compare
//	Returns: bool result of the comparison operation
func (c *_TObjectComparator) Compare(value1 any, operation string, value2 any) bool {
	operation = strings.ToUpper(operation)

	if operation == "=" || operation == "==" || operation == "EQ" {
		return c.AreEqual(value1, value2)
	}
	if operation == "!=" || operation == "<>" || operation == "NE" {
		return c.AreNotEqual(value1, value2)
	}
	if operation == "<" || operation == "LT" {
		return c.IsLess(value1, value2)
	}
	if operation == "<=" || operation == "LE" || operation == "LTE" {
		return c.AreEqual(value1, value2) || c.IsLess(value1, value2)
	}
	if operation == ">" || operation == "GT" {
		return c.IsGreater(value1, value2)
	}
	if operation == ">=" || operation == "GE" || operation == "GTE" {
		return c.AreEqual(value1, value2) || c.IsGreater(value1, value2)
	}
	if operation == "LIKE" {
		return c.Match(value1, value2)
	}

	return false
}

// AreEqual checks if two values are equal. The operation can be performed over values of any type.
//
//	Parameters:
//		- value1 interface the first value to compare
//		- value2 any the second value to compare
//	Returns: bool true if values are equal and false otherwise
func (c *_TObjectComparator) AreEqual(value1 any, value2 any) bool {
	if value1 == nil && value2 == nil {
		return true
	}
	if value1 == nil || value2 == nil {
		return false
	}

	if equatable, ok := value1.(data.IEquatable[any]); ok {
		return equatable.Equals(value2)
	}
	if equatable, ok := value2.(data.IEquatable[any]); ok {
		return equatable.Equals(value1)
	}

	if number1, ok := convert.DoubleConverter.ToNullableDouble(value1); ok {
		if number2, ok := convert.DoubleConverter.ToNullableDouble(value2); ok {
			return number1 == number2
		}
	}

	if str1, ok := convert.StringConverter.ToNullableString(value1); ok {
		if str2, ok := convert.StringConverter.ToNullableString(value2); ok {
			return str1 == str2
		}
	}

	return value1 == value2
}

// AreNotEqual checks if two values are NOT equal The operation can be performed over values of any type.
//
//	Parameters:
//		- value1 any the first value to compare
//		- value2 any the second value to compare
//	Returns: bool true if values are NOT equal and false otherwise
func (c *_TObjectComparator) AreNotEqual(value1 any, value2 any) bool {
	return !c.AreEqual(value1, value2)
}

// IsLess checks if first value is less than the second one. The operation can be performed over numbers or strings.
//
//	Parameters:
//		- value1 any the first value to compare
//		- value2 any the second value to compare
//	Returns: bool true if the first value is less than second and false otherwise.
func (c *_TObjectComparator) IsLess(value1 any, value2 any) bool {
	if number1, ok := convert.DoubleConverter.ToNullableDouble(value1); ok {
		if number2, ok := convert.DoubleConverter.ToNullableDouble(value2); ok {
			return number1 < number2
		}
	}
	return false
}

// IsGreater checks if first value is greater than the second one. The operation can be performed over numbers or strings.
//
//	Parameters:
//		- value1 any the first value to compare
//		- value2 any the second value to compare
//	Returns: bool true if the first value is greater than second and false otherwise.
func (c *_TObjectComparator) IsGreater(value1 any, value2 any) bool {
	if number1, ok := convert.DoubleConverter.ToNullableDouble(value1); ok {
		if number2, ok := convert.DoubleConverter.ToNullableDouble(value2); ok {
			return number1 > number2
		}
	}
	return false
}

// Match checks if string  views are matches
//
//	Parameters:
//		- value1 any a string value to match
//		- value1 any a string value to match
//	Returns: bool true if the value matches regular expression and false otherwise.
func (c *_TObjectComparator) Match(value1 any, value2 any) bool {
	if value1 == nil && value2 == nil {
		return true
	}
	if value1 == nil || value2 == nil {
		return false
	}

	str1 := convert.StringConverter.ToString(value1)
	str2 := convert.StringConverter.ToString(value2)

	matched, _ := regexp.MatchString(str2, str1)
	return matched
}
