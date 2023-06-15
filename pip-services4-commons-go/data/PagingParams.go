package data

// PagingParams is a data transfer object to pass paging parameters for queries.
//
//	The page is defined by two parameters:
//		the skip parameter defines number of items to skip.
//		the take parameter sets how many items to return in a page.
//	additionally, the optional total parameter tells to return total number of items in the query.
//	Remember: not all implementations support the total parameter because its generation may lead to
//	severe performance implications.
//
//	Example:
//		filter := NewFilterParamsFromTuples("type", "Type1");
//		paging := NewPagingParams(0, 100);
//
//		err, page = myDataClient.GetDataByFilter(context.Background(), filter, paging);
type PagingParams struct {
	Skip  int64 `json:"skip"`
	Take  int64 `json:"take"`
	Total bool  `json:"total"`
}

const DefaultSkip int64 = 0
const DefaultTake int64 = 50

//	NewEmptyPagingParams creates a new instance.
//	Returns: *PagingParams
func NewEmptyPagingParams() *PagingParams {
	return &PagingParams{Skip: DefaultSkip, Take: DefaultTake, Total: false}
}

//	NewPagingParams creates a new instance and sets its values.
//	Parameters:
//		- skip the number of items to skip.
//		- take the number of items to return.
//		- total true to return the total number of items.
//	Returns: *PagingParams
func NewPagingParams(skip, take int64, total bool) *PagingParams {
	return &PagingParams{
		Skip:  skip,
		Take:  take,
		Total: total,
	}
}

// NewPagingParamsFromValue converts specified value into PagingParams.
//	Parameters: value any value to be converted
//	Returns: *PagingParams a newly created PagingParams.
func NewPagingParamsFromValue(value any) *PagingParams {
	if v, ok := value.(*PagingParams); ok {
		return v
	}
	return NewPagingParamsFromMap(NewAnyValueMapFromValue(value))
}

// NewPagingParamsFromTuples creates a new PagingParams from a list of key-value pairs called tuples.
//	Parameters: tuples ...any a list of values where odd elements are
//		keys and the following even elements are values
//	Returns: *PagingParams a newly created PagingParams.
func NewPagingParamsFromTuples(tuples ...any) *PagingParams {
	return NewPagingParamsFromMap(NewAnyValueMapFromTuplesArray(tuples))
}

// NewPagingParamsFromMap creates a new PagingParams and sets it parameters from the specified map
//	Parameters: value AnyValueMap or StringValueMap to initialize this PagingParams
//	Returns: *PagingParams a newly created PagingParams.
func NewPagingParamsFromMap(value *AnyValueMap) *PagingParams {
	return &PagingParams{
		Skip:  value.GetAsLongWithDefault("skip", DefaultSkip),
		Take:  value.GetAsLongWithDefault("take", DefaultTake),
		Total: value.GetAsBooleanWithDefault("total", false),
	}
}

// GetSkip gets the number of items to skip.
//	Parameters: minSkip int64 the minimum number of items to skip.
//	Returns: int64 the number of items to skip.
func (c *PagingParams) GetSkip(minSkip int64) int64 {
	if c.Skip < minSkip {
		return minSkip
	}
	return c.Skip
}

// GetTake gets the number of items to return in a page.
//	Parameters: maxTake int64 the maximum number of items to return.
//	Returns int64 the number of items to return.
func (c *PagingParams) GetTake(maxTake int64) int64 {
	if c.Take < 0 {
		return 0
	}
	if c.Take > maxTake {
		return maxTake
	}
	return c.Take
}
