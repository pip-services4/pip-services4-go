package reflect

import (
	refl "reflect"
	"strings"
	"time"

	"github.com/pip-services4/pip-services4-commons-go/convert"
)

// TypeMatcher Helper class matches value types for equality.
// This class has symmetric implementation across all languages supported by Pip.Services toolkit
// and used to support dynamic data processing.
var TypeMatcher = &_TTypeMatcher{}

type _TTypeMatcher struct{}

// MatchValue matches expected type to a type of a value. The expected type can be specified by a type, type name or TypeCode.
//
//	Parameters:
//		- expectedType any an expected type to match.
//		- actualValue any a value to match its type to the expected one.
//	Returns: bool true if types are matching and false if they don't.
func (c *_TTypeMatcher) MatchValue(expectedType any, actualValue any) bool {
	if expectedType == nil {
		return true
	}
	if actualValue == nil {
		panic("Actual value cannot be nil")
	}

	// Check actual value by its type
	return c.MatchType(expectedType, refl.TypeOf(actualValue))
}

// MatchValueByName matches expected type to a type of a value.
//
//	Parameters:
//		- expectedType string an expected type name to match.
//		- actualValue any a value to match its type to the expected one.
//	Returns: bool true if types are matching and false if they don't.
func (c *_TTypeMatcher) MatchValueByName(expectedType string, actualValue any) bool {
	if expectedType == "" {
		return true
	}
	if actualValue == nil {
		panic("Actual value cannot be nil")
	}

	// Check actual value by its type
	return c.MatchTypeByName(expectedType, refl.TypeOf(actualValue))
}

// MatchType matches expected type to an actual type. The types can be specified as types, type names or TypeCode.
//
//	Parameters:
//		- expectedType any an expected type to match.
//		- actualType refl.Type n actual type to match.
//	Returns: bool true if types are matching and false if they don't.
func (c *_TTypeMatcher) MatchType(expectedType any, actualType refl.Type) bool {
	if expectedType == nil {
		return true
	}
	if actualType == nil {
		panic("Actual type cannot be null")
	}

	// Compare for matching types
	if expectedType == actualType {
		return true
	}

	// Extract inner value because Go implementations of Maps and Arrays are wrappers
	if innerType, ok := expectedType.(IValueWrapper); ok {
		expectedType = innerType.InnerValue()
	}

	// If expected value is type
	if typ, ok := expectedType.(refl.Type); ok {
		type2 := typ
		// Check pointer type as well
		if type2.Kind() == refl.Ptr {
			type2 = type2.Elem()
		}
		return actualType.AssignableTo(typ) || actualType.AssignableTo(type2)
	}

	// For strings compare string types
	if str, ok := expectedType.(string); ok {
		return c.MatchTypeByName(str, actualType)
	}

	// For typecodes compare them
	if typeCode, ok := expectedType.(convert.TypeCode); ok {
		//return convert.TypeConverter.ToTypeCode(actualType) == typeCode
		actualTypeCode := convert.TypeConverter.ToTypeCode(actualType)
		if typeCode == actualTypeCode {
			return true
		}
		// Special provisions for dynamic data
		if typeCode == convert.Integer &&
			(actualTypeCode == convert.Long || actualTypeCode == convert.Float || actualTypeCode == convert.Double) {
			return true
		}
		if typeCode == convert.Long &&
			(actualTypeCode == convert.Integer || actualTypeCode == convert.Float || actualTypeCode == convert.Double) {
			return true
		}
		if typeCode == convert.Float &&
			(actualTypeCode == convert.Integer || actualTypeCode == convert.Long || actualTypeCode == convert.Double) {
			return true
		}
		if typeCode == convert.Double &&
			(actualTypeCode == convert.Integer || actualTypeCode == convert.Long || actualTypeCode == convert.Float) {
			return true
		}
		if typeCode == convert.DateTime && actualTypeCode == convert.String {
			return true
		}
	}

	return false
}

// MatchTypeByName matches expected type to an actual type.
//
//	Parameters:
//		- expectedType string an expected type name to match.
//		- actualType refl.Type an actual type to match defined by type code.
//	Returns: bool true if types are matching and false if they don't.
func (c *_TTypeMatcher) MatchTypeByName(expectedType string, actualType refl.Type) bool {
	if expectedType == "" {
		return true
	}

	if actualType == nil {
		panic("Actual type cannot be null")
	}

	if actualType.Kind() == refl.Ptr {
		actualType = actualType.Elem()
	}

	expectedType = strings.ToLower(expectedType)
	actualTypeName := strings.ToLower(actualType.Name())
	actualTypeFullName := strings.ToLower(actualType.PkgPath() + "." + actualType.Name())
	actualTypeKind := actualType.Kind()

	if actualTypeName == expectedType || actualTypeFullName == expectedType {
		return true
	}

	if expectedType == "object" {
		return true
	}

	if expectedType == "int" || expectedType == "integer" {
		return actualTypeKind == refl.Int8 ||
			actualTypeKind == refl.Uint8 ||
			actualTypeKind == refl.Int16 ||
			actualTypeKind == refl.Uint16 ||
			actualTypeKind == refl.Int32 ||
			actualTypeKind == refl.Int
	}

	if expectedType == "long" {
		return actualTypeKind == refl.Int8 ||
			actualTypeKind == refl.Uint8 ||
			actualTypeKind == refl.Int16 ||
			actualTypeKind == refl.Uint16 ||
			actualTypeKind == refl.Int32 ||
			actualTypeKind == refl.Uint32 ||
			actualTypeKind == refl.Int64 ||
			actualTypeKind == refl.Uint64 ||
			actualTypeKind == refl.Int ||
			actualTypeKind == refl.Uint
	}

	if expectedType == "float" {
		return actualTypeKind == refl.Float32
	}

	if expectedType == "double" {
		return actualTypeKind == refl.Float32 ||
			actualTypeKind == refl.Float64
	}

	if expectedType == "string" {
		return actualTypeKind == refl.String
	}

	if expectedType == "bool" || expectedType == "boolean" {
		return actualTypeKind == refl.Bool
	}

	if expectedType == "date" || expectedType == "datetime" {
		return actualType == refl.TypeOf(time.Time{})
	}

	if expectedType == "timespan" || expectedType == "duration" {
		return actualTypeKind == refl.Int8 ||
			actualTypeKind == refl.Uint8 ||
			actualTypeKind == refl.Int16 ||
			actualTypeKind == refl.Uint16 ||
			actualTypeKind == refl.Int32 ||
			actualTypeKind == refl.Uint32 ||
			actualTypeKind == refl.Int64 ||
			actualTypeKind == refl.Uint64 ||
			actualTypeKind == refl.Int ||
			actualTypeKind == refl.Uint ||
			actualTypeKind == refl.Float32 ||
			actualTypeKind == refl.Float64 ||
			actualType == refl.TypeOf(time.Duration(1))
	}

	if expectedType == "map" || expectedType == "dict" || expectedType == "dictionary" {
		return actualTypeKind == refl.Map
	}

	if expectedType == "array" || expectedType == "list" {
		return actualTypeKind == refl.Array ||
			actualTypeKind == refl.Slice
	}

	if strings.HasSuffix(expectedType, "[]") {
		if actualTypeKind == refl.Slice || actualTypeKind == refl.Array {
			expectedType = expectedType[:len(expectedType)-2]
			actualType = actualType.Elem()
			return c.MatchTypeByName(expectedType, actualType)
		}
	}

	return false
}
