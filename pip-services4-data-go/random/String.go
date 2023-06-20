package random

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// String random generator for string values.
//	Example:
//		value1 := random.String.PickChar("ABC");     // Possible result: "C"
//		value2 := random.String.Pick(["A","B","C"]); // Possible result: "B"
var String = &_TRandomString{}

type _TRandomString struct{}

const _digits = "01234956789"
const _symbols = "_,.:-/.[].{},#-!,$=%.+^.&*-() "
const _alphaLower = "abcdefghijklmnopqrstuvwxyz"
const _alphaUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const _alpha = _alphaUpper + _alphaLower
const _chars = _alpha + _digits + _symbols

// PickChar picks a random character from a string.
//	Parameters:
//		- values: string to pick a char from
//	Returns: randomly picked char.
func (c *_TRandomString) PickChar(values string) byte {
	if len(values) == 0 {
		return 0
	}

	index := Integer.Next(0, len(values))
	return values[index]
}

// Pick picks a random string from an array of string.
//	Parameters:
//		- values: string[] strings to pick from.
//	Returns: randomly picked string.
func (c *_TRandomString) Pick(values []string) string {
	if len(values) == 0 {
		return ""
	}

	index := Integer.Next(0, len(values))
	return values[index]
}

// Distort distorts a string by randomly replacing characters in it.
//	Parameters:
//		-value: string - a string to distort.
//	Returns: a distored string.
func (c *_TRandomString) Distort(value string) string {
	if value == "" {
		return ""
	}

	value = strings.ToLower(value)

	//Capitalize the first letter of the string 'value'.
	if Boolean.Chance(1, 5) {
		r, n := utf8.DecodeRuneInString(value)
		value = string(unicode.ToUpper(r)) + value[n:]
	}

	//Add a symbol to the end of the string 'value'
	if Boolean.Chance(1, 3) {
		value = value + string([]byte{String.PickChar(_symbols)})
	}

	return value
}

// NextAlphaChar generates random alpha characted [A-Za-z]
//	Returns: a random characted.
func (c *_TRandomString) NextAlphaChar() byte {
	index := Integer.Next(0, len(_alpha))
	return _alpha[index]
}

// Next generates a random string, consisting of upper and lower case letters (of the English alphabet),
// digits (0-9), and symbols ("_,.:-/.[].{},#-!,$=%.+^.&*-() ").
//	Parameters:
//		- minLength: int - minimum string length.
//		- maxLength: int - maximum string length.
//	Returns: a random string.
func (c *_TRandomString) Next(minLength int, maxLength int) string {
	length := Integer.Next(minLength, maxLength)
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		index := Integer.Next(0, len(_chars))
		result[i] = _chars[index]
	}

	return string(result)
}
