package query

// SortParams Defines a field name and order used to sort query results.
//	see SortField
//
//	Example:
//		filter := NewFilterParamsFromTuples("type", "Type1");
//		paging := NewPagingParams(0, 100);
//		sorting := NewSortingParams(NewSortField("create_time", true));
//
//		myDataClient.getDataByFilter(filter, paging, sorting, (err, page) => {...});
type SortParams []SortField

// NewEmptySortParams creates a new instance.
//	Returns: *SortParams
func NewEmptySortParams() *SortParams {
	c := make(SortParams, 0, 10)
	return &c
}

// NewSortParams creates a new instance and initializes it with specified sort fields.
//	Parameters
//		- fields []SortField a list of fields to sort by.
//	Returns: *SortParams
func NewSortParams(fields []SortField) *SortParams {
	c := make(SortParams, len(fields))
	copy(c, fields)
	return &c
}
