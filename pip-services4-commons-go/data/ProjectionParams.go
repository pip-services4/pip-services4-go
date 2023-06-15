package data

import (
	"strings"
)

// ProjectionParams defines projection parameters with list if fields to include into query results.
// The parameters support two formats: dot format and nested format.
// The dot format is the standard way to define included fields and subfields
// using dot object notation: "field1,field2.field21,field2.field22.field221".
// As alternative the nested format offers a more compact representation: "field1,field2(field21,field22(field221))".
//
//	Example:
//		filter := NewFilterParamsFromTuples("type", "Type1");
//		paging := NewPagingParams(0, 100);
//		projection = NewProjectionParamsFromString("field1,field2(field21,field22)")
//
//		err, page := myDataClient.GetDataByFilter(context.Background(), filter, paging, projection);
type ProjectionParams struct {
	_values []string
}

// NewEmptyProjectionParams creates a new instance of the projection parameters and assigns its value.
//	Returns: *ProjectionParams
func NewEmptyProjectionParams() *ProjectionParams {
	return &ProjectionParams{
		_values: make([]string, 0, 10),
	}
}

// NewProjectionParamsFromStrings creates a new instance of the projection parameters and assigns its from string value.
//	Parameters: values []string
//	Returns: *ProjectionParams
func NewProjectionParamsFromStrings(values []string) *ProjectionParams {
	c := &ProjectionParams{
		_values: make([]string, len(values)),
	}
	copy(c._values, values)
	return c
}

// NewProjectionParamsFromAnyArray creates a new instance of the projection parameters and assigns
// its from AnyValueArray values.
//	Parameters: values *AnyValueArray
//	Returns: *ProjectionParams
func NewProjectionParamsFromAnyArray(values *AnyValueArray) *ProjectionParams {
	if values == nil {
		return NewEmptyProjectionParams()
	}

	c := &ProjectionParams{
		_values: make([]string, 0, values.Len()),
	}

	for index := 0; index < values.Len(); index++ {
		value := values.GetAsString(index)
		if value != "" {
			c._values = append(c._values, value)
		}
	}

	return c
}

// NewProjectionParamsFromValue converts specified value into ProjectionParams.
//	see AnyValueArray.fromValue
//	Parameters: value any value to be converted
//	Returns: *ProjectionParams a newly created ProjectionParams.
func NewProjectionParamsFromValue(value any) *ProjectionParams {
	return NewProjectionParamsFromAnyArray(NewAnyValueArrayFromValue(value))
}

// ParseProjectionParams create new ProjectionParams and set values from values
//	Parameters: values ...string a values to parse
//	Returns: *ProjectionParams
func ParseProjectionParams(values ...string) *ProjectionParams {
	c := NewEmptyProjectionParams()

	for index := 0; index < len(values); index++ {
		parseProjectionParamValue("", c, values[index])
	}

	return c
}

// Value return raw values []string
func (c *ProjectionParams) Value() []string {
	return c._values
}

// Len gets or sets the length of the array. This is a number one
// higher than the highest element defined in an array.
func (c *ProjectionParams) Len() int {
	return len(c._values)
}

// Get value by index
// Parameters:
//  - index int
//  an index of element
// Return string
func (c *ProjectionParams) Get(index int) (any, bool) {
	if c.IsValidIndex(index) {
		return c._values[index], true
	}
	return nil, false
}

// IsValidIndex checks that 0 <= index < len.
//	Parameters:
//		index int an index of the element to get.
// Returns: bool
func (c *ProjectionParams) IsValidIndex(index int) bool {
	return index >= 0 && index < c.Len()
}

// Put value in index position
//	Parameters:
//		- index int an index of element
//		- value string value
func (c *ProjectionParams) Put(index int, value string) bool {
	if index <= 0 && index <= c.Len() {
		after := c._values[index:]
		before := c._values[:index]
		c._values = append(make([]string, 0, len(c._values)+1), before...)
		c._values = append(c._values, value)
		c._values = append(c._values, after...)
		return true
	}
	return false
}

// Remove element by index
// Parameters:
//  - index int
//  an index of remove element
func (c *ProjectionParams) Remove(index int) {
	c._values = append(c._values[:index], c._values[index+1:]...)
}

// Push new element to an array.
//	Parameters: value string
func (c *ProjectionParams) Push(value string) {
	c._values = append(c._values, value)
}

// Append new elements to an array.
//	Parameters: value []string
func (c *ProjectionParams) Append(elements []string) {
	if elements != nil {
		c._values = append(c._values, elements...)
	}
}

// Clear elements
func (c *ProjectionParams) Clear() {
	c._values = make([]string, 0, 10)
}

// String returns a string representation of an array.
//	Returns: string
func (c *ProjectionParams) String() string {
	builder := strings.Builder{}
	if c.Len() == 0 {
		return ""
	}
	builder.WriteString(c._values[0])
	for index := 1; index < c.Len(); index++ {
		builder.WriteString(",")
		builder.WriteString(c._values[index])
	}

	return builder.String()
}

// parseProjectionParamValue Add parse value into exist ProjectionParams and add prefix
//	Parameters:
//		- prefix string prefix value
//		- c *ProjectionParams ProjectionParams instance wheare need to add value
//		- value string a values to parse
func parseProjectionParamValue(prefix string, c *ProjectionParams, value string) {
	if value != "" {
		value = strings.Trim(value, " \t\n\r")
	}

	openBracket := 0
	openBracketIndex := -1
	closeBracketIndex := -1
	commaIndex := -1

	breakCycleRequired := false
	for index := 0; index < len(value); index++ {
		switch value[index] {
		case '(':
			if openBracket == 0 {
				openBracketIndex = index
			}

			openBracket++
			break
		case ')':
			openBracket--

			if openBracket == 0 {
				closeBracketIndex = index

				if openBracketIndex >= 0 && closeBracketIndex > 0 {
					previousPrefix := prefix

					if prefix != "" {
						prefix = prefix + "." + value[:openBracketIndex]
					} else {
						prefix = value[:openBracketIndex]
					}

					subValue := value[openBracketIndex+1 : closeBracketIndex]
					parseProjectionParamValue(prefix, c, subValue)

					subValue = value[closeBracketIndex+1:]
					parseProjectionParamValue(previousPrefix, c, subValue)
					breakCycleRequired = true
				}
			}
			break
		case ',':
			if openBracket == 0 {
				commaIndex = index

				subValue := value[0:commaIndex]

				if subValue != "" {
					if prefix != "" {
						c.Push(prefix + "." + subValue)
					} else {
						c.Push(subValue)
					}
				}

				subValue = value[commaIndex+1:]

				if subValue != "" {
					parseProjectionParamValue(prefix, c, subValue)
					breakCycleRequired = true
				}
			}
			break
		}

		if breakCycleRequired {
			break
		}
	}

	if value != "" && openBracketIndex == -1 && commaIndex == -1 {
		if prefix != "" {
			c.Push(prefix + "." + value)
		} else {
			c.Push(value)
		}
	}
}
