package data

// FilterParams is a data transfer object used to pass filter parameters as simple key-value pairs.
//	see StringValueMap
//
//	Example:
//		filter := NewFilterParamsFromTuples(
//			"type", "Type1",
//			"from_create_time", time.Now(),
//			"to_create_time", time.Now().Add(10*time.Hour),
//			"completed", true,
//		)
//		paging = NewPagingParams(0, 100)
//
//		err, page = myDataClient.getDataByFilter(filter, paging)
type FilterParams struct {
	*StringValueMap
}

// NewEmptyFilterParams сreates a new instance.
//	Returns: *FilterParams
func NewEmptyFilterParams() *FilterParams {
	return &FilterParams{
		StringValueMap: NewEmptyStringValueMap(),
	}
}

// NewFilterParams creates a new instance and initialize it with elements from the specified map.
//	Parameters:
//		- value map[string]string a map to initialize this instance.
//	Returns: *FilterParams
func NewFilterParams(values map[string]string) *FilterParams {
	return &FilterParams{
		StringValueMap: NewStringValueMapFromMaps(values),
	}
}

// NewFilterParamsFromValue converts specified value into FilterParams.
//	Parameters:
//		- value interface value to be converted
//	Returns: a newly created FilterParams.
func NewFilterParamsFromValue(value any) *FilterParams {
	return &FilterParams{
		StringValueMap: NewStringValueMapFromValue(value),
	}
}

// NewFilterParamsFromTuples creates a new FilterParams from a list of key-value pairs called tuples.
//	Parameters:
//		- tuples ...any a list of values where odd
//		elements are keys and the following even elements are values
//	Returns: *FilterParams a newly created FilterParams.
func NewFilterParamsFromTuples(tuples ...any) *FilterParams {
	return &FilterParams{
		StringValueMap: NewStringValueMapFromTuplesArray(tuples),
	}
}

// NewFilterParamsFromString parses semicolon-separated key-value pairs and returns them as a FilterParams.
// see StringValueMap.FromString
//	Parameters:
//		- line string semicolon-separated key-value list to initialize FilterParams.
//	Returns: *FilterParams
func NewFilterParamsFromString(line string) *FilterParams {
	return &FilterParams{
		StringValueMap: NewStringValueMapFromString(line),
	}
}
